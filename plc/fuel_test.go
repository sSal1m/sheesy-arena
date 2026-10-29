package plc

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestFuelSnapshotRejectsReadFromPreviousReset(t *testing.T) {
	plc := new(ModbusPlc)
	plc.SetAddress("")
	epoch := plc.GetFuelSnapshot()
	plc.ResetMatch()
	var registers [registerCount]uint16
	registers[redHubTotal] = 99
	plc.publishFuelRegisters(registers, epoch)
	assert.True(t, plc.GetFuelSnapshot().ReadAt.IsZero())
	epoch = plc.GetFuelSnapshot()
	registers[redHubTotal] = 0
	plc.publishFuelRegisters(registers, epoch)
	snapshot := plc.GetFuelSnapshot()
	assert.True(t, snapshot.Healthy)
	assert.Equal(t, uint64(1), snapshot.Sequence)
	assert.True(t, time.Since(snapshot.ReadAt) < time.Second)
}

func TestFuelSnapshotConcurrentReaders(t *testing.T) {
	plc := new(ModbusPlc)
	var group sync.WaitGroup
	for i := 0; i < 4; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			for i := 0; i < 100; i++ {
				plc.GetFuelSnapshot()
			}
		}()
	}
	for i := 0; i < 100; i++ {
		var registers [registerCount]uint16
		registers[redHubTotal], registers[blueHubTotal] = uint16(i), uint16(i)
		plc.publishFuelRegisters(registers, plc.GetFuelSnapshot())
	}
	group.Wait()
	assert.Equal(t, 99, plc.GetFuelSnapshot().Red)
}
