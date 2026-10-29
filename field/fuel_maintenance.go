package field

import (
	"fmt"
	"github.com/Team254/cheesy-arena/model"
)

// BeginFuelMaintenance excludes checkpoint writes and operator commands while
// the web UI replaces/clears the event DB. Finish must be called on every path.
func (arena *Arena) BeginFuelMaintenance() (func(bool) error, error) {
	arena.fuelOperationMutex.Lock()
	arena.scoreMutex.Lock()
	if arena.fuel.Started && arena.fuel.Status != "committed" && arena.fuel.Status != "discarded" {
		arena.scoreMutex.Unlock()
		arena.fuelOperationMutex.Unlock()
		return nil, fmt.Errorf("post or discard the pending match before restoring or clearing the database")
	}
	previous := arena.fuel.Status
	arena.fuel.Status = "maintenance"
	arena.fuel.Revision++
	arena.scoreMutex.Unlock()
	return func(changed bool) error {
		defer arena.fuelOperationMutex.Unlock()
		if !changed {
			arena.scoreMutex.Lock()
			arena.fuel.Status = previous
			arena.fuel.Revision++
			arena.scoreMutex.Unlock()
			return nil
		}
		checkpoint, err := arena.Database.GetFuelCheckpoint()
		if err != nil {
			return err
		}
		arena.scoreMutex.Lock()
		if checkpoint != nil {
			arena.fuel.Mode = checkpoint.Mode
		}
		arena.fuelStorageError, arena.fuelPersistenceUncertain = "", false
		arena.scoreMutex.Unlock()
		arena.MatchState = PreMatch
		if err = arena.loadMatch(&model.Match{Type: model.Test, ShortName: "T", LongName: "Test Match"}); err != nil {
			return err
		}
		if err = arena.CreatePlayoffTournament(); err != nil {
			return err
		}
		if err = arena.UpdatePlayoffTournament(); err != nil {
			return err
		}
		if checkpoint != nil && checkpoint.Started && checkpoint.Status != "committed" && checkpoint.Status != "discarded" {
			err = arena.recoverFuel(checkpoint)
		} else {
			err = arena.persistFuelLocked()
		}
		arena.MatchLoadNotifier.Notify()
		arena.RealtimeScoreNotifier.Notify()
		arena.ManualFuelNotifier.Notify()
		return err
	}, nil
}
