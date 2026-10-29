package field

import (
	"errors"
	"fmt"
	"github.com/Team254/cheesy-arena/game"
	"github.com/Team254/cheesy-arena/model"
	"github.com/Team254/cheesy-arena/plc"
	"github.com/google/uuid"
	"log"
	"math"
	"time"
)

var errNoScoreChange = errors.New("no score change")

type FuelCounts struct {
	Source           string
	Auto             int
	Teleop           int
	RawAuto          int
	RawTeleop        int
	AutoAdjustment   int
	TeleopAdjustment int
	Suspended        string
}

type FuelState struct {
	HistoryId    string
	AttemptId    string
	MatchId      int
	MatchName    string
	Revision     uint64
	Mode         string
	Phase        MatchState
	CanAuto      bool
	CanTeleop    bool
	CanSetMode   bool
	CanRebase    bool
	Recovered    bool
	CheckpointAt time.Time
	StorageError string
	Red          FuelCounts
	Blue         FuelCounts
}

type FuelReply struct {
	RequestId string
	Error     string
	Duplicate bool
	State     FuelState
}

func cloneGameScore(score game.Score) game.Score {
	score.Fouls = append([]game.Foul{}, score.Fouls...)
	return score
}

func cloneCards(cards map[string]string) map[string]string {
	copy := make(map[string]string, len(cards))
	for key, value := range cards {
		copy[key] = value
	}
	return copy
}

// ScoreResultSnapshot never exposes mutable live score pointers to JSON writers.
func (arena *Arena) ScoreResultSnapshot() *model.MatchResult {
	arena.scoreMutex.RLock()
	defer arena.scoreMutex.RUnlock()
	return arena.scoreResultLocked()
}

func (arena *Arena) scoreResultLocked() *model.MatchResult {
	red := cloneGameScore(arena.RedRealtimeScore.CurrentScore)
	blue := cloneGameScore(arena.BlueRealtimeScore.CurrentScore)
	metadata := arena.fuel.FuelMetadata
	if arena.fuel.ResultAttemptId != "" {
		metadata.AttemptId = arena.fuel.ResultAttemptId
	}
	return &model.MatchResult{
		FuelContext: arena.fuel.AttemptId, FuelRevision: arena.fuel.Revision,
		MatchId: arena.fuel.Match.Id, MatchType: arena.fuel.Match.Type,
		RedScore: &red, BlueScore: &blue,
		RedCards: cloneCards(arena.RedRealtimeScore.Cards), BlueCards: cloneCards(arena.BlueRealtimeScore.Cards),
		PlayoffRedAllianceCard:  arena.RedRealtimeScore.PlayoffAllianceCard,
		PlayoffBlueAllianceCard: arena.BlueRealtimeScore.PlayoffAllianceCard,
		FuelControl:             &metadata,
	}
}

// EditLiveScore serializes the existing referee/scoring writers with FUEL updates.
func (arena *Arena) EditLiveScore(edit func(red, blue *RealtimeScore) error) error {
	arena.fuelOperationMutex.Lock()
	defer arena.fuelOperationMutex.Unlock()
	arena.scoreMutex.Lock()
	if arena.fuel.Status == "committing" || arena.fuel.Status == "committed" || arena.fuel.Status == "discarded" {
		arena.scoreMutex.Unlock()
		return fmt.Errorf("score is closed; load a match or edit the recorded result")
	}
	err := edit(arena.RedRealtimeScore, arena.BlueRealtimeScore)
	if err == errNoScoreChange {
		arena.scoreMutex.Unlock()
		return nil
	}
	if err == nil {
		arena.fuel.Revision++
	}
	arena.scoreMutex.Unlock()
	if err == nil {
		arena.RealtimeScoreNotifier.Notify()
	}
	return err
}

func effectiveFuel(source model.FuelAlliance, mode string) game.Hub {
	hub := source.Raw
	for shift := game.ShiftAuto; shift < game.ShiftCount; shift++ {
		if mode == model.FuelManual {
			hub.ShiftCounts[shift] = 0
		}
		hub.ShiftCounts[shift] += source.Adjustment[shift]
	}
	return hub
}

func (arena *Arena) projectFuelLocked() {
	arena.RedRealtimeScore.CurrentScore.Hub = effectiveFuel(arena.fuel.Red, arena.fuel.Mode)
	arena.BlueRealtimeScore.CurrentScore.Hub = effectiveFuel(arena.fuel.Blue, arena.fuel.Mode)
}

func (arena *Arena) newFuelAttemptLocked(match *model.Match) {
	mode := arena.fuel.Mode
	revision := arena.fuel.Revision + 1
	if mode != model.FuelManual {
		mode = model.FuelPlcCorrection
	}
	arena.fuel = model.FuelCheckpoint{
		FuelMetadata: model.FuelMetadata{AttemptId: uuid.NewString(), Mode: mode, RedSource: arena.automaticFuelSources["red"], BlueSource: arena.automaticFuelSources["blue"]},
		Match:        *match, Status: "prepared", Phase: int(PreMatch), Revision: revision,
	}
	arena.fuel.Red.Suspended = "Waiting for a fresh PLC zero after reset"
	arena.fuel.ResultAttemptId = arena.fuel.AttemptId
	arena.fuel.Blue.Suspended = arena.fuel.Red.Suspended
	if arena.fuel.RedSource == "arduino" {
		arena.fuel.Red.Suspended = ""
	}
	if arena.fuel.BlueSource == "arduino" {
		arena.fuel.Blue.Suspended = ""
	}
}

func (arena *Arena) resetFuelCursor() {
	arena.scoreMutex.Lock()
	defer arena.scoreMutex.Unlock()
	for index, source := range []*model.FuelAlliance{&arena.fuel.Red, &arena.fuel.Blue} {
		if index == 0 && arena.fuel.RedSource == "arduino" || index == 1 && arena.fuel.BlueSource == "arduino" {
			source.Suspended = ""
			continue
		}
		source.Armed = false
		source.Suspended = "Waiting for a fresh PLC zero after reset"
		source.LastCount = 0
	}
	arena.fuel.Revision++
}

func (arena *Arena) checkpointLocked() *model.FuelCheckpoint {
	checkpoint := arena.fuel
	checkpoint.Result = *arena.scoreResultLocked()
	checkpoint.UpdatedAt = time.Now()
	return &checkpoint
}

// CheckpointFuel is also used as a lifecycle barrier. The operation mutex is
// acquired before CAPTURING the snapshot, never after queuing an old snapshot.
func (arena *Arena) CheckpointFuel() error {
	arena.fuelOperationMutex.Lock()
	defer arena.fuelOperationMutex.Unlock()
	return arena.persistFuelLocked()
}

func (arena *Arena) persistFuelLocked() error {
	arena.scoreMutex.RLock()
	if arena.fuelPersistenceUncertain {
		err := fmt.Errorf("FUEL durability uncertain; restart to reconcile: %s", arena.fuelStorageError)
		arena.scoreMutex.RUnlock()
		return err
	}
	checkpoint := arena.checkpointLocked()
	arena.scoreMutex.RUnlock()
	err := arena.Database.SaveFuelCheckpoint(checkpoint, nil)
	arena.scoreMutex.Lock()
	if err != nil {
		arena.fuelStorageError = err.Error()
	} else {
		arena.fuelStorageError = ""
		arena.fuel.UpdatedAt = checkpoint.UpdatedAt
	}
	arena.scoreMutex.Unlock()
	return err
}

func (arena *Arena) runFuelCheckpoints() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	var previousRevision uint64
	var previousAttempt string
	for range ticker.C {
		arena.scoreMutex.RLock()
		revision, attempt, started := arena.fuel.Revision, arena.fuel.AttemptId, arena.fuel.Started
		arena.scoreMutex.RUnlock()
		if !started || revision == previousRevision && attempt == previousAttempt {
			continue
		}
		if err := arena.CheckpointFuel(); err != nil {
			log.Printf("FUEL checkpoint failed: %v", err)
			continue
		}
		previousRevision, previousAttempt = revision, attempt
	}
}

func (arena *Arena) FuelState() FuelState {
	arena.scoreMutex.RLock()
	defer arena.scoreMutex.RUnlock()
	return arena.fuelStateLocked()
}

func (arena *Arena) fuelStateLocked() FuelState {
	fuel := arena.fuel
	editable := fuel.Status == "active" && fuel.Started && arena.fuelStorageError == ""
	state := FuelState{
		HistoryId: fuel.ResultAttemptId,
		AttemptId: fuel.AttemptId, MatchId: fuel.Match.Id, MatchName: fuel.Match.LongName,
		Revision: fuel.Revision, Mode: fuel.Mode, Phase: MatchState(fuel.Phase),
		CanAuto:    editable && fuel.Phase >= int(AutoPeriod) && fuel.Phase <= int(PostMatch),
		CanTeleop:  editable && (fuel.Phase == int(TeleopPeriod) || fuel.Phase == int(PostMatch)),
		CanSetMode: fuel.Status == "prepared" && !fuel.Started && fuel.Phase == int(PreMatch),
		CanRebase:  editable && fuel.Mode == model.FuelPlcCorrection && !fuel.Recovered,
		Recovered:  fuel.Recovered, CheckpointAt: fuel.UpdatedAt, StorageError: arena.fuelStorageError,
	}
	counts := func(source model.FuelAlliance, owner string) FuelCounts {
		if owner == "" {
			owner = "plc"
		}
		hub := effectiveFuel(source, fuel.Mode)
		result := FuelCounts{Source: owner, Auto: hub.GetShiftCount(game.ShiftAuto, true), Teleop: hub.GetTeleopActiveFuelCount(),
			RawAuto: source.Raw.GetShiftCount(game.ShiftAuto, true), RawTeleop: source.Raw.GetTeleopActiveFuelCount(), Suspended: source.Suspended}
		result.AutoAdjustment = result.Auto - result.RawAuto
		result.TeleopAdjustment = result.Teleop - result.RawTeleop
		return result
	}
	state.Red, state.Blue = counts(fuel.Red, fuel.RedSource), counts(fuel.Blue, fuel.BlueSource)
	return state
}

// syncFuelPhase is cheap and performs no disk or network I/O in the arena loop.
func (arena *Arena) syncFuelPhase() {
	arena.scoreMutex.Lock()
	defer arena.scoreMutex.Unlock()
	arena.fuel.ClockStart = arena.MatchStartTime
	arena.fuel.Aborted, arena.fuel.StoppedAt = arena.matchAborted, arena.matchStopTime
	if arena.fuel.Phase != int(arena.MatchState) {
		arena.fuel.Phase = int(arena.MatchState)
		arena.fuel.Revision++
	}
	if arena.MatchState >= AutoPeriod && arena.MatchState <= PostMatch && arena.fuel.Status == "prepared" {
		arena.fuel.Started = true
		arena.fuel.Status = "active"
	}
}

func (arena *Arena) processFuelSnapshot(snapshot plc.FuelSnapshot, now time.Time) {
	arena.scoreMutex.Lock()
	defer arena.scoreMutex.Unlock()
	if arena.fuel.Mode == model.FuelManual || arena.fuel.Recovered || arena.fuel.Status != "active" && arena.fuel.Status != "prepared" {
		return
	}
	if arena.matchAborted && now.After(arena.matchStopTime.Add(time.Duration(game.ScoringGracePeriodSec)*time.Second)) {
		return
	}
	changed := false
	for index, source := range []*model.FuelAlliance{&arena.fuel.Red, &arena.fuel.Blue} {
		if index == 0 && arena.fuel.RedSource == "arduino" || index == 1 && arena.fuel.BlueSource == "arduino" {
			continue
		}
		count := snapshot.Red
		if index == 1 {
			count = snapshot.Blue
		}
		old := *source
		valid := snapshot.Healthy && !snapshot.ReadAt.IsZero() && now.Sub(snapshot.ReadAt) <= 500*time.Millisecond && !snapshot.ReadAt.After(now.Add(50*time.Millisecond))
		if !valid {
			if source.Armed {
				source.Suspended = "PLC disconnected or sample older than 500 ms"
			}
		} else if snapshot.Sequence > source.Sequence {
			if source.Armed && (snapshot.Connection != source.Connection || snapshot.Reset != source.Reset) {
				source.Suspended = "PLC connection/reset changed; operator rebase required"
			}
			if !source.Armed && count == 0 {
				source.Armed = true
				source.Suspended = ""
			}
			if source.Armed && source.Suspended == "" && count < source.LastCount {
				source.Suspended = "PLC counter decreased; operator rebase required"
			}
			if source.Armed && source.Suspended == "" {
				if arena.fuel.Started && arena.fuel.Status == "active" {
					total := 0
					for _, value := range source.Raw.ShiftCounts {
						total += value
					}
					source.Raw.UpdateState(total+count-source.LastCount, arena.fuel.ClockStart, snapshot.ReadAt)
				}
				source.LastCount = count
			}
			source.Sequence, source.Connection, source.Reset = snapshot.Sequence, snapshot.Connection, snapshot.Reset
		}
		changed = changed || old.Raw != source.Raw || old.Suspended != source.Suspended || old.Armed != source.Armed
	}
	if changed {
		arena.projectFuelLocked()
		arena.fuel.Revision++
	}
}

func changeFuel(source *model.FuelAlliance, mode, phase string, delta int, state MatchState, start, now time.Time) (int, int, error) {
	hub := effectiveFuel(*source, mode)
	before := hub.GetShiftCount(game.ShiftAuto, true)
	if phase == "teleop" {
		before = hub.GetTeleopActiveFuelCount()
	}
	if delta < -before || delta > math.MaxInt32-before {
		return before, before, fmt.Errorf("FUEL cannot be negative or overflow")
	}
	if phase == "auto" {
		source.Adjustment[game.ShiftAuto] += delta
		return before, before + delta, nil
	}
	if delta > 0 {
		shift := game.ShiftEndgame
		if state != PostMatch {
			current, _, _, ok := hub.GetCurrentShiftTiming(start, now)
			if !ok || current < game.ShiftTransition {
				current = game.ShiftTransition
			}
			shift = game.ShiftTransition
			for candidate := game.ShiftTransition; candidate <= current && candidate <= game.ShiftEndgame; candidate++ {
				// Only credited shifts can receive a scored manual FUEL.
				probe := game.Hub{WonAuto: hub.WonAuto}
				probe.ShiftCounts[candidate] = 1
				if probe.GetShiftCount(candidate, true) == 1 {
					shift = candidate
				}
			}
		}
		source.Adjustment[shift] += delta
	} else {
		remaining := -delta
		for shift := game.ShiftPostMatch; shift >= game.ShiftTransition && remaining > 0; shift-- {
			remove := hub.GetShiftCount(shift, true)
			if remove > remaining {
				remove = remaining
			}
			source.Adjustment[shift] -= remove
			remaining -= remove
		}
	}
	return before, before + delta, nil
}

// ApplyFuelCommand serializes operator writes without holding the score mutex
// during persistence. PLC deltas may continue; they never change adjustments.
func (arena *Arena) ApplyFuelCommand(kind string, command model.FuelCommand, actor string) FuelReply {
	arena.fuelOperationMutex.Lock()
	defer arena.fuelOperationMutex.Unlock()
	reply := FuelReply{RequestId: command.RequestId}
	fail := func(err error) FuelReply { reply.Error = err.Error(); reply.State = arena.FuelState(); return reply }
	if _, err := uuid.Parse(command.RequestId); err != nil {
		return fail(fmt.Errorf("invalid RequestId"))
	}
	if len(command.Operator) > 64 {
		return fail(fmt.Errorf("operator label too long"))
	}
	arena.scoreMutex.RLock()
	uncertain := arena.fuelPersistenceUncertain
	arena.scoreMutex.RUnlock()
	if uncertain {
		return fail(fmt.Errorf("FUEL persistence is uncertain; restart to reconcile the durable state"))
	}
	previous, err := arena.Database.GetFuelOperation(command.RequestId)
	if err != nil {
		return fail(err)
	}
	if previous != nil {
		if previous.Command != command || previous.Kind != kind {
			return fail(fmt.Errorf("RequestId reused with different contents"))
		}
		reply.Duplicate = true
		reply.State = arena.FuelState()
		return reply
	}
	arena.scoreMutex.RLock()
	candidate := arena.checkpointLocked()
	state := arena.fuelStateLocked()
	start := arena.fuel.ClockStart
	arena.scoreMutex.RUnlock()
	if command.AttemptId != state.AttemptId || command.MatchId != state.MatchId {
		return fail(fmt.Errorf("stale match/attempt; reload the current state"))
	}
	if state.StorageError != "" {
		return fail(fmt.Errorf("FUEL storage unavailable: %s", state.StorageError))
	}
	operation := &model.FuelOperation{Command: command, Kind: kind, Actor: actor, Time: time.Now()}
	source := &candidate.Red
	if command.Alliance == "blue" {
		source = &candidate.Blue
	}
	operation.ResultAttemptId = candidate.ResultAttemptId
	operation.PreviousBase = source.LastCount
	switch kind {
	case "adjustFuel":
		if command.Alliance != "red" && command.Alliance != "blue" {
			return fail(fmt.Errorf("invalid alliance"))
		}
		if command.Phase != "auto" && command.Phase != "teleop" {
			return fail(fmt.Errorf("invalid FUEL phase"))
		}
		if command.Delta != 1 && command.Delta != -1 && command.Delta != 5 && command.Delta != -5 {
			return fail(fmt.Errorf("delta must be +1, -1, +5 or -5"))
		}
		if command.Phase == "auto" && !state.CanAuto || command.Phase == "teleop" && !state.CanTeleop {
			return fail(fmt.Errorf("FUEL editing is closed for this match/phase"))
		}
		operation.Before, operation.After, err = changeFuel(source, candidate.Mode, command.Phase, command.Delta, state.Phase, start, time.Now())
		if err != nil {
			return fail(err)
		}
	case "setFuelMode":
		if !state.CanSetMode || state.Red.Auto != 0 || state.Red.Teleop != 0 || state.Blue.Auto != 0 || state.Blue.Teleop != 0 {
			return fail(fmt.Errorf("mode can only change before an empty match starts"))
		}
		if command.Mode != model.FuelManual && command.Mode != model.FuelPlcCorrection {
			return fail(fmt.Errorf("invalid FUEL mode"))
		}
		candidate.Mode = command.Mode
	case "rebasePlcFuel":
		if !state.CanRebase || command.Alliance != "red" && command.Alliance != "blue" || source.Suspended == "" {
			return fail(fmt.Errorf("PLC rebase is not available"))
		}
		if command.Alliance == "red" && candidate.RedSource == "arduino" || command.Alliance == "blue" && candidate.BlueSource == "arduino" {
			return fail(fmt.Errorf("this alliance is owned by Arduino, not PLC"))
		}
		snapshot := arena.Plc.GetFuelSnapshot()
		age := time.Since(snapshot.ReadAt)
		if !snapshot.Healthy || snapshot.ReadAt.IsZero() || age > 500*time.Millisecond || age < -50*time.Millisecond {
			return fail(fmt.Errorf("a fresh healthy PLC sample is required"))
		}
		source.LastCount = snapshot.Red
		if command.Alliance == "blue" {
			source.LastCount = snapshot.Blue
		}
		source.Sequence, source.Connection, source.Reset = snapshot.Sequence, snapshot.Connection, snapshot.Reset
		source.Armed, source.Suspended = true, ""
	default:
		return fail(fmt.Errorf("unknown FUEL command"))
	}
	candidate.Revision++
	operation.Mode, operation.BaseCount, operation.SourceSequence = candidate.Mode, source.LastCount, source.Sequence
	candidate.Result.RedScore.Hub = effectiveFuel(candidate.Red, candidate.Mode)
	candidate.Result.BlueScore.Hub = effectiveFuel(candidate.Blue, candidate.Mode)
	metadata := candidate.FuelMetadata
	if candidate.ResultAttemptId != "" {
		metadata.AttemptId = candidate.ResultAttemptId
	}
	candidate.Result.FuelControl = &metadata
	operation.Revision = candidate.Revision
	if err = arena.Database.SaveFuelCheckpoint(candidate, operation); err != nil {
		arena.scoreMutex.Lock()
		arena.fuelPersistenceUncertain = true
		arena.fuelStorageError = err.Error()
		arena.scoreMutex.Unlock()
		return fail(err)
	}
	arena.scoreMutex.Lock()
	arena.fuel.Mode = candidate.Mode
	arena.fuel.Red.Adjustment, arena.fuel.Blue.Adjustment = candidate.Red.Adjustment, candidate.Blue.Adjustment
	if kind == "rebasePlcFuel" {
		target := &arena.fuel.Red
		if command.Alliance == "blue" {
			target = &arena.fuel.Blue
		}
		target.LastCount, target.Sequence, target.Connection, target.Reset = source.LastCount, source.Sequence, source.Connection, source.Reset
		target.Armed, target.Suspended = true, ""
	}
	arena.projectFuelLocked()
	arena.fuel.Revision++
	arena.fuel.UpdatedAt = candidate.UpdatedAt
	reply.State = arena.fuelStateLocked()
	arena.scoreMutex.Unlock()
	arena.RealtimeScoreNotifier.Notify()
	arena.ManualFuelNotifier.Notify()
	return reply
}

// beginFuelCommit is used through BeginFuelCommit to exclude new commands until
// the existing result publishing pipeline has completed or failed.
func (arena *Arena) BeginFuelCommit() (*model.MatchResult, func(bool), error) {
	arena.fuelOperationMutex.Lock()
	arena.scoreMutex.Lock()
	if arena.MatchState != PostMatch && (arena.MatchState != PreMatch || arena.fuel.Started) || arena.fuel.Status == "committed" || arena.fuel.Status == "discarded" {
		arena.scoreMutex.Unlock()
		arena.fuelOperationMutex.Unlock()
		return nil, nil, fmt.Errorf("match is not awaiting results")
	}
	if arena.fuelGracePendingLocked() {
		arena.scoreMutex.Unlock()
		arena.fuelOperationMutex.Unlock()
		return nil, nil, fmt.Errorf("wait for the FUEL scoring grace period before posting")
	}
	arena.fuel.Status = "committing"
	arena.fuel.Revision++
	result := arena.scoreResultLocked()
	arena.scoreMutex.Unlock()
	if err := arena.persistFuelLocked(); err != nil {
		arena.scoreMutex.Lock()
		arena.fuel.Status = "active"
		arena.scoreMutex.Unlock()
		arena.fuelOperationMutex.Unlock()
		return nil, nil, err
	}
	// If a process stopped partway through committing, update the same result,
	// rather than creating a second PlayNumber on retry.
	previous, err := arena.Database.GetMatchResultForMatch(result.MatchId)
	if err != nil {
		arena.scoreMutex.Lock()
		arena.fuel.Status = "active"
		arena.scoreMutex.Unlock()
		arena.fuelOperationMutex.Unlock()
		return nil, nil, err
	}
	if previous != nil && previous.FuelControl != nil && previous.FuelControl.AttemptId == result.FuelControl.AttemptId {
		result.Id, result.PlayNumber = previous.Id, previous.PlayNumber
	}
	return result, func(success bool) {
		arena.scoreMutex.Lock()
		arena.fuel.Status = "active"
		if success {
			arena.fuel.Status = "committed"
		}
		arena.fuel.Revision++
		arena.scoreMutex.Unlock()
		if err := arena.persistFuelLocked(); err != nil {
			log.Printf("FUEL commit checkpoint failed: %v", err)
		}
		arena.fuelOperationMutex.Unlock()
		arena.ManualFuelNotifier.Notify()
	}, nil
}

func (arena *Arena) recoverFuel(checkpoint *model.FuelCheckpoint) error {
	if checkpoint == nil || !checkpoint.Started || checkpoint.Status == "committed" || checkpoint.Status == "discarded" {
		return nil
	}
	if checkpoint.Result.RedScore == nil || checkpoint.Result.BlueScore == nil {
		return fmt.Errorf("invalid FUEL recovery scores")
	}
	arena.CurrentMatch = &checkpoint.Match
	arena.fuel = *checkpoint
	// Rotate the browser context; retain the original attempt in result metadata
	// through RecoveryAttemptId so unfinished result commits remain idempotent.
	if arena.fuel.ResultAttemptId == "" {
		arena.fuel.ResultAttemptId = checkpoint.AttemptId
	}
	arena.fuel.AttemptId = uuid.NewString()
	arena.fuel.Status, arena.fuel.Phase = "active", int(PostMatch)
	arena.fuel.Recovered = true
	arena.fuel.Revision++
	arena.MatchState = PostMatch
	arena.MatchStartTime = time.Time{}
	arena.RedRealtimeScore.CurrentScore = cloneGameScore(*checkpoint.Result.RedScore)
	arena.BlueRealtimeScore.CurrentScore = cloneGameScore(*checkpoint.Result.BlueScore)
	arena.RedRealtimeScore.Cards = cloneCards(checkpoint.Result.RedCards)
	arena.BlueRealtimeScore.Cards = cloneCards(checkpoint.Result.BlueCards)
	arena.RedRealtimeScore.PlayoffAllianceCard = checkpoint.Result.PlayoffRedAllianceCard
	arena.BlueRealtimeScore.PlayoffAllianceCard = checkpoint.Result.PlayoffBlueAllianceCard
	for _, station := range arena.AllianceStations {
		station.Bypass = false
	}
	for _, score := range []*game.Score{checkpoint.Result.RedScore, checkpoint.Result.BlueScore} {
		for _, foul := range score.Fouls {
			if foul.FoulId >= arena.NextFoulId {
				arena.NextFoulId = foul.FoulId + 1
			}
		}
	}
	for station, teamId := range map[string]int{"R1": checkpoint.Match.Red1, "R2": checkpoint.Match.Red2, "R3": checkpoint.Match.Red3, "B1": checkpoint.Match.Blue1, "B2": checkpoint.Match.Blue2, "B3": checkpoint.Match.Blue3} {
		if teamId != 0 {
			team, err := arena.Database.GetTeamById(teamId)
			if err != nil {
				return err
			}
			arena.AllianceStations[station].Team = team
		}
	}
	return arena.persistFuelLocked()
}

func (arena *Arena) fuelGracePendingLocked() bool {
	if !arena.fuel.Started || arena.fuel.Recovered || arena.fuel.Mode == model.FuelManual || !arena.Plc.IsEnabled() && arena.fuel.RedSource != "arduino" && arena.fuel.BlueSource != "arduino" {
		return false
	}
	end := arena.MatchStartTime.Add(game.GetDurationToTeleopEnd())
	if arena.matchAborted {
		end = arena.matchStopTime
	}
	return time.Now().Before(end.Add(time.Duration(game.ScoringGracePeriodSec) * time.Second))
}
