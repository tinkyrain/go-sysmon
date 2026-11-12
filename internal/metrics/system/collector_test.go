package system

import (
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tinkyrain/go-sysmon/internal/procfs"
)

func systemFiles() map[string]string {
	return map[string]string{
		hostnameFile:  "hostname",
		osreleaseFile: "osrelease",
		ostypeFile:    "ostype",
		loadAvgFile:   "0.15 0.25 0.30 1/456 12345",
		uptimeFile:    "509124.78 2025143.51",
	}
}

func TestCollectSuccess(t *testing.T) {
	expected := Stats{
		Info: Info{
			Hostname: "hostname",
			Kernel:   "osrelease",
			OS:       "ostype",
			Arch:     runtime.GOARCH,
		},
		LoadAvg: LoadAvg{
			OneMin:       0.15,
			FiveMin:      0.25,
			FifteenMin:   0.30,
			RunningProcs: 1,
			TotalProcs:   456,
			LastPID:      12345,
		},
		Uptime: Uptime{
			Total: 509124.78,
			Idle:  2025143.51,
		},
	}

	collector := New(procfs.New(tempDirWithFiles(t, systemFiles())))

	result, err := collector.Collect()

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestCollectError(t *testing.T) {
	cases := []struct {
		name          string
		targetFile    string
		expectedError error
	}{
		{
			name:          "Hostname file not found",
			targetFile:    hostnameFile,
			expectedError: os.ErrNotExist,
		},
		{
			name:          "Osrelease file not found",
			targetFile:    osreleaseFile,
			expectedError: os.ErrNotExist,
		},
		{
			name:          "Ostype file not found",
			targetFile:    ostypeFile,
			expectedError: os.ErrNotExist,
		},
		{
			name:          "Loadavg file not found",
			targetFile:    loadAvgFile,
			expectedError: os.ErrNotExist,
		},
		{
			name:          "Uptime file not found",
			targetFile:    uptimeFile,
			expectedError: os.ErrNotExist,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			files := systemFiles()
			delete(files, tc.targetFile)

			collector := New(procfs.New(tempDirWithFiles(t, files)))

			result, err := collector.Collect()

			require.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, Stats{}, result)
		})
	}
}
