package plc

import "time"

// FuelSnapshot is published only after a complete successful register read.
// Sequence is host-side ordering; the PLC protocol has no hardware match ID.
type FuelSnapshot struct {
	Red        int
	Blue       int
	Sequence   uint64
	Connection uint64
	Reset      uint64
	ReadAt     time.Time
	Healthy    bool
}

func (plc *ModbusPlc) GetFuelSnapshot() FuelSnapshot {
	plc.fuelMutex.RLock()
	defer plc.fuelMutex.RUnlock()
	return plc.fuelSnapshot
}

func (plc *ModbusPlc) publishFuelRegisters(registers [registerCount]uint16, epoch FuelSnapshot) {
	plc.fuelMutex.Lock()
	defer plc.fuelMutex.Unlock()
	if epoch.Reset != plc.fuelSnapshot.Reset || epoch.Connection != plc.fuelSnapshot.Connection {
		return
	}
	plc.fuelSnapshot.Red = int(registers[redHubTotal])
	plc.fuelSnapshot.Blue = int(registers[blueHubTotal])
	plc.fuelSnapshot.Sequence++
	plc.fuelSnapshot.ReadAt = time.Now()
	plc.fuelSnapshot.Healthy = true
}
