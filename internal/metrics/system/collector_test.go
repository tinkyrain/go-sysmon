package system

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tinkyrain/go-sysmon/internal/procfs"
)

func systemFiles() map[string]string {
	return map[string]string{
		"sys/kernel/hostname":  "hostname",
		"sys/kernel/osrelease": "osrelease",
		"sys/kernel/ostype":    "ostype",
		"loadavg":              "0.15 0.25 0.30 1/456 12345",
		"uptime":               "509124.78 2025143.51",
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

	fs := procfs.New(TempDirWithFiles(t, systemFiles(), 0o755, 0o600))
	collector := New(fs)

	result, err := collector.Collect()

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}
