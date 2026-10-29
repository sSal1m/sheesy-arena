package fueldevice

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
)

const FrameSize = 88
const (
	Hello        byte = 1
	Challenge    byte = 2
	Auth         byte = 3
	StateMessage byte = 4
	EventMessage byte = 5
	Ack          byte = 6
	Heartbeat    byte = 7
)
const (
	FlagStarted uint32 = 1 << iota
	FlagAuto
	FlagTeleop
	FlagHubActive
	FlagScoring
	FlagAborted
	FlagComplete
	FlagManual
)
const (
	AckAccepted uint32 = iota
	AckDuplicate
	AckRejected
	AckStorageError
)

type Frame struct {
	Type     byte
	Sequence uint32
	Attempt  [16]byte
	Lease    uint32
	Count    uint32
	Flags    uint32
	Value    uint32
	Device   string
	Nonce    [16]byte
}

func (f Frame) Encode(key []byte) []byte {
	b := make([]byte, FrameSize)
	copy(b, "F2")
	b[2], b[3] = 1, f.Type
	binary.BigEndian.PutUint32(b[4:8], f.Sequence)
	copy(b[8:24], f.Attempt[:])
	for index, value := range []uint32{f.Lease, f.Count, f.Flags, f.Value} {
		binary.BigEndian.PutUint32(b[24+index*4:28+index*4], value)
	}
	copy(b[40:56], f.Device)
	copy(b[56:72], f.Nonce[:])
	mac := hmac.New(sha256.New, key)
	mac.Write(b[:72])
	copy(b[72:], mac.Sum(nil)[:16])
	return b
}

func Decode(b, key []byte) (Frame, error) {
	var f Frame
	if len(b) != FrameSize || string(b[:2]) != "F2" || b[2] != 1 {
		return f, fmt.Errorf("invalid frame/version")
	}
	mac := hmac.New(sha256.New, key)
	mac.Write(b[:72])
	if !hmac.Equal(b[72:], mac.Sum(nil)[:16]) {
		return f, fmt.Errorf("invalid message authentication")
	}
	f.Type, f.Sequence = b[3], binary.BigEndian.Uint32(b[4:8])
	copy(f.Attempt[:], b[8:24])
	f.Lease, f.Count, f.Flags, f.Value = binary.BigEndian.Uint32(b[24:28]), binary.BigEndian.Uint32(b[28:32]), binary.BigEndian.Uint32(b[32:36]), binary.BigEndian.Uint32(b[36:40])
	f.Device = strings.TrimRight(string(b[40:56]), "\x00")
	copy(f.Nonce[:], b[56:72])
	return f, nil
}

func AttemptBytes(id string) [16]byte {
	var result [16]byte
	b, _ := hex.DecodeString(strings.ReplaceAll(id, "-", ""))
	copy(result[:], b)
	return result
}

func AttemptString(id [16]byte) string {
	b := hex.EncodeToString(id[:])
	return b[:8] + "-" + b[8:12] + "-" + b[12:16] + "-" + b[16:20] + "-" + b[20:]
}
