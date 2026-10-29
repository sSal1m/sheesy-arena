package model

import (
	"encoding/json"
	"github.com/Team254/cheesy-arena/game"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
)

func TestFuelAtomicCheckpointAndOperation(t *testing.T) {
	database, err := OpenDatabase(filepath.Join(t.TempDir(), "event.db"))
	require.NoError(t, err)
	defer database.Close()
	checkpoint := &FuelCheckpoint{FuelMetadata: FuelMetadata{AttemptId: "a", Mode: FuelManual}, Revision: 1}
	checkpoint.Red.Adjustment[game.ShiftAuto] = 5
	operation := &FuelOperation{Command: FuelCommand{RequestId: "request", AttemptId: "a"}, After: 5}
	require.NoError(t, database.SaveFuelCheckpoint(checkpoint, operation))
	checkpoint.Red.Adjustment[game.ShiftAuto] = 99
	require.Error(t, database.SaveFuelCheckpoint(checkpoint, operation))
	stored, err := database.GetFuelCheckpoint()
	require.NoError(t, err)
	assert.Equal(t, 5, stored.Red.Adjustment[game.ShiftAuto])
	history, err := database.GetFuelHistory("a")
	require.NoError(t, err)
	assert.Len(t, history, 1)
}

func TestFuelLegacyResultAndMetadataReconciliation(t *testing.T) {
	var result MatchResult
	require.NoError(t, json.Unmarshal([]byte(`{"MatchId":1,"RedScore":{"Hub":{"ShiftCounts":[25,5]}},"BlueScore":{}}`), &result))
	assert.Nil(t, result.FuelControl)
	require.NoError(t, ReconcileFuelMetadata(&result, nil))
	metadata := &FuelMetadata{Mode: FuelPlcCorrection}
	metadata.Red.Raw.ShiftCounts[game.ShiftAuto] = 22
	require.NoError(t, ReconcileFuelMetadata(&result, metadata))
	assert.Equal(t, 3, result.FuelControl.Red.Adjustment[game.ShiftAuto])
	assert.Equal(t, 25, result.RedScoreSummary().AutoFuelPoints)
	result.RedScore.Hub.ShiftCounts[game.ShiftAuto] = -1
	require.Error(t, ReconcileFuelMetadata(&result, metadata))
}

func TestMatchResultCloneOwnsMutableScoreData(t *testing.T) {
	result := NewMatchResult()
	result.RedScore.Fouls = []game.Foul{{FoulId: 1}}
	result.RedCards["254"] = "yellow"
	result.FuelControl = &FuelMetadata{Mode: FuelManual}
	cloned := result.Clone()
	assert.NotNil(t, cloned.BlueScore.Fouls) // Existing editors append to the JSON array.
	result.RedScore.Fouls[0].FoulId = 2
	result.RedCards["254"] = "red"
	result.FuelControl.Mode = FuelPlcCorrection
	assert.Equal(t, 1, cloned.RedScore.Fouls[0].FoulId)
	assert.Equal(t, "yellow", cloned.RedCards["254"])
	assert.Equal(t, FuelManual, cloned.FuelControl.Mode)
}
