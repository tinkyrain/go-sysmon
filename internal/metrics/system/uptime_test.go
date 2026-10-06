package system

import (
	"os"
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

func TestReadUptimeScanError(t *testing.T) {
	files := uptimeFiles()
	delete(files, uptimeFile)

	reader := uptimeReader{fs: procfs.New(tempDirWithFiles(t, files, 0o755, 0o600))}

	result, err := reader.read()

	require.ErrorIs(t, err, os.ErrNotExist)
	assert.Equal(t, Uptime{}, result)
}

func TestReadUptimeEmptyDataError(t *testing.T) {
	files := uptimeFiles()
	files[uptimeFile] = ""

	reader := uptimeReader{fs: procfs.New(tempDirWithFiles(t, files, 0o755, 0o600))}

	result, err := reader.read()

	require.ErrorIs(t, err, ErrEmptyUptimeData)
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

func TestParseUptimeLineError(t *testing.T) {
	cases := []struct {
		name          string
		line          string
		expectedError error
	}{
		{
			name:          "Incorrect data in the line",
			line:          "2025143.5",
			expectedError: ErrIncorrectUptimeData,
		},
		{
			name:          "Total uptime metric is incorrect",
			line:          "test 509124.78",
			expectedError: strconv.ErrSyntax,
		},
		{
			name:          "Idle uptime metric is incorrect",
			line:          "509124.78 test",
			expectedError: strconv.ErrSyntax,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := parseUptimeLine(tc.line)

			require.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, Uptime{}, result)
		})
	}
}
