// Package fueldevice implements the authenticated, bounded Arduino FUEL wire protocol.
package fueldevice

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

type DeviceConfig struct {
	Id       string
	Alliance string
	IP       string
	Key      string
	Source   string
}

type Config struct {
	Enabled bool
	Listen  string
	Devices []DeviceConfig
}

func LoadConfig(path string) (Config, error) {
	cfg := Config{Listen: "10.0.100.5:1760"}
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return cfg, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return cfg, err
	}
	if info.Size() > 8192 {
		return cfg, fmt.Errorf("Arduino configuration exceeds 8192 bytes")
	}
	d := json.NewDecoder(io.LimitReader(f, 8193))
	d.DisallowUnknownFields()
	if err = d.Decode(&cfg); err != nil {
		return cfg, err
	}
	var extra any
	if err = d.Decode(&extra); err != io.EOF {
		return cfg, fmt.Errorf("unexpected trailing configuration")
	}
	return cfg, cfg.Validate()
}

func (cfg Config) Validate() error {
	if !cfg.Enabled {
		return nil
	}
	host, port, err := net.SplitHostPort(cfg.Listen)
	portNumber, portErr := strconv.Atoi(port)
	if err != nil || net.ParseIP(host) == nil || portErr != nil || portNumber < 1024 || portNumber > 65535 || portNumber == 1160 || portNumber == 1750 || portNumber == 8080 {
		return fmt.Errorf("Arduino listen must use a numeric local address and a dedicated TCP port (default 1760; reserved Arena ports excluded)")
	}
	if len(cfg.Devices) != 2 {
		return fmt.Errorf("exactly one RED and one BLUE device required")
	}
	seen := map[string]bool{}
	for _, device := range cfg.Devices {
		if len(device.Id) < 1 || len(device.Id) > 16 || strings.Trim(device.Id, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_") != "" {
			return fmt.Errorf("invalid device ID")
		}
		if device.Alliance != "red" && device.Alliance != "blue" {
			return fmt.Errorf("invalid alliance")
		}
		ip := net.ParseIP(device.IP)
		if ip == nil || ip.To4() == nil || ip.IsUnspecified() || ip.IsMulticast() {
			return fmt.Errorf("invalid device IPv4 address")
		}
		if device.Source != "plc" && device.Source != "arduino" {
			return fmt.Errorf("source must be plc or arduino")
		}
		key, err := hex.DecodeString(device.Key)
		if err != nil || len(key) != 32 || strings.Trim(device.Key, "0") == "" {
			return fmt.Errorf("device %s needs a unique random 32-byte hex key", device.Id)
		}
		for _, identity := range []string{"id:" + device.Id, "alliance:" + device.Alliance, "ip:" + ip.String(), "key:" + strings.ToLower(device.Key)} {
			if seen[identity] {
				return fmt.Errorf("duplicate device identity, alliance, IP or key")
			}
			seen[identity] = true
		}
	}
	return nil
}
