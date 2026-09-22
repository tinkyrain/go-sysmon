package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/tinkyrain/go-sysmon/internal/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func configPath(t *testing.T, filename, filecontent string, perm os.FileMode) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, filename)
	require.NoError(t, os.WriteFile(path, []byte(filecontent), perm))
	return path
}

func getDefaultConfig() config.Config {
	return config.Config{
		Interval: 3 * time.Second,
		ProcRoot: "/proc",
	}
}

func TestConfigSuccessLoadFromFillFile(t *testing.T) {
	path := "testdata/config.toml"
	expected := config.Config{
		Interval: 8 * time.Second,
		ProcRoot: "/test_proc/",
	}

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

func TestConfigParseError(t *testing.T) {
	path := configPath(t, "error_config.toml", "interval : '12s'\nerror: df", 0o600)
	_, err := config.Load(path)
	assert.ErrorContains(t, err, "failed to parse config file")
}

func TestConfigReadError(t *testing.T) {
	path := configPath(t, "error_read_config.toml", "interval : '12s'\nerror: df", 0)
	_, err := config.Load(path)
	require.Error(t, err)
}

// will return default config
func TestConfigSuccessLoadFromBlankFile(t *testing.T) {
	path := configPath(t, "blank_config.toml", "", 0o600)
	expected := getDefaultConfig()

	conf, err := config.Load(path)

	require.NoError(t, err)
	assert.Equal(t, expected, conf)
}

func TestConfigPartFill(t *testing.T) {
	path := configPath(t, "invalid_interval_config.toml", "interval = '12s'", 0o600)
	expected := getDefaultConfig()
	expected.Interval = 12 * time.Second

	conf, err := config.Load(path)

	require.NoError(t, err)
	assert.Equal(t, expected, conf)
}

func TestConfigValidateIntervalValue(t *testing.T) {
	path := configPath(t, "invalid_interval_config.toml", "interval = '3ms'\nproc_root = '/proc_root'", 0o600)
	_, err := config.Load(path)
	assert.ErrorIs(t, err, config.ErrInvalidInterval)
}

func TestConfigValidateProcRootValue(t *testing.T) {
	path := configPath(t, "invalid_proc_root_config.toml", "interval = '1s'\nproc_root = ''", 0o600)
	_, err := config.Load(path)
	assert.ErrorIs(t, err, config.ErrBlankProcRoot)
}
