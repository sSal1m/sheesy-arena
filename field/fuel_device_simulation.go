package field

import (
	"fmt"
	"github.com/Team254/cheesy-arena/game"
	"github.com/Team254/cheesy-arena/model"
	"time"
)

type FuelSimulation struct {
	Simulation              bool
	Red, Blue               FuelCounts
	RedSummary, BlueSummary *game.ScoreSummary
}

// This sandbox has no route to live scores, MatchResult or the event DB.
// Synthetic timestamps select fixed game phases; no second match timer runs.
func (arena *Arena) SimulateFuelDevice(alliance, phase string, count int, reset bool) (FuelSimulation, error) {
	arena.fuelSimulationMutex.Lock()
	defer arena.fuelSimulationMutex.Unlock()
	if alliance != "red" && alliance != "blue" || count < 0 || count > 5 {
		return FuelSimulation{}, fmt.Errorf("invalid simulation input")
	}
	start := time.Unix(1000, 0)
	var at time.Time
	switch phase {
	case "auto":
		at = start.Add(time.Second)
	case "transition":
		at = start.Add(game.GetDurationToTeleopStart() + time.Second)
	case "shift1":
		at = start.Add(game.GetDurationToTeleopStart() + time.Duration(game.MatchTiming.TransitionShiftDurationSec+game.ScoringGracePeriodSec+1)*time.Second)
	case "shift2":
		at = start.Add(game.GetDurationToTeleopStart() + time.Duration(game.MatchTiming.TransitionShiftDurationSec+game.MatchTiming.ShiftDurationSec+game.ScoringGracePeriodSec+1)*time.Second)
	case "endgame":
		at = start.Add(game.GetDurationToTeleopEnd() - time.Second)
	case "post-match":
		at = start.Add(game.GetDurationToTeleopEnd() + time.Second)
	default:
		return FuelSimulation{}, fmt.Errorf("invalid simulation phase")
	}
	if reset {
		arena.fuelSimulation = [2]model.FuelAlliance{}
	}
	arena.fuelSimulation[1].Raw.WonAuto = true
	source := &arena.fuelSimulation[0]
	if alliance == "blue" {
		source = &arena.fuelSimulation[1]
	}
	total := 0
	for _, v := range source.Raw.ShiftCounts {
		total += v
	}
	if total > 100000-count {
		return FuelSimulation{}, fmt.Errorf("reset the simulation before adding more events")
	}
	source.Raw.UpdateState(total+count, start, at)
	red, blue := game.Score{Hub: arena.fuelSimulation[0].Raw}, game.Score{Hub: arena.fuelSimulation[1].Raw}
	counts := func(hub game.Hub) FuelCounts {
		return FuelCounts{Source: "simulation", Auto: hub.GetShiftCount(game.ShiftAuto, true), Teleop: hub.GetTeleopActiveFuelCount()}
	}
	return FuelSimulation{Simulation: true, Red: counts(red.Hub), Blue: counts(blue.Hub), RedSummary: red.Summarize(&blue), BlueSummary: blue.Summarize(&red)}, nil
}
