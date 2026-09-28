package system

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tinkyrain/go-sysmon/internal/procfs"
)

func uptimeFiles() map[string]string {
	return map[string]string{
		"uptime": "509124.78 2025143.51",
	}
}

func TestReadUptimeSuccess(t *testing.T) {
	expected := Uptime{
		Total: 509124.78,
		Idle:  2025143.51,
	}

	fs := procfs.New(TempDirWithFiles(t, uptimeFiles(), 0o755, 0o600))
	reader := UptimeReader{fs: fs}

	result, err := reader.Read()

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestReadUptimeFileError(t *testing.T) {
	expected := Uptime{}

	filepath := TempDirWithFiles(t, uptimeFiles(), 0, 0)
	fs := procfs.New(filepath)
	reader := UptimeReader{fs: fs}

	result, err := reader.Read()

	require.Error(t, err)
	assert.Equal(t, expected, result)
}

func TestReadUptimeBlankFileError(t *testing.T) {
	expected := Uptime{}

	files := uptimeFiles()
	files["uptime"] = ""
	filepath := TempDirWithFiles(t, files, 0o600, 0o600)
	fs := procfs.New(filepath)
	reader := UptimeReader{fs: fs}

	result, err := reader.Read()

	require.ErrorIs(t, err, ErrInsufficientUptime)
	assert.Equal(t, expected, result)
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

			require.Error(t, err)
			require.Equal(t, Uptime{}, result)
		})
	}
}
