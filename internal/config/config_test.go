package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func configPath(t *testing.T, fileContent string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, configFileName)
	require.NoError(t, os.WriteFile(path, []byte(fileContent), 0o600))
	return path
}

func defaultConfig() Config {
	return Config{
		Interval: 3 * time.Second,
		ProcRoot: "/proc",
	}
}

func TestLoadSuccess(t *testing.T) {
	cases := []struct {
		name        string
		fileContent string
		expected    Config
	}{
		{
			name: "Full config file",
			fileContent: `title = "Example full config file"

interval = "8s"
proc_root = "/test_proc/"`,
			expected: Config{
				Interval: 8 * time.Second,
				ProcRoot: "/test_proc/",
			},
		},
		{
			name:        "Blank config file",
			fileContent: "",
			expected:    defaultConfig(),
		},
		{
			name:        "Only interval is filled",
			fileContent: "interval = '12s'",
			expected: Config{
				Interval: 12 * time.Second,
				ProcRoot: "/proc",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Load(configPath(t, tc.fileContent))

			require.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestLoadFileNotFound(t *testing.T) {
	result, err := Load(filepath.Join(t.TempDir(), configFileName))

	require.NoError(t, err)
	assert.Equal(t, defaultConfig(), result)
}

func TestLoadValidateError(t *testing.T) {
	cases := []struct {
		name          string
		fileContent   string
		expectedError error
	}{
		{
			name:          "Interval is less than allowed",
			fileContent:   "interval = '3ms'\nproc_root = '/proc'",
			expectedError: ErrInvalidInterval,
		},
		{
			name:          "Proc root is blank",
			fileContent:   "interval = '1s'\nproc_root = ''",
			expectedError: ErrBlankProcRoot,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Load(configPath(t, tc.fileContent))

			require.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, Config{}, result)
		})
	}
}

func TestLoadParseError(t *testing.T) {
	var parseError toml.ParseError

	result, err := Load(configPath(t, "interval : '12s'\nerror: df"))

	require.ErrorAs(t, err, &parseError)
	assert.Equal(t, Config{}, result)
}

func TestLoadReadError(t *testing.T) {
	path := filepath.Join(t.TempDir(), configFileName)
	require.NoError(t, os.Mkdir(path, 0o755))

	result, err := Load(path)

	require.Error(t, err)
	assert.Equal(t, Config{}, result)
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
