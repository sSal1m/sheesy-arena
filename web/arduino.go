package web

import (
	"encoding/json"
	"github.com/Team254/cheesy-arena/fueldevice"
	"github.com/Team254/cheesy-arena/model"
	"io"
	"net/http"
)

func (web *Web) arduinoHandler(w http.ResponseWriter, r *http.Request) {
	if !web.userIsAdmin(w, r) {
		return
	}
	t, err := web.parseFiles("templates/arduino.html", "templates/base.html")
	if err != nil {
		handleWebErr(w, err)
		return
	}
	if err = t.ExecuteTemplate(w, "base", struct{ *model.EventSettings }{web.arena.EventSettings}); err != nil {
		handleWebErr(w, err)
	}
}

func (web *Web) arduinoStatusHandler(w http.ResponseWriter, r *http.Request) {
	if !web.userIsAdmin(w, r) {
		return
	}
	status := fueldevice.Status{Devices: []fueldevice.DeviceStatus{}}
	if web.arena.FuelDevices != nil {
		status = web.arena.FuelDevices.Status()
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func (web *Web) arduinoSimulationHandler(w http.ResponseWriter, r *http.Request) {
	if !web.userIsAdmin(w, r) {
		return
	}
	var args struct {
		Alliance, Phase string
		Count           int
		Reset           bool
	}
	decoder := json.NewDecoder(io.LimitReader(r.Body, 1024))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&args); err != nil {
		http.Error(w, "Invalid simulation input", 400)
		return
	}
	var extra any
	if decoder.Decode(&extra) != io.EOF {
		http.Error(w, "Unexpected simulation input", 400)
		return
	}
	result, err := web.arena.SimulateFuelDevice(args.Alliance, args.Phase, args.Count, args.Reset)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
