package metrics

import (
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tinkyrain/go-sysmon/internal/metrics/cpu"
	"github.com/tinkyrain/go-sysmon/internal/metrics/disk"
	"github.com/tinkyrain/go-sysmon/internal/metrics/memory"
	"github.com/tinkyrain/go-sysmon/internal/metrics/system"
	"github.com/tinkyrain/go-sysmon/internal/statfs"
)

func metricsFiles() map[string]string {
	return map[string]string{
		"sys/kernel/hostname":  "hostname",
		"sys/kernel/osrelease": "osrelease",
		"sys/kernel/ostype":    "ostype",
		"loadavg":              "0.15 0.25 0.30 1/456 12345",
		"uptime":               "509124.78 2025143.51",
		"mounts":               "/dev/nvme0n1p5 / ext4 rw,relatime 0 0",
		"meminfo": `
MemTotal:       13003440 kB
MemFree:          199488 kB
MemAvailable:    5121912 kB
SwapTotal:       4194300 kB
SwapFree:        3383196 kB
Dirty:              1788 kB
		`,
		"stat": `
cpu  298830 3399 75896 7741115 5083 0 1586 0 0 0
cpu0 17496 72 6050 650357 437 0 376 0 0 0
cpu1 10757 57 3837 661024 159 0 1073 0 0 0
intr 30062225 138 0 0 0 0 0 0 1 0 7 0 0 0 0 0 0 0 0 0
`,
	}
}

func TestCollectSuccess(t *testing.T) {
	expected := Snapshot{
		Time: time.Now(),
		System: system.Stats{
			Info: system.Info{
				Hostname: "hostname",
				Kernel:   "osrelease",
				OS:       "ostype",
				Arch:     runtime.GOARCH,
			},
			LoadAvg: system.LoadAvg{
				OneMin:       0.15,
				FiveMin:      0.25,
				FifteenMin:   0.30,
				RunningProcs: 1,
				TotalProcs:   456,
				LastPID:      12345,
			},
			Uptime: system.Uptime{
				Total: 509124.78,
				Idle:  2025143.51,
			},
		},
		CPU: cpu.Stats{
			Usage: cpu.Usage{
				Total: 0,
				Cores: []cpu.CoreUsage{
					{ID: "cpu0", Usage: 0},
					{ID: "cpu1", Usage: 0},
				},
			},
		},
		Memory: memory.Stats{
			Memory: memory.MemInfo{
				Total:         13315522560,
				Available:     5244837888,
				SwapTotal:     4294963200,
				SwapAvailable: 3464392704,
			},
		},
		Disk: disk.Stats{
			Mounts: []disk.Mount{
				{
					Path:      "/",
					Total:     120,
					Available: 50,
				},
			},
		},
	}

	dir := tempDirWithFiles(t, metricsFiles())
	collector := New(dir, func(path string) (statfs.Stats, error) {
		result := statfs.Stats{}
		result.BlockSize = 10
		result.Blocks = 12
		result.Available = 5
		return result, nil
	})

	snapshot, err := collector.Collect()

	snapshot.Time = expected.Time

	require.NoError(t, err)
	assert.Equal(t, expected, snapshot)
}

func TestCollectError(t *testing.T) {
	cases := []struct {
		name       string
		targetFile string
	}{
		{
			name:       "File hostname is not exist",
			targetFile: "sys/kernel/hostname",
		},
		{
			name:       "File osrelease is not exist",
			targetFile: "sys/kernel/osrelease",
		},
		{
			name:       "File ostype is not exist",
			targetFile: "sys/kernel/ostype",
		},
		{
			name:       "File loadavg is not exist",
			targetFile: "loadavg",
		},
		{
			name:       "File uptime is not exist",
			targetFile: "uptime",
		},
		{
			name:       "File mounts is not exist",
			targetFile: "mounts",
		},
		{
			name:       "File meminfo is not exist",
			targetFile: "meminfo",
		},
		{
			name:       "File stat is not exist",
			targetFile: "stat",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			files := metricsFiles()
			delete(files, tc.targetFile)

			collector := New(tempDirWithFiles(t, files), func(path string) (statfs.Stats, error) {
				result := statfs.Stats{}
				result.BlockSize = 10
				result.Blocks = 12
				result.Available = 5
				return result, nil
			})

			snapshot, err := collector.Collect()

			require.ErrorIs(t, err, os.ErrNotExist)
			assert.Equal(t, Snapshot{}, snapshot)
		})
	}
}
