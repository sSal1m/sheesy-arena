package field

import (
	"encoding/hex"
	"github.com/Team254/cheesy-arena/fueldevice"
	"github.com/Team254/cheesy-arena/game"
	"github.com/Team254/cheesy-arena/model"
	"github.com/Team254/cheesy-arena/plc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func sensorConfig() fueldevice.Config {
	return fueldevice.Config{Enabled: true, Listen: "127.0.0.1:1760", Devices: []fueldevice.DeviceConfig{
		{Id: "fuel-red-01", Alliance: "red", IP: "127.0.0.1", Key: strings.Repeat("ab", 32), Source: "arduino"},
		{Id: "fuel-blue-01", Alliance: "blue", IP: "127.0.0.2", Key: strings.Repeat("cd", 32), Source: "arduino"},
	}}
}
func sensorEvent(arena *Arena, alliance string, sequence uint32) fueldevice.Event {
	device := "fuel-red-01"
	if alliance == "blue" {
		device = "fuel-blue-01"
	}
	state := arena.FuelDeviceState(alliance, time.Now())
	return fueldevice.Event{Device: device, Alliance: alliance, Attempt: state.Attempt, Boot: 1, Sequence: sequence, Detected: sequence, Expected: state, ValidLease: true}
}

func TestArduinoOwnershipManualCorrectionAndPhases(t *testing.T) {
	arena := setupTestArena(t)
	require.NoError(t, arena.ConfigureFuelDevices(sensorConfig()))
	fuelPhase(arena, AutoPeriod)
	for _, alliance := range []string{"red", "blue"} {
		assert.True(t, arena.ApplyFuelDeviceEvent(sensorEvent(arena, alliance, 1)).Accepted)
		require.Empty(t, arena.ApplyFuelCommand("adjustFuel", fuelCommand(arena, alliance, "auto", 5), "operator").Error)
	}
	assert.Equal(t, 6, arena.FuelState().Red.Auto)
	assert.Equal(t, 6, arena.FuelState().Blue.Auto)
	// PLC readings cannot overwrite, increment or consume Arduino-owned counts.
	arena.processFuelSnapshot(plc.FuelSnapshot{Sequence: 1, Connection: 1, Healthy: true, ReadAt: time.Now()}, time.Now())
	arena.processFuelSnapshot(plc.FuelSnapshot{Sequence: 2, Connection: 1, Red: 200, Blue: 300, Healthy: true, ReadAt: time.Now()}, time.Now())
	assert.Equal(t, 6, arena.FuelState().Red.Auto)
	assert.True(t, arena.ApplyFuelDeviceEvent(sensorEvent(arena, "red", 2)).Accepted)
	assert.Equal(t, 7, arena.FuelState().Red.Auto)
	arena.handleAutoWinner()
	fuelPhase(arena, TeleopPeriod)
	// RED won AUTO, so after transition grace RED is inactive in Shift1.
	arena.MatchStartTime = time.Now().Add(-game.GetDurationToTeleopStart() - time.Duration(game.MatchTiming.TransitionShiftDurationSec+game.ScoringGracePeriodSec+1)*time.Second)
	arena.syncFuelPhase()
	assert.False(t, arena.ApplyFuelDeviceEvent(sensorEvent(arena, "red", 3)).Accepted)
	assert.True(t, arena.ApplyFuelDeviceEvent(sensorEvent(arena, "blue", 2)).Accepted)
	assert.Equal(t, 1, arena.FuelState().Blue.Teleop)
	late := sensorEvent(arena, "blue", 3)
	late.ValidLease = false
	assert.NotEmpty(t, arena.ApplyFuelDeviceEvent(late).Error)
	late = sensorEvent(arena, "blue", 3)
	late.Attempt = "old-attempt"
	assert.NotEmpty(t, arena.ApplyFuelDeviceEvent(late).Error)
	wrong := sensorEvent(arena, "red", 4)
	wrong.Device = "fuel-blue-01"
	assert.NotEmpty(t, arena.ApplyFuelDeviceEvent(wrong).Error)
	// Actual phase changes invalidate an otherwise valid old state lease.
	oldState := sensorEvent(arena, "blue", 4)
	arena.MatchStartTime = time.Now().Add(-game.GetDurationToTeleopEnd() - 4*time.Second)
	arena.MatchState = PostMatch
	arena.syncFuelPhase()
	assert.NotEmpty(t, arena.ApplyFuelDeviceEvent(oldState).Error)
	assert.False(t, arena.ApplyFuelDeviceEvent(sensorEvent(arena, "blue", 5)).Accepted)
	_, finish, err := arena.BeginFuelCommit()
	require.NoError(t, err)
	finish(true)
	assert.False(t, arena.ApplyFuelDeviceEvent(sensorEvent(arena, "blue", 6)).Accepted)
}

func TestArduinoSourcePolicyMixedAndManualMode(t *testing.T) {
	arena := setupTestArena(t)
	cfg := sensorConfig()
	cfg.Devices[1].Source = "plc"
	require.NoError(t, arena.ConfigureFuelDevices(cfg))
	fuelPhase(arena, AutoPeriod)
	assert.True(t, arena.ApplyFuelDeviceEvent(sensorEvent(arena, "red", 1)).Accepted)
	assert.False(t, arena.ApplyFuelDeviceEvent(sensorEvent(arena, "blue", 1)).Accepted)
	sample := plc.FuelSnapshot{Sequence: 1, Connection: 1, Healthy: true, ReadAt: time.Now()}
	arena.processFuelSnapshot(sample, time.Now())
	sample.Sequence++
	sample.Red = 100
	sample.Blue = 10
	sample.ReadAt = time.Now()
	arena.processFuelSnapshot(sample, time.Now())
	assert.Equal(t, 1, arena.FuelState().Red.Auto)
	assert.Equal(t, 10, arena.FuelState().Blue.Auto)
	require.NoError(t, arena.AbortMatch())
	require.NoError(t, arena.ResetMatch())
	require.NoError(t, arena.LoadTestMatch())
	setFuelMode(t, arena, model.FuelManual)
	fuelPhase(arena, AutoPeriod)
	assert.False(t, arena.ApplyFuelDeviceEvent(sensorEvent(arena, "red", 2)).Accepted)
	require.Empty(t, arena.ApplyFuelCommand("adjustFuel", fuelCommand(arena, "red", "auto", 1), "a").Error)
	assert.Equal(t, 1, arena.FuelState().Red.Auto)
}

func TestArduinoDurableDedupAndSafeRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "event.db")
	arena, err := NewArena(path)
	require.NoError(t, err)
	require.NoError(t, arena.ConfigureFuelDevices(sensorConfig()))
	fuelPhase(arena, AutoPeriod)
	event := sensorEvent(arena, "red", 1)
	assert.True(t, arena.ApplyFuelDeviceEvent(event).Accepted)
	assert.True(t, arena.ApplyFuelDeviceEvent(event).Duplicate)
	event.Boot = 2
	assert.True(t, arena.ApplyFuelDeviceEvent(event).Accepted)
	require.NoError(t, arena.Database.Close())
	recovered, err := NewArena(path)
	require.NoError(t, err)
	defer recovered.Database.Close()
	require.NoError(t, recovered.ConfigureFuelDevices(sensorConfig()))
	assert.Equal(t, 2, recovered.FuelState().Red.Auto)
	assert.True(t, recovered.ApplyFuelDeviceEvent(event).Duplicate)
	assert.Zero(t, recovered.FuelDeviceState("red", time.Now()).Flags&fueldevice.FlagScoring)
	event.Sequence++
	assert.False(t, recovered.ApplyFuelDeviceEvent(event).Accepted)
}

type wireSensor struct {
	t     *testing.T
	conn  net.Conn
	key   []byte
	id    string
	nonce [16]byte
	seq   uint32
	state fueldevice.Frame
}

func connectSensor(t *testing.T, config fueldevice.DeviceConfig, boot uint32) *wireSensor {
	dialer := net.Dialer{LocalAddr: &net.TCPAddr{IP: net.ParseIP(config.IP)}, Timeout: time.Second}
	conn, err := dialer.Dial("tcp4", "127.0.0.1:1760")
	require.NoError(t, err)
	key, _ := hex.DecodeString(config.Key)
	sensor := &wireSensor{t: t, conn: conn, key: key, id: config.Id}
	t.Cleanup(func() { conn.Close() })
	_, err = conn.Write((fueldevice.Frame{Type: fueldevice.Hello, Device: config.Id, Value: boot, Flags: 0x010000}).Encode(key))
	require.NoError(t, err)
	challenge := sensor.read()
	require.Equal(t, fueldevice.Challenge, challenge.Type)
	sensor.nonce = challenge.Nonce
	proof := (fueldevice.Frame{Type: fueldevice.Auth, Sequence: 1, Device: config.Id, Nonce: sensor.nonce, Value: boot}).Encode(key)
	// Deliberately fragment authentication across TCP writes.
	for start := 0; start < len(proof); start += 7 {
		end := start + 7
		if end > len(proof) {
			end = len(proof)
		}
		_, err = conn.Write(proof[start:end])
		require.NoError(t, err)
	}
	sensor.seq = 1
	sensor.until(fueldevice.StateMessage)
	return sensor
}
func (s *wireSensor) read() fueldevice.Frame {
	s.conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	b := make([]byte, fueldevice.FrameSize)
	_, err := io.ReadFull(s.conn, b)
	require.NoError(s.t, err)
	f, err := fueldevice.Decode(b, s.key)
	require.NoError(s.t, err)
	return f
}
func (s *wireSensor) until(kind byte) fueldevice.Frame {
	for {
		f := s.read()
		if f.Type == fueldevice.StateMessage {
			s.state = f
		}
		if f.Type == kind {
			return f
		}
	}
}
func (s *wireSensor) event(number uint32) fueldevice.Frame {
	s.seq++
	f := fueldevice.Frame{Type: fueldevice.EventMessage, Sequence: s.seq, Device: s.id, Nonce: s.nonce, Attempt: s.state.Attempt, Lease: s.state.Lease, Count: number, Value: number}
	_, err := s.conn.Write(f.Encode(s.key))
	require.NoError(s.t, err)
	return s.until(fueldevice.Ack)
}

func TestArduinoTcpTwoDevicesReconnectAndRestart(t *testing.T) {
	arena := setupTestArena(t)
	cfg := sensorConfig()
	require.NoError(t, arena.ConfigureFuelDevices(cfg))
	fuelPhase(arena, AutoPeriod)
	server, err := fueldevice.NewServer(cfg, arena)
	require.NoError(t, err)
	require.NoError(t, server.Start())
	defer server.Close()
	red, blue := connectSensor(t, cfg.Devices[0], 1), connectSensor(t, cfg.Devices[1], 1)
	var group sync.WaitGroup
	for _, device := range []*wireSensor{red, blue} {
		group.Add(1)
		go func(s *wireSensor) {
			defer group.Done()
			ack := s.event(1)
			assert.Equal(t, fueldevice.AckAccepted, ack.Flags)
		}(device)
	}
	group.Wait()
	assert.Equal(t, 1, arena.FuelState().Red.Auto)
	assert.Equal(t, 1, arena.FuelState().Blue.Auto)
	red.conn.Close()
	assert.Equal(t, fueldevice.AckAccepted, blue.event(2).Flags)
	red = connectSensor(t, cfg.Devices[0], 1)
	assert.Equal(t, fueldevice.AckDuplicate, red.event(1).Flags)
	red.conn.Close()
	red = connectSensor(t, cfg.Devices[0], 2)
	assert.Equal(t, fueldevice.AckAccepted, red.event(1).Flags)
	assert.Equal(t, 2, arena.FuelState().Red.Auto)
	assert.Equal(t, 2, arena.FuelState().Blue.Auto)
	// A coalesced duplicate + fresh event is parsed as two bounded frames.
	frames := []byte{}
	for _, number := range []uint32{1, 2} {
		red.seq++
		frames = append(frames, (fueldevice.Frame{Type: fueldevice.EventMessage, Sequence: red.seq, Device: red.id, Nonce: red.nonce, Attempt: red.state.Attempt, Lease: red.state.Lease, Count: number, Value: number}).Encode(red.key)...)
	}
	_, err = red.conn.Write(frames)
	require.NoError(t, err)
	assert.Equal(t, fueldevice.AckDuplicate, red.until(fueldevice.Ack).Flags)
	assert.Equal(t, fueldevice.AckAccepted, red.until(fueldevice.Ack).Flags)
	status := server.Status()
	assert.Equal(t, "1.0.0", status.Devices[0].Firmware)
	assert.Equal(t, 3, arena.FuelState().Red.Auto)
	// New connection cannot accidentally reuse the previous connection's lease.
	oldLease := red.state.Lease
	red.conn.Close()
	red = connectSensor(t, cfg.Devices[0], 2)
	assert.NotEqual(t, oldLease, red.state.Lease)
	red.state.Lease = oldLease
	assert.Equal(t, fueldevice.AckRejected, red.event(3).Flags)
}

func TestArduinoSimulationNeverTouchesLiveScores(t *testing.T) {
	arena := setupTestArena(t)
	fuelPhase(arena, AutoPeriod)
	before := arena.ScoreResultSnapshot()
	result, err := arena.SimulateFuelDevice("red", "shift1", 1, false)
	require.NoError(t, err)
	assert.True(t, result.Simulation)
	assert.Equal(t, 1, result.Red.Teleop)
	result, err = arena.SimulateFuelDevice("blue", "shift1", 1, false)
	require.NoError(t, err)
	assert.Zero(t, result.Blue.Teleop)
	assert.Equal(t, before, arena.ScoreResultSnapshot())
}

func TestArduinoTcpRejectsIdentityDelayAndSessionReplay(t *testing.T) {
	arena := setupTestArena(t)
	cfg := sensorConfig()
	require.NoError(t, arena.ConfigureFuelDevices(cfg))
	fuelPhase(arena, AutoPeriod)
	server, err := fueldevice.NewServer(cfg, arena)
	require.NoError(t, err)
	require.NoError(t, server.Start())
	defer server.Close()
	for _, test := range []struct{ ip, id, key string }{
		{"127.0.0.1", cfg.Devices[1].Id, cfg.Devices[1].Key},       // BLUE identity on RED IP.
		{"127.0.0.3", cfg.Devices[0].Id, cfg.Devices[0].Key},       // Unknown source IP.
		{"127.0.0.1", cfg.Devices[0].Id, strings.Repeat("12", 32)}, // Wrong key.
	} {
		dialer := net.Dialer{LocalAddr: &net.TCPAddr{IP: net.ParseIP(test.ip)}}
		c, err := dialer.Dial("tcp4", cfg.Listen)
		require.NoError(t, err)
		key, _ := hex.DecodeString(test.key)
		_, err = c.Write((fueldevice.Frame{Type: fueldevice.Hello, Device: test.id, Value: 1}).Encode(key))
		require.NoError(t, err)
		c.SetReadDeadline(time.Now().Add(time.Second))
		_, err = c.Read(make([]byte, 88))
		require.Error(t, err)
		c.Close()
	}
	red := connectSensor(t, cfg.Devices[0], 1)
	lease := red.state.Lease
	for red.state.Lease == lease {
		red.until(fueldevice.StateMessage)
	}
	red.state.Lease = lease
	assert.Equal(t, fueldevice.AckRejected, red.event(1).Flags)
	red.until(fueldevice.StateMessage)
	assert.Equal(t, fueldevice.AckAccepted, red.event(2).Flags)
	// TCP-level replay with a repeated connection sequence closes the session.
	_, err = red.conn.Write((fueldevice.Frame{Type: fueldevice.EventMessage, Sequence: red.seq, Device: red.id, Nonce: red.nonce, Attempt: red.state.Attempt, Lease: red.state.Lease, Count: 3, Value: 3}).Encode(red.key))
	require.NoError(t, err)
	red.conn.SetReadDeadline(time.Now().Add(time.Second))
	_, err = red.conn.Read(make([]byte, 88))
	require.Error(t, err)
	assert.Equal(t, 1, arena.FuelState().Red.Auto)
}

func TestArduinoEventsCannotReachReplayOfSameMatch(t *testing.T) {
	arena := setupTestArena(t)
	require.NoError(t, arena.ConfigureFuelDevices(sensorConfig()))
	fuelPhase(arena, AutoPeriod)
	stale := sensorEvent(arena, "red", 1)
	require.NoError(t, arena.AbortMatch())
	require.NoError(t, arena.ResetMatch())
	require.NoError(t, arena.LoadTestMatch())
	fuelPhase(arena, AutoPeriod)
	assert.False(t, arena.ApplyFuelDeviceEvent(stale).Accepted)
	assert.Zero(t, arena.FuelState().Red.Auto)
}

func TestArduinoAndManualConcurrentTransactions(t *testing.T) {
	arena := setupTestArena(t)
	require.NoError(t, arena.ConfigureFuelDevices(sensorConfig()))
	fuelPhase(arena, AutoPeriod)
	assert.Error(t, arena.ConfigureFuelDevices(sensorConfig()))
	var group sync.WaitGroup
	for number := uint32(1); number <= 20; number++ {
		event := sensorEvent(arena, "red", number)
		command := fuelCommand(arena, "red", "auto", 1)
		group.Add(2)
		go func() { defer group.Done(); assert.True(t, arena.ApplyFuelDeviceEvent(event).Accepted) }()
		go func() {
			defer group.Done()
			assert.Empty(t, arena.ApplyFuelCommand("adjustFuel", command, "operator").Error)
		}()
	}
	group.Wait()
	state := arena.FuelState()
	assert.Equal(t, 40, state.Red.Auto)
	assert.Equal(t, 20, state.Red.RawAuto)
	assert.Equal(t, 20, state.Red.AutoAdjustment)
	checkpoint, err := arena.Database.GetFuelCheckpoint()
	require.NoError(t, err)
	assert.Equal(t, 40, checkpoint.Result.RedScoreSummary().Score)
}

func TestArduinoNormalAndAbortScoringGrace(t *testing.T) {
	for _, abort := range []bool{false, true} {
		t.Run(map[bool]string{false: "normal", true: "abort"}[abort], func(t *testing.T) {
			arena := setupTestArena(t)
			require.NoError(t, arena.ConfigureFuelDevices(sensorConfig()))
			fuelPhase(arena, AutoPeriod)
			if abort {
				require.NoError(t, arena.AbortMatch())
			} else {
				arena.MatchState = PostMatch
				arena.MatchStartTime = time.Now().Add(-game.GetDurationToTeleopEnd() - time.Second)
			}
			arena.syncFuelPhase()
			state := arena.FuelDeviceState("red", time.Now())
			assert.NotZero(t, state.Flags&fueldevice.FlagComplete)
			assert.NotZero(t, state.Flags&fueldevice.FlagScoring)
			assert.Zero(t, state.Flags&fueldevice.FlagHubActive)
			if abort {
				assert.NotZero(t, state.Flags&fueldevice.FlagAborted)
			}
			assert.True(t, arena.ApplyFuelDeviceEvent(sensorEvent(arena, "red", 1)).Accepted)
			_, _, err := arena.BeginFuelCommit()
			require.ErrorContains(t, err, "grace")
			if abort {
				arena.matchStopTime = time.Now().Add(-4 * time.Second)
			} else {
				arena.MatchStartTime = time.Now().Add(-game.GetDurationToTeleopEnd() - 4*time.Second)
			}
			arena.syncFuelPhase()
			assert.False(t, arena.ApplyFuelDeviceEvent(sensorEvent(arena, "red", 2)).Accepted)
		})
	}
}

func TestArduinoStorageFailureDoesNotApplyEvent(t *testing.T) {
	arena := setupTestArena(t)
	require.NoError(t, arena.ConfigureFuelDevices(sensorConfig()))
	fuelPhase(arena, AutoPeriod)
	event := sensorEvent(arena, "red", 1)
	require.NoError(t, arena.Database.Close())
	outcome := arena.ApplyFuelDeviceEvent(event)
	assert.True(t, outcome.StorageError)
	assert.False(t, outcome.Accepted)
	assert.Zero(t, arena.FuelState().Red.Auto)
}

func TestArduinoJournalUsesPhaseTotalAcrossShifts(t *testing.T) {
	arena := setupTestArena(t)
	require.NoError(t, arena.ConfigureFuelDevices(sensorConfig()))
	fuelPhase(arena, TeleopPeriod)
	assert.True(t, arena.ApplyFuelDeviceEvent(sensorEvent(arena, "red", 1)).Accepted)
	arena.MatchStartTime = time.Now().Add(-game.GetDurationToTeleopStart() - time.Duration(game.MatchTiming.TransitionShiftDurationSec+game.ScoringGracePeriodSec+1)*time.Second)
	arena.syncFuelPhase()
	assert.True(t, arena.ApplyFuelDeviceEvent(sensorEvent(arena, "red", 2)).Accepted)
	history, err := arena.Database.GetFuelHistory(arena.FuelState().HistoryId)
	require.NoError(t, err)
	found := false
	for _, op := range history {
		if op.EventSequence == 2 {
			found = true
			assert.Equal(t, 1, op.Before)
			assert.Equal(t, 2, op.After)
		}
	}
	assert.True(t, found)
	assert.Equal(t, 2, arena.FuelState().Red.Teleop)
}
