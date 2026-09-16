package config_test

import (
	"go-sysmon/internal/config"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getDefaultConfig() config.Config {
	return config.Config{
		Interval: 3 * time.Second,
		ProcRoot: "/proc",
	}
}

func TestConfigParseError(t *testing.T) {
	path := "testdata/error_config.toml"
	_, err := config.Load(path)
	assert.ErrorContains(t, err, "failed to parse config file")
}

func TestConfigReadError(t *testing.T) {
	path := "testdata/error_read_config"
	_, err := config.Load(path)
	require.Error(t, err)
}

func TestConfigSuccessLoadFromFillFile(t *testing.T) {
	path := "testdata/fill_config.toml"
	expected := config.Config{
		Interval: 8 * time.Second,
		ProcRoot: "/test_proc/",
	}

	conf, err := config.Load(path)

	require.NoError(t, err)
	assert.Equal(t, expected, conf)
}

// will return default config
func TestConfigSuccessLoadFromBlankFile(t *testing.T) {
	path := "testdata/blank_config.toml"
	expected := getDefaultConfig()

	conf, err := config.Load(path)

	require.NoError(t, err)
	assert.Equal(t, expected, conf)
}

func TestConfigFileNotFound(t *testing.T) {
	path := "testdata/aodjdaof.toml"
	expected := getDefaultConfig()

	conf, err := config.Load(path)

	require.NoError(t, err)
	assert.Equal(t, expected, conf)
}

func TestConfigValidateIntervalValue(t *testing.T) {
	path := "testdata/invalid_interval_config.toml"
	_, err := config.Load(path)
	assert.ErrorIs(t, err, config.ErrInvalidInterval)
}

func TestConfigValidateProcRootValue(t *testing.T) {
	path := "testdata/invalid_proc_root_config.toml"
	_, err := config.Load(path)
	assert.ErrorIs(t, err, config.ErrBlankProcRoot)
}
