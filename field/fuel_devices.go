package field

import (
	"fmt"
	"github.com/Team254/cheesy-arena/fueldevice"
	"github.com/Team254/cheesy-arena/game"
	"github.com/Team254/cheesy-arena/model"
	"github.com/google/uuid"
	"time"
)

// ConfigureFuelDevices is called once, before Arena.Run starts. Ownership cannot
// change during play. Recovery preserves the pending match's original metadata.
func (arena *Arena) ConfigureFuelDevices(config fueldevice.Config) error {
	if err := config.Validate(); err != nil {
		return err
	}
	arena.fuelOperationMutex.Lock()
	defer arena.fuelOperationMutex.Unlock()
	arena.scoreMutex.Lock()
	if arena.fuel.Started && !arena.fuel.Recovered {
		arena.scoreMutex.Unlock()
		return fmt.Errorf("cannot configure devices during an active or pending match")
	}
	arena.automaticFuelSources, arena.fuelDeviceIds = map[string]string{}, map[string]string{}
	if config.Enabled {
		for _, d := range config.Devices {
			arena.automaticFuelSources[d.Alliance], arena.fuelDeviceIds[d.Alliance] = d.Source, d.Id
		}
	}
	if !arena.fuel.Started {
		arena.fuel.RedSource, arena.fuel.BlueSource = arena.automaticFuelSources["red"], arena.automaticFuelSources["blue"]
		if arena.fuel.RedSource == "arduino" {
			arena.fuel.Red.Suspended = ""
		}
		if arena.fuel.BlueSource == "arduino" {
			arena.fuel.Blue.Suspended = ""
		}
		arena.fuel.Revision++
	}
	arena.scoreMutex.Unlock()
	return arena.persistFuelLocked()
}

func (arena *Arena) FuelDeviceState(alliance string, now time.Time) fueldevice.State {
	arena.scoreMutex.RLock()
	defer arena.scoreMutex.RUnlock()
	return arena.fuelDeviceStateLocked(alliance, now)
}

func (arena *Arena) fuelDeviceStateLocked(alliance string, now time.Time) fueldevice.State {
	fuel := arena.fuel
	state := fueldevice.State{Attempt: fuel.AttemptId, MatchId: fuel.Match.Id, Phase: fuel.Phase, Shift: int(game.ShiftCount)}
	source, owner := fuel.Red, fuel.RedSource
	if alliance == "blue" {
		source, owner = fuel.Blue, fuel.BlueSource
	}
	state.Source = owner
	if state.Source == "" {
		state.Source = "plc"
	}
	if fuel.Started {
		state.Flags |= fueldevice.FlagStarted
	}
	if fuel.Phase == int(AutoPeriod) {
		state.Flags |= fueldevice.FlagAuto
	}
	if fuel.Phase == int(TeleopPeriod) {
		state.Flags |= fueldevice.FlagTeleop
	}
	if fuel.Aborted {
		state.Flags |= fueldevice.FlagAborted
	}
	if fuel.Phase == int(PostMatch) {
		state.Flags |= fueldevice.FlagComplete
	}
	if fuel.Mode == model.FuelManual {
		state.Flags |= fueldevice.FlagManual
	}
	if fuel.Recovered || !fuel.Started || fuel.Status != "active" {
		return state
	}
	hub := source.Raw
	physical, _, _, ok := hub.GetCurrentShiftTiming(fuel.ClockStart, now)
	probe := game.Hub{WonAuto: hub.WonAuto}
	if ok {
		probe.ShiftCounts[physical] = 1
		if probe.GetShiftCount(physical, true) > 0 && !fuel.Aborted {
			state.Flags |= fueldevice.FlagHubActive
		}
	}
	shift, ok := hub.ScoringShiftAt(fuel.ClockStart, now)
	if fuel.Aborted && !now.Before(fuel.StoppedAt.Add(time.Duration(game.ScoringGracePeriodSec)*time.Second)) {
		ok = false
	}
	if ok {
		state.Shift = int(shift)
		probe = game.Hub{WonAuto: hub.WonAuto}
		probe.ShiftCounts[shift] = 1
		if probe.GetShiftCount(shift, true) > 0 && owner == "arduino" && fuel.Mode != model.FuelManual && arena.fuelStorageError == "" {
			state.Flags |= fueldevice.FlagScoring
		}
	}
	return state
}

// ApplyFuelDeviceEvent uses the same durability/serialization boundary as manual
// actions. Sensor events update the raw source only; adjustments are immutable.
func (arena *Arena) ApplyFuelDeviceEvent(event fueldevice.Event) fueldevice.Outcome {
	arena.fuelOperationMutex.Lock()
	defer arena.fuelOperationMutex.Unlock()
	fail := func(message string, storage bool) fueldevice.Outcome {
		return fueldevice.Outcome{Error: message, StorageError: storage}
	}
	if event.Alliance != "red" && event.Alliance != "blue" || event.Boot == 0 || event.Sequence == 0 {
		return fail("invalid sensor event identity", false)
	}
	id := uuid.NewSHA1(uuid.NameSpaceOID, []byte(fmt.Sprintf("arduino:%s:%d:%d:%s", event.Device, event.Boot, event.Sequence, event.Attempt))).String()
	arena.scoreMutex.RLock()
	uncertain := arena.fuelPersistenceUncertain
	arena.scoreMutex.RUnlock()
	if uncertain {
		return fail("FUEL persistence uncertain; restart to reconcile", true)
	}
	previous, err := arena.Database.GetFuelOperation(id)
	if err != nil {
		return fail(err.Error(), true)
	}
	if previous != nil {
		return fueldevice.Outcome{Accepted: true, Duplicate: true}
	}
	arena.scoreMutex.RLock()
	candidate := arena.checkpointLocked()
	current := arena.fuelDeviceStateLocked(event.Alliance, time.Now())
	device := arena.fuelDeviceIds[event.Alliance]
	arena.scoreMutex.RUnlock()
	if event.Device != device || device == "" {
		return fail("device is not assigned to this alliance", false)
	}
	if event.Attempt != candidate.AttemptId || !event.ValidLease || current != event.Expected {
		return fail("stale match, phase or state lease", false)
	}
	if current.Flags&fueldevice.FlagScoring == 0 {
		return fail("scoring disabled, inactive HUB, or another source owns this alliance", false)
	}
	source := &candidate.Red
	if event.Alliance == "blue" {
		source = &candidate.Blue
	}
	shift := game.Shift(current.Shift)
	effective := effectiveFuel(*source, candidate.Mode)
	before := effective.GetShiftCount(game.ShiftAuto, true)
	if shift != game.ShiftAuto {
		before = effective.GetTeleopActiveFuelCount()
	}
	if before >= 2147483647 {
		return fail("FUEL overflow", false)
	}
	source.Raw.ShiftCounts[shift]++
	source.Suspended = ""
	candidate.Revision++
	candidate.Result.RedScore.Hub = effectiveFuel(candidate.Red, candidate.Mode)
	candidate.Result.BlueScore.Hub = effectiveFuel(candidate.Blue, candidate.Mode)
	metadata := candidate.FuelMetadata
	metadata.AttemptId = candidate.ResultAttemptId
	candidate.Result.FuelControl = &metadata
	phase := "teleop"
	if shift == game.ShiftAuto {
		phase = "auto"
	}
	op := &model.FuelOperation{DeviceId: event.Device, Boot: event.Boot, EventSequence: event.Sequence, ResultAttemptId: candidate.ResultAttemptId, Mode: candidate.Mode, Kind: "sensorFuel", Actor: "arduino:" + event.Device, Time: time.Now(), Before: before, After: before + 1, Revision: candidate.Revision, Command: model.FuelCommand{RequestId: id, AttemptId: event.Attempt, MatchId: candidate.Match.Id, Alliance: event.Alliance, Phase: phase, Delta: 1}}
	if err = arena.Database.SaveFuelCheckpoint(candidate, op); err != nil {
		arena.scoreMutex.Lock()
		arena.fuelPersistenceUncertain = true
		arena.fuelStorageError = err.Error()
		arena.scoreMutex.Unlock()
		return fail(err.Error(), true)
	}
	arena.scoreMutex.Lock()
	target := &arena.fuel.Red
	if event.Alliance == "blue" {
		target = &arena.fuel.Blue
	}
	// The selected source excludes PLC writes; the operation lock excludes all
	// other sensor/manual/result/lifecycle writers, so no raw update is lost.
	target.Raw.ShiftCounts[shift]++
	target.Suspended = ""
	arena.projectFuelLocked()
	arena.fuel.Revision++
	arena.fuel.UpdatedAt = candidate.UpdatedAt
	arena.scoreMutex.Unlock()
	arena.RealtimeScoreNotifier.Notify()
	arena.ManualFuelNotifier.Notify()
	return fueldevice.Outcome{Accepted: true}
}
