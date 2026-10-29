// Copyright 2014 Team 254. All Rights Reserved.
// Author: pat@patfairbank.com (Patrick Fairbank)

package main

import (
	"flag"
	"github.com/Team254/cheesy-arena/field"
	"github.com/Team254/cheesy-arena/fueldevice"
	"github.com/Team254/cheesy-arena/network"
	"github.com/Team254/cheesy-arena/web"
	"log"
)

const eventDbPath = "./event.db"
const httpPort = 8080

// Main entry point for the application.
func main() {
	var deviceConfig string
	flag.StringVar(&deviceConfig, "arduino-config", "arduino.local.json", "Arduino Ethernet configuration (missing file disables devices)")
	flag.BoolVar(&network.DevMode, "dev", false, "Bind driver station listeners to all IP addresses for development")
	flag.Parse()

	arena, err := field.NewArena(eventDbPath)
	if err != nil {
		log.Fatalln("Error during startup: ", err)
	}
	config, err := fueldevice.LoadConfig(deviceConfig)
	if err != nil {
		log.Fatalln("Invalid Arduino configuration:", err)
	}
	if err = arena.ConfigureFuelDevices(config); err != nil {
		log.Fatalln("Arduino scoring setup:", err)
	}
	arena.FuelDevices, err = fueldevice.NewServer(config, arena)
	if err != nil {
		log.Fatalln(err)
	}
	if err = arena.FuelDevices.Start(); err != nil {
		log.Fatalln("Arduino listener:", err)
	}
	defer arena.FuelDevices.Close()

	// Start the web server in a separate goroutine.
	web := web.NewWeb(arena)
	go web.ServeWebInterface(httpPort)

	// Run the arena state machine in the main thread.
	arena.Run()
}
