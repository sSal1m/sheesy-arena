package web

import (
	"bytes"
	"encoding/json"
	"github.com/Team254/cheesy-arena/field"
	"github.com/Team254/cheesy-arena/model"
	"github.com/Team254/cheesy-arena/websocket"
	"io"
	"log"
	"net/http"
	"sort"
)

func (web *Web) manualFuelHandler(w http.ResponseWriter, r *http.Request) {
	if !web.userIsAdmin(w, r) {
		return
	}
	template, err := web.parseFiles("templates/manual_fuel.html", "templates/base.html")
	if err != nil {
		handleWebErr(w, err)
		return
	}
	data := struct{ *model.EventSettings }{web.arena.EventSettings}
	if err = template.ExecuteTemplate(w, "base", data); err != nil {
		handleWebErr(w, err)
	}
}

func (web *Web) fuelActor(r *http.Request) (string, bool) {
	if web.arena.EventSettings.AdminPassword == "" {
		return "open-access", true
	}
	session := web.getUserSessionFromCookie(r)
	if session != nil && session.Username == adminUser {
		return session.Username, true
	}
	return "", false
}

func (web *Web) manualFuelWebsocketHandler(w http.ResponseWriter, r *http.Request) {
	if !web.userIsAdmin(w, r) {
		return
	}
	ws, err := websocket.NewWebsocket(w, r)
	if err != nil {
		handleWebErr(w, err)
		return
	}
	defer closeWebsocket(ws)
	go ws.HandleNotifiers(web.arena.ManualFuelNotifier)
	for {
		kind, data, err := ws.Read()
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			return
		}
		encoded, err := json.Marshal(data)
		if err != nil {
			writeWebsocketError(ws, "Invalid command")
			continue
		}
		var command model.FuelCommand
		decoder := json.NewDecoder(bytes.NewReader(encoded))
		decoder.DisallowUnknownFields()
		if err = decoder.Decode(&command); err != nil {
			writeWebsocketError(ws, "Invalid command fields: "+err.Error())
			continue
		}
		actor, authorized := web.fuelActor(r)
		if !authorized {
			writeWebsocketMessage(ws, "fuelResult", field.FuelReply{RequestId: command.RequestId, Error: "Admin session expired; sign in again", State: web.arena.FuelState()})
			continue
		}
		if kind == "getFuelState" {
			writeWebsocketMessage(ws, "manualFuelState", web.arena.FuelState())
			continue
		}
		writeWebsocketMessage(ws, "fuelResult", web.arena.ApplyFuelCommand(kind, command, actor))
	}
}

func (web *Web) manualFuelHistoryHandler(w http.ResponseWriter, r *http.Request) {
	if !web.userIsAdmin(w, r) {
		return
	}
	operations, err := web.arena.Database.GetFuelHistory(r.URL.Query().Get("attemptId"))
	if err != nil {
		handleWebErr(w, err)
		return
	}
	sort.Slice(operations, func(i, j int) bool { return operations[i].Time.After(operations[j].Time) })
	if len(operations) > 100 {
		operations = operations[:100]
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(operations); err != nil {
		log.Println(err)
	}
}
