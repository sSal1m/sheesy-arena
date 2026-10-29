package field

import (
	"fmt"
	"github.com/Team254/cheesy-arena/game"
	"github.com/Team254/cheesy-arena/model"
	"strconv"
)

func (arena *Arena) AddLiveFoul(alliance string, major bool) error {
	return arena.EditLiveScore(func(red, blue *RealtimeScore) error {
		score := blue
		if alliance == "red" {
			score = red
		} else if alliance != "blue" {
			return fmt.Errorf("invalid alliance")
		}
		score.CurrentScore.Fouls = append(score.CurrentScore.Fouls, game.Foul{FoulId: arena.NextFoulId, IsMajor: major})
		arena.NextFoulId++
		return nil
	})
}

func (arena *Arena) EditLiveFoul(alliance, action string, index, teamId, ruleId int) error {
	return arena.EditLiveScore(func(red, blue *RealtimeScore) error {
		score := blue
		if alliance == "red" {
			score = red
		} else if alliance != "blue" {
			return fmt.Errorf("invalid alliance")
		}
		fouls := &score.CurrentScore.Fouls
		if index < 0 || index >= len(*fouls) {
			return errNoScoreChange
		}
		switch action {
		case "toggleFoulType":
			(*fouls)[index].IsMajor = !(*fouls)[index].IsMajor
			(*fouls)[index].RuleId = 0
		case "deleteFoul":
			*fouls = append((*fouls)[:index], (*fouls)[index+1:]...)
		case "updateFoulTeam":
			if (*fouls)[index].TeamId == teamId {
				(*fouls)[index].TeamId = 0
			} else {
				(*fouls)[index].TeamId = teamId
			}
		case "updateFoulRule":
			(*fouls)[index].RuleId = ruleId
		default:
			return fmt.Errorf("invalid foul action")
		}
		return nil
	})
}

func (arena *Arena) SetLiveCard(alliance string, teamId int, card string) error {
	return arena.EditLiveScore(func(red, blue *RealtimeScore) error {
		score := blue
		if alliance == "red" {
			score = red
		} else if alliance != "blue" {
			return fmt.Errorf("invalid alliance")
		}
		if arena.CurrentMatch.Type == model.Playoff {
			score.PlayoffAllianceCard = card
			score.Cards = make(map[string]string)
			score.CurrentScore.PlayoffDq = card == "red" || card == "dq"
		} else {
			score.Cards[strconv.Itoa(teamId)] = card
		}
		return nil
	})
}

func (arena *Arena) SetFoulsCommitted() {
	arena.fuelOperationMutex.Lock()
	defer arena.fuelOperationMutex.Unlock()
	arena.scoreMutex.Lock()
	arena.RedRealtimeScore.FoulsCommitted, arena.BlueRealtimeScore.FoulsCommitted = true, true
	arena.scoreMutex.Unlock()
	arena.ScoringStatusNotifier.Notify()
}

// ReplaceLiveResult is the existing full-result editor's boundary. Unlike delta
// commands, a full replacement must reject a stale revision to avoid lost edits.
func (arena *Arena) ReplaceLiveResult(result *model.MatchResult) error {
	arena.fuelOperationMutex.Lock()
	defer arena.fuelOperationMutex.Unlock()
	arena.scoreMutex.RLock()
	if arena.fuelPersistenceUncertain {
		arena.scoreMutex.RUnlock()
		return fmt.Errorf("FUEL persistence is uncertain; restart to reconcile the durable state")
	}
	candidate := arena.checkpointLocked()
	arena.scoreMutex.RUnlock()
	if candidate.Status == "committed" || candidate.Status == "discarded" {
		return fmt.Errorf("result is closed; edit the recorded match instead")
	}
	arena.scoreMutex.RLock()
	gracePending := arena.fuelGracePendingLocked()
	arena.scoreMutex.RUnlock()
	if arena.MatchState == PostMatch && gracePending {
		return fmt.Errorf("wait for the FUEL scoring grace period before replacing the result")
	}
	if arena.MatchState != PreMatch && arena.MatchState != PostMatch {
		return fmt.Errorf("use the scoring panels during a running match; full result editing is only available before/after play")
	}
	if result.MatchId != candidate.Match.Id {
		return fmt.Errorf("stale match")
	}
	if result.RedScore == nil || result.BlueScore == nil {
		return fmt.Errorf("missing score")
	}
	if result.FuelContext != candidate.AttemptId || result.FuelRevision != candidate.Revision {
		// Preserve legacy, pre-match editor clients which have no FUEL context.
		if candidate.Started || result.FuelContext != "" {
			return fmt.Errorf("score changed; reload the result editor")
		}
	}
	if candidate.AutoWinnerFrozen && (result.RedScore.Hub.WonAuto != candidate.Red.Raw.WonAuto || result.BlueScore.Hub.WonAuto != candidate.Blue.Raw.WonAuto) {
		return fmt.Errorf("auto winner is frozen for the current match")
	}
	metadata := candidate.FuelMetadata
	if err := model.ReconcileFuelMetadata(result, &metadata); err != nil {
		return err
	}
	candidate.Red, candidate.Blue = metadata.Red, metadata.Blue
	candidate.Result = *result
	candidate.Revision++
	if err := arena.Database.SaveFuelCheckpoint(candidate, nil); err != nil {
		arena.scoreMutex.Lock()
		arena.fuelPersistenceUncertain = true
		arena.fuelStorageError = err.Error()
		arena.scoreMutex.Unlock()
		return err
	}
	arena.scoreMutex.Lock()
	arena.fuel.Red, arena.fuel.Blue = candidate.Red, candidate.Blue
	arena.fuel.Revision++
	arena.RedRealtimeScore.CurrentScore = cloneGameScore(*result.RedScore)
	arena.BlueRealtimeScore.CurrentScore = cloneGameScore(*result.BlueScore)
	arena.RedRealtimeScore.Cards = cloneCards(result.RedCards)
	arena.BlueRealtimeScore.Cards = cloneCards(result.BlueCards)
	arena.RedRealtimeScore.PlayoffAllianceCard, arena.BlueRealtimeScore.PlayoffAllianceCard = result.PlayoffRedAllianceCard, result.PlayoffBlueAllianceCard
	arena.scoreMutex.Unlock()
	arena.RealtimeScoreNotifier.Notify()
	arena.ManualFuelNotifier.Notify()
	return nil
}
