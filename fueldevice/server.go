package fueldevice

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

// State is derived by Arena; the transport never runs a match clock.
type State struct {
	Source  string
	Attempt string
	MatchId int
	Flags   uint32
	Phase   int
	Shift   int
}
type Event struct {
	Device, Alliance, Attempt string
	Boot, Sequence, Detected  uint32
	Expected                  State
	ValidLease                bool
}
type Outcome struct {
	Accepted, Duplicate bool
	Error               string
	StorageError        bool
}
type Provider interface {
	FuelDeviceState(alliance string, now time.Time) State
	ApplyFuelDeviceEvent(event Event) Outcome
}
type DeviceStatus struct {
	ActiveSource                            string
	State                                   State
	Id, Alliance, IP, Source                string
	Connected                               bool
	LastSeen                                time.Time
	Detected, Boot                          uint32
	Accepted, Rejected, Duplicates, Dropped uint64
	LastError, Firmware                     string
}
type Status struct {
	Enabled bool
	Listen  string
	Devices []DeviceStatus
}

type Server struct {
	config       Config
	provider     Provider
	mutex        sync.Mutex
	status       map[string]*DeviceStatus
	connections  map[string]net.Conn
	listener     net.Listener
	wg           sync.WaitGroup
	closing      bool
	leaseCounter uint32
	slots        chan struct{}
}

func NewServer(config Config, provider Provider) (*Server, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}
	s := &Server{config: config, provider: provider, status: map[string]*DeviceStatus{}, connections: map[string]net.Conn{}, slots: make(chan struct{}, 8)}
	for _, d := range config.Devices {
		s.status[d.Id] = &DeviceStatus{Id: d.Id, Alliance: d.Alliance, IP: d.IP, Source: d.Source}
	}
	return s, nil
}

func (s *Server) Start() error {
	if !s.config.Enabled {
		return nil
	}
	l, err := net.Listen("tcp4", s.config.Listen)
	if err != nil {
		return err
	}
	s.mutex.Lock()
	s.listener = l
	s.mutex.Unlock()
	s.wg.Add(1)
	go s.accept(l)
	return nil
}

func (s *Server) Close() {
	s.mutex.Lock()
	s.closing = true
	if s.listener != nil {
		s.listener.Close()
	}
	for _, c := range s.connections {
		c.Close()
	}
	s.mutex.Unlock()
	s.wg.Wait()
}

func (s *Server) Status() Status {
	s.mutex.Lock()
	result := Status{Enabled: s.config.Enabled, Listen: s.config.Listen, Devices: []DeviceStatus{}}
	for _, d := range s.config.Devices {
		result.Devices = append(result.Devices, *s.status[d.Id])
	}
	s.mutex.Unlock()
	for i := range result.Devices {
		result.Devices[i].State = s.provider.FuelDeviceState(result.Devices[i].Alliance, time.Now())
		result.Devices[i].ActiveSource = result.Devices[i].State.Source
	}
	return result
}

func (s *Server) accept(listener net.Listener) {
	defer s.wg.Done()
	for {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		select {
		case s.slots <- struct{}{}:
		default:
			conn.Close()
			continue
		}
		s.wg.Add(1)
		go func() { defer s.wg.Done(); defer func() { <-s.slots }(); defer conn.Close(); s.serve(conn) }()
	}
}

func (s *Server) serve(conn net.Conn) {
	host, _, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		return
	}
	var device DeviceConfig
	for _, candidate := range s.config.Devices {
		if net.ParseIP(host).Equal(net.ParseIP(candidate.IP)) {
			device = candidate
			break
		}
	}
	if device.Id == "" {
		return
	}
	key, _ := hex.DecodeString(device.Key)
	read := func() (Frame, error) {
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		buffer := make([]byte, FrameSize)
		if _, err := io.ReadFull(conn, buffer); err != nil {
			return Frame{}, err
		}
		f, err := Decode(buffer, key)
		if err == nil && f.Device != device.Id {
			err = fmt.Errorf("device/IP mismatch")
		}
		return f, err
	}
	var writeMutex sync.Mutex
	var outgoing uint32
	write := func(f Frame) error {
		writeMutex.Lock()
		defer writeMutex.Unlock()
		outgoing++
		f.Sequence, f.Device = outgoing, device.Id
		conn.SetWriteDeadline(time.Now().Add(time.Second))
		data := f.Encode(key)
		for len(data) > 0 {
			n, err := conn.Write(data)
			if err != nil {
				return err
			}
			if n == 0 {
				return io.ErrShortWrite
			}
			data = data[n:]
		}
		return nil
	}
	fail := func(err error) { s.mutex.Lock(); s.status[device.Id].LastError = err.Error(); s.mutex.Unlock() }
	hello, err := read()
	if err != nil {
		fail(err)
		return
	}
	if hello.Type != Hello || hello.Sequence != 0 || hello.Value == 0 {
		fail(fmt.Errorf("invalid HELLO/boot ID"))
		return
	}
	var nonce [16]byte
	if _, err = rand.Read(nonce[:]); err != nil {
		fail(err)
		return
	}
	if err = write(Frame{Type: Challenge, Nonce: nonce}); err != nil {
		fail(err)
		return
	}
	proof, err := read()
	if err != nil || proof.Type != Auth || proof.Sequence != 1 || proof.Nonce != nonce || proof.Value != hello.Value {
		fail(fmt.Errorf("challenge authentication failed"))
		return
	}
	s.mutex.Lock()
	if s.closing {
		s.mutex.Unlock()
		return
	}
	if previous := s.connections[device.Id]; previous != nil {
		previous.Close()
	}
	s.connections[device.Id] = conn
	status := s.status[device.Id]
	if status.Boot != 0 && status.Boot != hello.Value {
		status.LastError = "Device rebooted; unsent sensor events may have been lost"
	}
	status.Connected, status.Boot, status.LastSeen = true, hello.Value, time.Now()
	status.Firmware = fmt.Sprintf("%d.%d.%d", hello.Flags>>16, (hello.Flags>>8)&255, hello.Flags&255)
	s.mutex.Unlock()
	defer func() {
		s.mutex.Lock()
		if s.connections[device.Id] == conn {
			delete(s.connections, device.Id)
			status.Connected = false
		}
		s.mutex.Unlock()
	}()
	// One short-lived lease is shared by state publication and event validation.
	// Old leases are rejected, including late packets at phase/Hub transitions.
	var leaseMutex sync.Mutex
	var lease uint32
	var leased State
	var expires time.Time
	var lastSent time.Time
	issueState := func() error {
		leaseMutex.Lock()
		defer leaseMutex.Unlock()
		now := time.Now()
		state := s.provider.FuelDeviceState(device.Alliance, now)
		if state == leased && now.Before(expires.Add(-250*time.Millisecond)) && now.Sub(lastSent) < 200*time.Millisecond {
			return nil
		}
		if state != leased || !now.Before(expires.Add(-250*time.Millisecond)) {
			s.mutex.Lock()
			s.leaseCounter++
			if s.leaseCounter == 0 {
				s.leaseCounter++
			}
			lease = s.leaseCounter
			s.mutex.Unlock()
			leased, expires = state, now.Add(750*time.Millisecond)
		}
		err := write(Frame{Type: StateMessage, Nonce: nonce, Attempt: AttemptBytes(state.Attempt), Lease: lease, Count: uint32(state.MatchId), Flags: state.Flags, Value: uint32(state.Phase)<<8 | uint32(state.Shift)})
		if err == nil {
			lastSent = now
		}
		return err
	}
	if err = issueState(); err != nil {
		fail(err)
		return
	}
	done := make(chan struct{})
	writerDone := make(chan struct{})
	go func() {
		defer close(writerDone)
		ticker := time.NewTicker(20 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				if err := issueState(); err != nil {
					conn.Close()
					return
				}
			}
		}
	}()
	defer func() { close(done); conn.Close(); <-writerDone }()
	sequence := uint32(1)
	windowStart, windowCount := time.Now(), 0
	for {
		f, err := read()
		if err != nil {
			fail(err)
			return
		}
		if f.Nonce != nonce || sequence == ^uint32(0) || f.Sequence != sequence+1 {
			fail(fmt.Errorf("stale session/transport sequence"))
			return
		}
		sequence = f.Sequence
		if time.Since(windowStart) >= time.Second {
			windowStart, windowCount = time.Now(), 0
		}
		windowCount++
		if windowCount > 500 {
			fail(fmt.Errorf("device rate limit exceeded"))
			return
		}
		s.mutex.Lock()
		status.LastSeen, status.Detected = time.Now(), f.Count
		s.mutex.Unlock()
		if f.Type == Heartbeat {
			s.mutex.Lock()
			status.Dropped = uint64(f.Value)
			s.mutex.Unlock()
			continue
		}
		if f.Type != EventMessage || f.Value == 0 || f.Count < f.Value {
			fail(fmt.Errorf("invalid event"))
			return
		}
		leaseMutex.Lock()
		expected := leased
		valid := f.Lease == lease && time.Now().Before(expires) && f.Attempt == AttemptBytes(expected.Attempt)
		leaseMutex.Unlock()
		outcome := s.provider.ApplyFuelDeviceEvent(Event{Device: device.Id, Alliance: device.Alliance, Attempt: AttemptString(f.Attempt), Boot: hello.Value, Sequence: f.Value, Detected: f.Count, Expected: expected, ValidLease: valid})
		code := AckRejected
		s.mutex.Lock()
		if outcome.Duplicate {
			code = AckDuplicate
			status.Duplicates++
		} else if outcome.Accepted {
			code = AckAccepted
			status.Accepted++
		} else {
			status.Rejected++
			status.LastError = outcome.Error
		}
		if outcome.StorageError {
			code = AckStorageError
		}
		s.mutex.Unlock()
		if err = write(Frame{Type: Ack, Nonce: nonce, Attempt: f.Attempt, Lease: f.Lease, Value: f.Value, Flags: code}); err != nil {
			fail(err)
			return
		}
	}
}
