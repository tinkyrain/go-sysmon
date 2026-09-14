package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Interval time.Duration `toml:"interval"`  // Interval in seconds for collection metrics. Example "3s"
	ProcRoot string        `toml:"proc_root"` // ProcRoot for metricsPath
}

const (
	appName        = "go-sysmon"
	configFileName = "config.toml"
)

func defaultSettings() Config {
	return Config{
		Interval: 3 * time.Second,
		ProcRoot: "/proc",
	}
}

func DefaultPath() (string, error) {
	systemConfigPath, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(systemConfigPath, appName, configFileName), nil
}

func Load(path string) (Config, error) {
	config := defaultSettings()

	data, err := os.ReadFile(path) // read config file
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return config, nil // Error not return, but we are used default settings
		}
		return config, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	// Decode toml cfg in structure
	if err := toml.Unmarshal(data, &config); err != nil {
		return config, fmt.Errorf("failed to parse config file %s: %w", path, err)
	}
	if err := validate(config); err != nil {
		return config, fmt.Errorf("failed to validate config file %s: %w", path, err)
	}

	return config, nil
}

func validate(config Config) error {
	if config.Interval <= 100*time.Millisecond {
		return fmt.Errorf("interval %v must be greater than 100ms", config.Interval)
	}
	if config.ProcRoot == "" {
		return errors.New("path cannot be empty")
	}
	return nil
}
