// Copyright 2014 Team 254. All Rights Reserved.
// Author: pat@patfairbank.com (Patrick Fairbank)
//
// Web handlers for the referee interface.

package web

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Team254/cheesy-arena/field"
	"github.com/Team254/cheesy-arena/game"
	"github.com/Team254/cheesy-arena/model"
	"github.com/Team254/cheesy-arena/websocket"
	"github.com/mitchellh/mapstructure"
)

// Renders the referee interface for assigning fouls.
func (web *Web) refereePanelHandler(w http.ResponseWriter, r *http.Request) {
	if !web.userIsAdmin(w, r) {
		return
	}

	template, err := web.parseFiles("templates/referee_panel.html", "templates/base.html")
	if err != nil {
		handleWebErr(w, err)
		return
	}

	data := struct {
		*model.EventSettings
	}{web.arena.EventSettings}
	err = template.ExecuteTemplate(w, "base_no_navbar", data)
	if err != nil {
		handleWebErr(w, err)
		return
	}
}

// Renders a partial template for when the foul list is updated.
func (web *Web) refereePanelFoulListHandler(w http.ResponseWriter, r *http.Request) {
	template, err := web.parseFiles("templates/referee_panel_foul_list.html")
	if err != nil {
		handleWebErr(w, err)
		return
	}

	result := web.arena.ScoreResultSnapshot()
	data := struct {
		Match     *model.Match
		RedFouls  []game.Foul
		BlueFouls []game.Foul
		Rules     map[int]*game.Rule
	}{
		web.arena.CurrentMatch,
		result.RedScore.Fouls,
		result.BlueScore.Fouls,
		game.GetAllRules(),
	}
	err = template.ExecuteTemplate(w, "referee_panel_foul_list", data)
	if err != nil {
		handleWebErr(w, err)
		return
	}
}

// The websocket endpoint for the refereee interface client to send control commands and receive status updates.
func (web *Web) refereePanelWebsocketHandler(w http.ResponseWriter, r *http.Request) {
	if !web.userIsAdmin(w, r) {
		return
	}

	ws, err := websocket.NewWebsocket(w, r)
	if err != nil {
		handleWebErr(w, err)
		return
	}
	defer closeWebsocket(ws)

	// Subscribe the websocket to the notifiers whose messages will be passed on to the client, in a separate goroutine.
	go ws.HandleNotifiers(
		web.arena.MatchLoadNotifier,
		web.arena.MatchTimeNotifier,
		web.arena.RealtimeScoreNotifier,
		web.arena.ScoringStatusNotifier,
		web.arena.ReloadDisplaysNotifier,
		web.arena.ArenaStatusNotifier,
	)

	// Loop, waiting for commands and responding to them, until the client closes the connection.
	for {
		messageType, data, err := ws.Read()
		if err != nil {
			if err == io.EOF {
				// Client has closed the connection; nothing to do here.
				return
			}
			log.Println(err)
			return
		}

		switch messageType {
		case "addFoul":
			args := struct {
				Alliance string
				IsMajor  bool
			}{}
			err = mapstructure.Decode(data, &args)
			if err != nil {
				writeWebsocketError(ws, err.Error())
				continue
			}

			if err := web.arena.AddLiveFoul(args.Alliance, args.IsMajor); err != nil {
				writeWebsocketError(ws, err.Error())
			}
		case "toggleFoulType", "updateFoulTeam", "updateFoulRule", "deleteFoul":
			args := struct {
				Alliance string
				Index    int
				TeamId   int
				RuleId   int
			}{}
			err = mapstructure.Decode(data, &args)
			if err != nil {
				writeWebsocketError(ws, err.Error())
				continue
			}

			if err := web.arena.EditLiveFoul(args.Alliance, messageType, args.Index, args.TeamId, args.RuleId); err != nil {
				writeWebsocketError(ws, err.Error())
			}
		case "card":
			args := struct {
				Alliance string
				TeamId   int
				Card     string
			}{}
			err = mapstructure.Decode(data, &args)
			if err != nil {
				writeWebsocketError(ws, err.Error())
				continue
			}

			if err := web.arena.SetLiveCard(args.Alliance, args.TeamId, args.Card); err != nil {
				writeWebsocketError(ws, err.Error())
			}
		case "toggleBypass":
			station, ok := data.(string)
			if !ok {
				writeWebsocketError(ws, fmt.Sprintf("Failed to parse '%s' message.", messageType))
				continue
			}
			err = web.arena.ToggleBypass(station)
			if err != nil {
				writeWebsocketError(ws, err.Error())
				continue
			}
		case "signalVolunteers":
			web.arena.SignalVolunteers()
		case "signalReset":
			web.arena.SignalReset()
		case "commitAndPost":
			if web.arena.MatchState != field.PostMatch {
				// Don't allow committing the fouls until the match is over.
				continue
			}
			web.arena.SetFoulsCommitted()

			err = web.commitPostAndLoadNextMatch()
			if err != nil {
				writeWebsocketError(ws, err.Error())
				continue
			}
		default:
			writeWebsocketError(ws, fmt.Sprintf("Invalid message type '%s'.", messageType))
		}
	}
}
