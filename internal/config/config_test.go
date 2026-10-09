package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

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

func getDefaultConfig() Config {
	return Config{
		Interval: 3 * time.Second,
		ProcRoot: "/proc",
	}
}

func TestConfigSuccessLoadFromFillFile(t *testing.T) {
	path := "testdata/config.toml"
	expected := Config{
		Interval: 8 * time.Second,
		ProcRoot: "/test_proc/",
	}

	conf, err := Load(path)

	require.NoError(t, err)
	assert.Equal(t, expected, conf)
}

func TestConfigFileNotFound(t *testing.T) {
	path := "testdata/aodjdaof.toml"
	expected := getDefaultConfig()

	conf, err := Load(path)

	require.NoError(t, err)
	assert.Equal(t, expected, conf)
}

func TestConfigParseError(t *testing.T) {
	path := configPath(t, "error_config.toml", "interval : '12s'\nerror: df", 0o600)
	_, err := Load(path)
	assert.ErrorContains(t, err, "failed to parse config file")
}

func TestConfigReadError(t *testing.T) {
	path := configPath(t, "error_read_config.toml", "interval : '12s'\nerror: df", 0)
	_, err := Load(path)
	require.Error(t, err)
}

// will return default config
func TestConfigSuccessLoadFromBlankFile(t *testing.T) {
	path := configPath(t, "blank_config.toml", "", 0o600)
	expected := getDefaultConfig()

	conf, err := Load(path)

	require.NoError(t, err)
	assert.Equal(t, expected, conf)
}

func TestConfigPartFill(t *testing.T) {
	path := configPath(t, "invalid_interval_config.toml", "interval = '12s'", 0o600)
	expected := getDefaultConfig()
	expected.Interval = 12 * time.Second

	conf, err := Load(path)

	require.NoError(t, err)
	assert.Equal(t, expected, conf)
}

func TestConfigValidateIntervalValue(t *testing.T) {
	path := configPath(t, "invalid_interval_config.toml", "interval = '3ms'\nproc_root = '/proc_root'", 0o600)
	_, err := Load(path)
	assert.ErrorIs(t, err, ErrInvalidInterval)
}

func TestConfigValidateProcRootValue(t *testing.T) {
	path := configPath(t, "invalid_proc_root_config.toml", "interval = '1s'\nproc_root = ''", 0o600)
	_, err := Load(path)
	assert.ErrorIs(t, err, ErrBlankProcRoot)
}

func TestDefaultPathSuccess(t *testing.T) {
	tmpDir := t.TempDir()

	var expectedBase string

	switch runtime.GOOS {
	case "darwin":
		t.Setenv("HOME", tmpDir)
		expectedBase = filepath.Join(tmpDir, "Library", "Application Support")
	case "windows":
		t.Setenv("AppData", tmpDir)
		expectedBase = tmpDir
	default:
		t.Setenv("XDG_CONFIG_HOME", tmpDir)
		expectedBase = tmpDir
	}

	expected := filepath.Join(expectedBase, appName, configFileName)

	result, err := DefaultPath()

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestDefaultPathError(t *testing.T) {
	t.Setenv("HOME", "")
	t.Setenv("AppData", "")
	t.Setenv("XDG_CONFIG_HOME", "")

	result, err := DefaultPath()

	assert.Equal(t, "", result)
	assert.Error(t, err)
}
