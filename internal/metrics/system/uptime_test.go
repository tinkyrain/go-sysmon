package system

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tinkyrain/go-sysmon/internal/procfs"
)

func uptimeFiles() map[string]string {
	return map[string]string{
		uptimeFile: "509124.78 2025143.51",
	}
}

func TestReadUptimeSuccess(t *testing.T) {
	expected := Uptime{
		Total: 509124.78,
		Idle:  2025143.51,
	}

	reader := uptimeReader{fs: procfs.New(tempDirWithFiles(t, uptimeFiles(), 0o755, 0o600))}

	result, err := reader.read()

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestReadUptimeFileError(t *testing.T) {
	dir := tempDirWithFiles(t, map[string]string{}, 0o755, 0o600)
	require.NoError(t, os.Mkdir(filepath.Join(dir, uptimeFile), 0o755))

	reader := uptimeReader{fs: procfs.New(dir)}

	result, err := reader.read()

	require.Error(t, err)
	assert.Equal(t, Uptime{}, result)
}

func TestReadUptimeBlankFileError(t *testing.T) {
	files := uptimeFiles()
	files[uptimeFile] = ""

	reader := uptimeReader{fs: procfs.New(tempDirWithFiles(t, files, 0o755, 0o600))}

	result, err := reader.read()

	require.Error(t, err, ErrInsufficientUptime)
	assert.Equal(t, Uptime{}, result)
}

func TestParseUptimeLineSuccess(t *testing.T) {
	expected := Uptime{
		Total: 509124.78,
		Idle:  2025143.51,
	}

	result, err := parseUptimeLine("509124.78 2025143.51")

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestParseUptimeLineIncorrectMetricsCountErr(t *testing.T) {
	result, err := parseUptimeLine("123")

	require.ErrorIs(t, err, ErrInsufficientUptime)
	require.Equal(t, Uptime{}, result)
}

func TestParseUptimeLineParseMetricsError(t *testing.T) {
	data := []string{
		"test 2025143.51",
		"509124.78 test",
	}

	for _, line := range data {
		t.Run(line, func(t *testing.T) {
			result, err := parseUptimeLine(line)

			require.ErrorIs(t, err, strconv.ErrSyntax)
			require.Equal(t, Uptime{}, result)
		})
	}
}
