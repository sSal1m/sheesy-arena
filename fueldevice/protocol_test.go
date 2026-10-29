package fueldevice

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func validConfig() Config {
	return Config{Enabled: true, Listen: "127.0.0.1:1760", Devices: []DeviceConfig{
		{Id: "red", Alliance: "red", IP: "127.0.0.1", Key: strings.Repeat("ab", 32), Source: "arduino"},
		{Id: "blue", Alliance: "blue", IP: "127.0.0.2", Key: strings.Repeat("cd", 32), Source: "plc"},
	}}
}

func TestConfigurationValidation(t *testing.T) {
	require.NoError(t, validConfig().Validate())
	customPort := validConfig()
	customPort.Listen = "127.0.0.1:1761"
	require.NoError(t, customPort.Validate())
	for _, change := range []func(*Config){
		func(c *Config) { c.Devices[1].Id = "red" }, func(c *Config) { c.Devices[1].Alliance = "red" }, func(c *Config) { c.Devices[1].IP = "127.0.0.1" },
		func(c *Config) { c.Devices[1].Key = c.Devices[0].Key }, func(c *Config) { c.Devices[0].Source = "both" }, func(c *Config) { c.Devices[0].Key = strings.Repeat("0", 64) }, func(c *Config) { c.Listen = "0.0.0.0:8080" },
	} {
		cfg := validConfig()
		change(&cfg)
		require.Error(t, cfg.Validate())
	}
	path := filepath.Join(t.TempDir(), "devices.json")
	missing, err := LoadConfig(path)
	require.NoError(t, err)
	assert.False(t, missing.Enabled)
	require.NoError(t, os.WriteFile(path, []byte(`{"Unexpected":1}`), 0600))
	_, err = LoadConfig(path)
	require.Error(t, err)
	data, _ := json.Marshal(validConfig())
	require.NoError(t, os.WriteFile(path, data, 0600))
	_, err = LoadConfig(path)
	require.NoError(t, err)
}

func TestAuthenticatedFrameRejectsTamperAndWrongKey(t *testing.T) {
	key := bytes.Repeat([]byte{0xab}, 32)
	f := Frame{Type: EventMessage, Sequence: 5, Attempt: AttemptBytes("aabaabab-1234-4567-890a-bbbbbbbbbbbb"), Lease: 35, Count: 125, Value: 4, Device: "red", Nonce: [16]byte{1, 2, 3}}
	encoded := f.Encode(key)
	require.Len(t, encoded, 88)
	decoded, err := Decode(encoded, key)
	require.NoError(t, err)
	assert.Equal(t, f, decoded)
	assert.Equal(t, "aabaabab-1234-4567-890a-bbbbbbbbbbbb", AttemptString(f.Attempt))
	for index := range encoded {
		corrupted := append([]byte{}, encoded...)
		corrupted[index] ^= 1
		_, err = Decode(corrupted, key)
		require.Error(t, err)
	}
	_, err = Decode(encoded, bytes.Repeat([]byte{2}, 32))
	require.Error(t, err)
	_, err = Decode(encoded[:87], key)
	require.Error(t, err)
}
