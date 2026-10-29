package field

import (
	"github.com/Team254/cheesy-arena/game"
	"github.com/Team254/cheesy-arena/model"
	"github.com/Team254/cheesy-arena/plc"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func fuelCommand(arena *Arena, alliance, phase string, delta int) model.FuelCommand {
	state := arena.FuelState()
	return model.FuelCommand{RequestId: uuid.NewString(), AttemptId: state.AttemptId, MatchId: state.MatchId, Alliance: alliance, Phase: phase, Delta: delta, Operator: "Unit operator"}
}

func setFuelMode(t *testing.T, arena *Arena, mode string) {
	command := fuelCommand(arena, "", "", 0)
	command.Mode = mode
	require.Empty(t, arena.ApplyFuelCommand("setFuelMode", command, "test-admin").Error)
}

func fuelPhase(arena *Arena, state MatchState) {
	arena.MatchState = state
	arena.MatchStartTime = time.Now()
	if state == TeleopPeriod {
		arena.MatchStartTime = time.Now().Add(-game.GetDurationToTeleopStart() - time.Second)
	}
	arena.syncFuelPhase()
}

func TestManualFuelFourCounters(t *testing.T) {
	for _, alliance := range []string{"red", "blue"} {
		for _, phase := range []string{"auto", "teleop"} {
			t.Run(alliance+"/"+phase, func(t *testing.T) {
				arena := setupTestArena(t)
				setFuelMode(t, arena, model.FuelManual)
				fuelPhase(arena, AutoPeriod)
				if phase == "teleop" {
					fuelPhase(arena, TeleopPeriod)
				}
				require.Empty(t, arena.ApplyFuelCommand("adjustFuel", fuelCommand(arena, alliance, phase, 1), "a").Error)
				result := arena.ScoreResultSnapshot()
				score := result.RedScoreSummary()
				if alliance == "blue" {
					score = result.BlueScoreSummary()
				}
				assert.Equal(t, 1, score.Score)
				if phase == "auto" {
					assert.Equal(t, 1, score.AutoFuelPoints)
				} else {
					assert.Equal(t, 1, score.TeleopFuelPoints)
				}
				require.Empty(t, arena.ApplyFuelCommand("adjustFuel", fuelCommand(arena, alliance, phase, -1), "b").Error)
				require.NotEmpty(t, arena.ApplyFuelCommand("adjustFuel", fuelCommand(arena, alliance, phase, -1), "b").Error)
				assert.Zero(t, arena.ScoreResultSnapshot().RedScoreSummary().Score)
				assert.Zero(t, arena.ScoreResultSnapshot().BlueScoreSummary().Score)
			})
		}
	}
}

func TestManualFuelPlcCorrectionIndependentCursor(t *testing.T) {
	arena := setupTestArena(t)
	fuelPhase(arena, AutoPeriod)
	sample := plc.FuelSnapshot{Sequence: 1, Connection: 1, Healthy: true, ReadAt: time.Now()}
	arena.processFuelSnapshot(sample, sample.ReadAt)
	sample.Sequence++
	sample.Red = 25
	sample.Blue = 10
	sample.ReadAt = time.Now()
	arena.processFuelSnapshot(sample, sample.ReadAt)
	for i := 0; i < 3; i++ {
		require.Empty(t, arena.ApplyFuelCommand("adjustFuel", fuelCommand(arena, "red", "auto", 1), "a").Error)
	}
	assert.Equal(t, 28, arena.FuelState().Red.Auto)
	sample.Sequence++
	sample.Red = 27
	sample.ReadAt = time.Now()
	arena.processFuelSnapshot(sample, sample.ReadAt)
	assert.Equal(t, 30, arena.FuelState().Red.Auto)
	assert.Equal(t, 27, arena.FuelState().Red.RawAuto)
	assert.Equal(t, 3, arena.FuelState().Red.AutoAdjustment)
	assert.Equal(t, 10, arena.FuelState().Blue.Auto)
	require.Empty(t, arena.ApplyFuelCommand("adjustFuel", fuelCommand(arena, "red", "auto", -5), "a").Error)
	sample.Sequence++
	sample.Red = 28
	sample.ReadAt = time.Now()
	arena.processFuelSnapshot(sample, sample.ReadAt)
	assert.Equal(t, 26, arena.FuelState().Red.Auto)
	// Repeated and out-of-order hardware samples do not score again.
	arena.processFuelSnapshot(sample, sample.ReadAt)
	sample.Sequence--
	sample.Red = 99
	arena.processFuelSnapshot(sample, sample.ReadAt)
	assert.Equal(t, 26, arena.FuelState().Red.Auto)
}

func TestManualFuelIgnoresOnlyPlcFuel(t *testing.T) {
	arena := setupTestArena(t)
	setFuelMode(t, arena, model.FuelManual)
	fake := &FakePlc{isEnabled: true, ftaReady: true, redHubCount: 99}
	arena.Plc = fake
	fuelPhase(arena, AutoPeriod)
	require.Empty(t, arena.ApplyFuelCommand("adjustFuel", fuelCommand(arena, "red", "auto", 1), "a").Error)
	arena.handlePlcInputOutput()
	assert.Equal(t, 1, arena.FuelState().Red.Auto)
	assert.True(t, fake.redHubMotor)
	fake.redEStops[0] = true
	arena.handlePlcInputOutput()
	assert.True(t, arena.AllianceStations["R1"].EStop)
}

func TestManualFuelSourceFaultsRequireRebase(t *testing.T) {
	for _, fault := range []string{"decrease", "reconnect", "reset", "stale", "unhealthy"} {
		t.Run(fault, func(t *testing.T) {
			arena := setupTestArena(t)
			fuelPhase(arena, AutoPeriod)
			sample := plc.FuelSnapshot{Sequence: 1, Connection: 1, Healthy: true, ReadAt: time.Now()}
			arena.processFuelSnapshot(sample, sample.ReadAt)
			sample.Sequence++
			sample.Red = 25
			sample.ReadAt = time.Now()
			arena.processFuelSnapshot(sample, sample.ReadAt)
			require.Empty(t, arena.ApplyFuelCommand("adjustFuel", fuelCommand(arena, "red", "auto", 1), "a").Error)
			sample.Sequence++
			sample.ReadAt = time.Now()
			switch fault {
			case "decrease":
				sample.Red = 2
			case "reconnect":
				sample.Connection++
			case "reset":
				sample.Reset++
			case "stale":
				sample.ReadAt = time.Now().Add(-time.Second)
			case "unhealthy":
				sample.Healthy = false
			}
			arena.processFuelSnapshot(sample, time.Now())
			assert.NotEmpty(t, arena.FuelState().Red.Suspended)
			assert.Equal(t, 26, arena.FuelState().Red.Auto)
			// Healthy data alone must not silently resume a suspended source.
			sample.Sequence++
			sample.ReadAt = time.Now()
			sample.Healthy = true
			sample.Red += 2
			arena.processFuelSnapshot(sample, time.Now())
			assert.Equal(t, 26, arena.FuelState().Red.Auto)
			fake := &FakePlc{fuelSnapshot: &sample}
			arena.Plc = fake
			require.Empty(t, arena.ApplyFuelCommand("rebasePlcFuel", fuelCommand(arena, "red", "", 0), "a").Error)
			assert.Equal(t, 26, arena.FuelState().Red.Auto)
			sample.Sequence++
			sample.ReadAt = time.Now()
			sample.Red += 2
			arena.processFuelSnapshot(sample, sample.ReadAt)
			assert.Equal(t, 28, arena.FuelState().Red.Auto)
		})
	}
}

func TestManualFuelDoesNotAcceptOldNonzeroAtNewMatch(t *testing.T) {
	arena := setupTestArena(t)
	fuelPhase(arena, AutoPeriod)
	sample := plc.FuelSnapshot{Sequence: 1, Connection: 1, Red: 90, Healthy: true, ReadAt: time.Now()}
	arena.processFuelSnapshot(sample, sample.ReadAt)
	assert.Zero(t, arena.FuelState().Red.Auto)
	assert.NotEmpty(t, arena.FuelState().Red.Suspended)
	sample.Sequence++
	sample.Red = 0
	arena.processFuelSnapshot(sample, sample.ReadAt)
	sample.Sequence++
	sample.Red = 3
	arena.processFuelSnapshot(sample, sample.ReadAt)
	assert.Equal(t, 3, arena.FuelState().Red.Auto)
}

func TestManualFuelPhaseRulesAndFrozenAutoWinner(t *testing.T) {
	arena := setupTestArena(t)
	setFuelMode(t, arena, model.FuelManual)
	require.NotEmpty(t, arena.ApplyFuelCommand("adjustFuel", fuelCommand(arena, "red", "auto", 1), "a").Error)
	fuelPhase(arena, AutoPeriod)
	require.NotEmpty(t, arena.ApplyFuelCommand("adjustFuel", fuelCommand(arena, "red", "teleop", 1), "a").Error)
	require.Empty(t, arena.ApplyFuelCommand("adjustFuel", fuelCommand(arena, "red", "auto", 5), "a").Error)
	fuelPhase(arena, PausePeriod)
	require.Empty(t, arena.ApplyFuelCommand("adjustFuel", fuelCommand(arena, "red", "auto", 1), "a").Error)
	arena.handleAutoWinner()
	require.True(t, arena.RedRealtimeScore.CurrentScore.Hub.WonAuto)
	fuelPhase(arena, TeleopPeriod)
	// The red Hub is inactive in Shift1. Manual scored FUEL must still count.
	arena.MatchStartTime = time.Now().Add(-game.GetDurationToTeleopStart() - 15*time.Second)
	require.Empty(t, arena.ApplyFuelCommand("adjustFuel", fuelCommand(arena, "red", "teleop", 5), "a").Error)
	assert.Equal(t, 5, arena.FuelState().Red.Teleop)
	assert.Equal(t, 5, arena.ScoreResultSnapshot().RedScore.Hub.ShiftCounts[game.ShiftTransition])
	for i := 0; i < 2; i++ {
		require.Empty(t, arena.ApplyFuelCommand("adjustFuel", fuelCommand(arena, "blue", "auto", 5), "a").Error)
	}
	assert.True(t, arena.ScoreResultSnapshot().RedScore.Hub.WonAuto)
	assert.Equal(t, "R", arena.AllianceStations["R1"].GameData)
	command := fuelCommand(arena, "", "", 0)
	command.Mode = model.FuelPlcCorrection
	require.NotEmpty(t, arena.ApplyFuelCommand("setFuelMode", command, "a").Error)
	fuelPhase(arena, PostMatch)
	require.Empty(t, arena.ApplyFuelCommand("adjustFuel", fuelCommand(arena, "red", "teleop", 1), "a").Error)
	assert.Equal(t, 1, arena.ScoreResultSnapshot().RedScore.Hub.ShiftCounts[game.ShiftEndgame])
}

func TestManualFuelIdempotencyAndConcurrentOperators(t *testing.T) {
	arena := setupTestArena(t)
	setFuelMode(t, arena, model.FuelManual)
	fuelPhase(arena, AutoPeriod)
	command := fuelCommand(arena, "red", "auto", 1)
	require.Empty(t, arena.ApplyFuelCommand("adjustFuel", command, "a").Error)
	assert.True(t, arena.ApplyFuelCommand("adjustFuel", command, "a").Duplicate)
	command.Delta = 5
	require.NotEmpty(t, arena.ApplyFuelCommand("adjustFuel", command, "a").Error)
	var group sync.WaitGroup
	for i := 0; i < 20; i++ {
		command := fuelCommand(arena, "red", "auto", 1)
		group.Add(1)
		go func() {
			defer group.Done()
			assert.Empty(t, arena.ApplyFuelCommand("adjustFuel", command, "operator").Error)
		}()
	}
	group.Add(1)
	go func() {
		defer group.Done()
		for i := 0; i < 100; i++ {
			arena.ScoreResultSnapshot()
			arena.generateRealtimeScoreMessage()
			arena.FuelState()
		}
	}()
	group.Wait()
	assert.Equal(t, 21, arena.FuelState().Red.Auto)
	history, err := arena.Database.GetFuelHistory(command.AttemptId)
	require.NoError(t, err)
	assert.Len(t, history, 22) // Includes mode selection; duplicate request is not a new action.
}

func TestManualFuelCommitBarrierAndNewAttempt(t *testing.T) {
	arena := setupTestArena(t)
	setFuelMode(t, arena, model.FuelManual)
	fuelPhase(arena, AutoPeriod)
	require.Empty(t, arena.ApplyFuelCommand("adjustFuel", fuelCommand(arena, "red", "auto", 1), "a").Error)
	fuelPhase(arena, PostMatch)
	stale := fuelCommand(arena, "red", "auto", 1)
	result, finish, err := arena.BeginFuelCommit()
	require.NoError(t, err)
	assert.Equal(t, 1, result.RedScoreSummary().Score)
	finished := make(chan FuelReply, 1)
	go func() { finished <- arena.ApplyFuelCommand("adjustFuel", stale, "a") }()
	finish(true)
	assert.NotEmpty(t, (<-finished).Error)
	require.NoError(t, arena.ResetMatch())
	require.NoError(t, arena.LoadTestMatch())
	assert.NotEqual(t, stale.AttemptId, arena.FuelState().AttemptId)
	assert.Equal(t, stale.MatchId, arena.FuelState().MatchId) // Test matches share ID 0.
	assert.NotEmpty(t, arena.ApplyFuelCommand("adjustFuel", stale, "a").Error)
	assert.Zero(t, arena.FuelState().Red.Auto)
	assert.Equal(t, 1, result.RedScoreSummary().Score) // Posted snapshot is independent.
}

func TestManualFuelRecoveryAndDedupAcrossRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "event.db")
	arena, err := NewArena(path)
	require.NoError(t, err)
	setFuelMode(t, arena, model.FuelManual)
	fuelPhase(arena, AutoPeriod)
	command := fuelCommand(arena, "red", "auto", 5)
	require.Empty(t, arena.ApplyFuelCommand("adjustFuel", command, "a").Error)
	require.NoError(t, arena.Database.Close())
	recovered, err := NewArena(path)
	require.NoError(t, err)
	defer recovered.Database.Close()
	assert.Equal(t, PostMatch, recovered.MatchState)
	assert.True(t, recovered.FuelState().Recovered)
	assert.Equal(t, 5, recovered.FuelState().Red.Auto)
	assert.NotEqual(t, command.AttemptId, recovered.FuelState().AttemptId)
	assert.True(t, recovered.MatchStartTime.IsZero())
	assert.True(t, recovered.ApplyFuelCommand("adjustFuel", command, "a").Duplicate)
	command.RequestId = uuid.NewString()
	assert.NotEmpty(t, recovered.ApplyFuelCommand("adjustFuel", command, "a").Error)
	require.Empty(t, recovered.ApplyFuelCommand("adjustFuel", fuelCommand(recovered, "red", "auto", -1), "a").Error)
	assert.Equal(t, 4, recovered.FuelState().Red.Auto)
}

func TestManualFuelClosedDatabaseDoesNotApply(t *testing.T) {
	arena := setupTestArena(t)
	setFuelMode(t, arena, model.FuelManual)
	fuelPhase(arena, AutoPeriod)
	require.NoError(t, arena.Database.Close())
	assert.NotEmpty(t, arena.ApplyFuelCommand("adjustFuel", fuelCommand(arena, "red", "auto", 1), "a").Error)
	assert.Zero(t, arena.FuelState().Red.Auto)
}

func TestManualFuelExistingBonusesAndG206(t *testing.T) {
	arena := setupTestArena(t)
	setFuelMode(t, arena, model.FuelManual)
	fuelPhase(arena, AutoPeriod)
	for i := 0; i < 20; i++ {
		require.Empty(t, arena.ApplyFuelCommand("adjustFuel", fuelCommand(arena, "red", "auto", 5), "a").Error)
	}
	assert.True(t, arena.ScoreResultSnapshot().RedScoreSummary().EnergizedBonusRankingPoint)
	var ruleId int
	for id, rule := range game.GetAllRules() {
		if rule.RuleNumber == "G206" {
			ruleId = id
			break
		}
	}
	require.NotZero(t, ruleId)
	require.NoError(t, arena.EditLiveScore(func(red, blue *RealtimeScore) error {
		red.CurrentScore.Fouls = []game.Foul{{RuleId: ruleId}}
		return nil
	}))
	assert.False(t, arena.ScoreResultSnapshot().RedScoreSummary().EnergizedBonusRankingPoint)
}

func TestManualFuelPlcPhaseAndGraceAttribution(t *testing.T) {
	arena := setupTestArena(t)
	fuelPhase(arena, AutoPeriod)
	start := arena.MatchStartTime
	sample := plc.FuelSnapshot{Sequence: 1, Connection: 1, Healthy: true, ReadAt: start}
	arena.processFuelSnapshot(sample, start)
	sample.Sequence++
	sample.Red = 5
	sample.ReadAt = start.Add(game.GetDurationToAutoEnd() + time.Second/2)
	arena.processFuelSnapshot(sample, sample.ReadAt)
	assert.Equal(t, 5, arena.FuelState().Red.Auto)
	require.Empty(t, arena.ApplyFuelCommand("adjustFuel", fuelCommand(arena, "red", "auto", 1), "a").Error)
	arena.handleAutoWinner()
	arena.MatchState = TeleopPeriod
	arena.syncFuelPhase()
	sample.Sequence++
	sample.Red = 8
	sample.ReadAt = start.Add(game.GetDurationToTeleopStart() + 4*time.Second)
	arena.processFuelSnapshot(sample, sample.ReadAt)
	assert.Equal(t, 6, arena.FuelState().Red.Auto)
	assert.Equal(t, 3, arena.FuelState().Red.Teleop)
	sample.Sequence++
	sample.Red = 10
	sample.ReadAt = start.Add(game.GetDurationToTeleopEnd() + time.Second)
	arena.processFuelSnapshot(sample, sample.ReadAt)
	assert.Equal(t, 2, arena.ScoreResultSnapshot().RedScore.Hub.ShiftCounts[game.ShiftPostMatch])
	sample.Sequence++
	sample.Red = 12
	sample.ReadAt = start.Add(game.GetDurationToTeleopEnd() + 4*time.Second)
	arena.processFuelSnapshot(sample, sample.ReadAt)
	assert.Equal(t, 5, arena.FuelState().Red.Teleop)
}

func TestManualFuelPlcCommitWaitsForGrace(t *testing.T) {
	arena := setupTestArena(t)
	arena.Plc = &FakePlc{isEnabled: true}
	fuelPhase(arena, AutoPeriod)
	arena.MatchState = PostMatch
	arena.syncFuelPhase()
	arena.MatchStartTime = time.Now().Add(-game.GetDurationToTeleopEnd() - time.Second)
	_, _, err := arena.BeginFuelCommit()
	require.ErrorContains(t, err, "grace")
	arena.MatchStartTime = time.Now().Add(-game.GetDurationToTeleopEnd() - 4*time.Second)
	_, finish, err := arena.BeginFuelCommit()
	require.NoError(t, err)
	finish(false)
}

func TestManualFuelCrashAfterDurableWriteBeforeAck(t *testing.T) {
	path := filepath.Join(t.TempDir(), "event.db")
	arena, err := NewArena(path)
	require.NoError(t, err)
	setFuelMode(t, arena, model.FuelManual)
	fuelPhase(arena, AutoPeriod)
	command := fuelCommand(arena, "red", "auto", 5)
	checkpoint := arena.checkpointLocked()
	checkpoint.Red.Adjustment[game.ShiftAuto] = 5
	checkpoint.Result.RedScore.Hub = effectiveFuel(checkpoint.Red, checkpoint.Mode)
	checkpoint.Revision++
	require.NoError(t, arena.Database.SaveFuelCheckpoint(checkpoint, &model.FuelOperation{Kind: "adjustFuel", Command: command, Actor: "a", After: 5}))
	require.NoError(t, arena.Database.Close())
	recovered, err := NewArena(path)
	require.NoError(t, err)
	defer recovered.Database.Close()
	assert.Equal(t, 5, recovered.FuelState().Red.Auto)
	reply := recovered.ApplyFuelCommand("adjustFuel", command, "a")
	assert.True(t, reply.Duplicate)
	assert.Equal(t, 5, reply.State.Red.Auto)
}

func TestManualFuelModeSurvivesIdleRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "event.db")
	arena, err := NewArena(path)
	require.NoError(t, err)
	setFuelMode(t, arena, model.FuelManual)
	require.NoError(t, arena.Database.Close())
	reopened, err := NewArena(path)
	require.NoError(t, err)
	defer reopened.Database.Close()
	assert.Equal(t, model.FuelManual, reopened.FuelState().Mode)
	assert.Equal(t, PreMatch, reopened.MatchState)
	assert.False(t, reopened.FuelState().Recovered)
}

func TestManualFuelDatabaseMaintenanceBarrier(t *testing.T) {
	arena := setupTestArena(t)
	oldAttempt := arena.FuelState().AttemptId
	finish, err := arena.BeginFuelMaintenance()
	require.NoError(t, err)
	require.NoError(t, finish(true))
	assert.NotEqual(t, oldAttempt, arena.FuelState().AttemptId)
	fuelPhase(arena, AutoPeriod)
	_, err = arena.BeginFuelMaintenance()
	require.ErrorContains(t, err, "pending match")
}

func TestManualFuelDurabilityUncertaintyBlocksRetries(t *testing.T) {
	arena := setupTestArena(t)
	setFuelMode(t, arena, model.FuelManual)
	fuelPhase(arena, AutoPeriod)
	command := fuelCommand(arena, "red", "auto", 1)
	require.Empty(t, arena.ApplyFuelCommand("adjustFuel", command, "a").Error)
	arena.fuelPersistenceUncertain = true
	assert.Contains(t, arena.ApplyFuelCommand("adjustFuel", command, "a").Error, "uncertain")
	assert.ErrorContains(t, arena.ReplaceLiveResult(arena.ScoreResultSnapshot()), "uncertain")
	assert.Equal(t, 1, arena.FuelState().Red.Auto)
}

func TestManualFuelPartialResultCommitRecovery(t *testing.T) {
	path := filepath.Join(t.TempDir(), "event.db")
	arena, err := NewArena(path)
	require.NoError(t, err)
	setFuelMode(t, arena, model.FuelManual)
	fuelPhase(arena, AutoPeriod)
	require.Empty(t, arena.ApplyFuelCommand("adjustFuel", fuelCommand(arena, "red", "auto", 5), "a").Error)
	fuelPhase(arena, PostMatch)
	result, _, err := arena.BeginFuelCommit()
	require.NoError(t, err)
	result.PlayNumber = 1
	require.NoError(t, arena.Database.CreateMatchResult(result))
	// Simulate process death after the result write but before the commit ACK.
	require.NoError(t, arena.Database.SaveFuelCheckpoint(arena.checkpointLocked(), nil))
	require.NoError(t, arena.Database.Close())
	recovered, err := NewArena(path)
	require.NoError(t, err)
	defer recovered.Database.Close()
	retried, finish, err := recovered.BeginFuelCommit()
	require.NoError(t, err)
	assert.Equal(t, result.Id, retried.Id)
	assert.Equal(t, 1, retried.PlayNumber)
	assert.Equal(t, 5, retried.RedScoreSummary().Score)
	finish(true)
	assert.False(t, recovered.FuelState().CanAuto)
}

func TestManualFuelRejectsOverflowAndUnnecessaryRebase(t *testing.T) {
	arena := setupTestArena(t)
	fuelPhase(arena, AutoPeriod)
	sample := plc.FuelSnapshot{Sequence: 1, Connection: 1, Healthy: true, ReadAt: time.Now()}
	arena.processFuelSnapshot(sample, sample.ReadAt)
	assert.Contains(t, arena.ApplyFuelCommand("rebasePlcFuel", fuelCommand(arena, "red", "", 0), "a").Error, "not available")
	arena.fuel.Red.Adjustment[game.ShiftAuto] = 2147483647
	arena.projectFuelLocked()
	assert.NotEmpty(t, arena.ApplyFuelCommand("adjustFuel", fuelCommand(arena, "red", "auto", 1), "a").Error)
	assert.Equal(t, 2147483647, arena.FuelState().Red.Auto)
}

func TestManualFuelAbortStopsAutomaticIntakeAfterGrace(t *testing.T) {
	arena := setupTestArena(t)
	fuelPhase(arena, AutoPeriod)
	sample := plc.FuelSnapshot{Sequence: 1, Connection: 1, Healthy: true, ReadAt: time.Now()}
	arena.processFuelSnapshot(sample, sample.ReadAt)
	require.NoError(t, arena.AbortMatch())
	arena.syncFuelPhase()
	arena.matchStopTime = time.Now().Add(-4 * time.Second)
	sample.Sequence++
	sample.Red = 10
	sample.ReadAt = time.Now()
	arena.processFuelSnapshot(sample, sample.ReadAt)
	assert.Zero(t, arena.FuelState().Red.Auto)
	require.Empty(t, arena.ApplyFuelCommand("adjustFuel", fuelCommand(arena, "red", "teleop", 1), "a").Error)
	assert.Equal(t, 1, arena.ScoreResultSnapshot().RedScore.Hub.ShiftCounts[game.ShiftEndgame])
}
