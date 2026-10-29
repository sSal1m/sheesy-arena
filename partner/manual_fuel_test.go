package partner

import (
	"github.com/Team254/cheesy-arena/game"
	"github.com/Team254/cheesy-arena/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestManualFuelCanonicalTbaBreakdown(t *testing.T) {
	result := model.NewMatchResult()
	result.RedScore.Hub.ShiftCounts = [game.ShiftCount]int{28, 12, 0, 0, 0, 0, 5, 2}
	result.FuelControl = &model.FuelMetadata{Mode: model.FuelPlcCorrection}
	result.FuelControl.Red.Raw.ShiftCounts[game.ShiftAuto] = 25
	result.FuelControl.Red.Adjustment[game.ShiftAuto] = 3
	breakdown, err := createTbaScoringBreakdown(&model.EventSettings{}, &model.Match{Type: model.Qualification}, result, "red")
	require.NoError(t, err)
	hub := breakdown["hubScore"].(map[string]any)
	assert.EqualValues(t, 28, hub["autoCount"])
	assert.EqualValues(t, 19, hub["teleopCount"])
	assert.EqualValues(t, 47, hub["totalCount"])
	assert.EqualValues(t, 47, breakdown["totalPoints"])
}
