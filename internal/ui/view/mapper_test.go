package view

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tinkyrain/go-sysmon/internal/metrics"
	"github.com/tinkyrain/go-sysmon/internal/metrics/cpu"
	"github.com/tinkyrain/go-sysmon/internal/metrics/disk"
	"github.com/tinkyrain/go-sysmon/internal/metrics/memory"
	"github.com/tinkyrain/go-sysmon/internal/metrics/system"
)

func TestNewRenderDataFormatsSystem(t *testing.T) {
	t.Parallel()

	snapshot := metrics.Snapshot{
		Time: time.Date(2026, 11, 10, 14, 32, 5, 0, time.UTC),
		System: system.Stats{
			Info: system.Info{Hostname: "thinkpad", OS: "Arch Linux", Kernel: "6.12.1", Arch: "x86_64"},
			LoadAvg: system.LoadAvg{
				OneMin: 1.5, FiveMin: 0.75, FifteenMin: 0.2,
				RunningProcs: 3, TotalProcs: 412, LastPID: 90210,
			},
			Uptime: system.Uptime{Total: 509124.78, Idle: 3600},
		},
	}

	data := NewRenderData(snapshot)

	assert.Equal(t, "10.11.2026 14:32", data.Time)
	assert.Equal(t, int64(1794321125000000000), data.Timestamp)
	assert.Equal(t, []Field{
		{Name: "Hostname", Value: "thinkpad"},
		{Name: "OS", Value: "Arch Linux"},
		{Name: "Kernel", Value: "6.12.1"},
		{Name: "Arch", Value: "x86_64"},
	}, data.System.Info)
	assert.Equal(t, []Field{
		{Name: "One min", Value: "1.50"},
		{Name: "Five min", Value: "0.75"},
		{Name: "Fifteen min", Value: "0.20"},
		{Name: "Running procs", Value: "3"},
		{Name: "Total procs", Value: "412"},
		{Name: "Last run PID", Value: "90210"},
	}, data.System.LoadAvg)
	assert.Equal(t, Uptime{Total: "5d 21h 25m", Idle: "1h 0m"}, data.System.Uptime)
}

func TestNewRenderDataSortsCoresByLoad(t *testing.T) {
	t.Parallel()

	snapshot := metrics.Snapshot{CPU: cpu.Stats{Usage: cpu.Usage{
		Total: 37.9,
		Cores: []cpu.CoreUsage{
			{ID: "cpu0", Usage: 12},
			{ID: "cpu1", Usage: 88.6},
			{ID: "cpu2", Usage: 50},
		},
	}}}

	data := NewRenderData(snapshot)

	assert.Equal(t, Gauge{Label: "Total", Percent: 37}, data.CPU.Total)
	assert.Equal(t, []Gauge{
		{Label: "cpu1", Percent: 88},
		{Label: "cpu2", Percent: 50},
		{Label: "cpu0", Percent: 12},
	}, data.CPU.Cores)
}

func TestNewRenderDataKeepsEqualCoresInIDOrder(t *testing.T) {
	t.Parallel()

	snapshot := metrics.Snapshot{CPU: cpu.Stats{Usage: cpu.Usage{Cores: []cpu.CoreUsage{
		{ID: "cpu2"}, {ID: "cpu0"}, {ID: "cpu1"},
	}}}}

	data := NewRenderData(snapshot)

	require.Len(t, data.CPU.Cores, 3)
	assert.Equal(t, "cpu0", data.CPU.Cores[0].Label)
	assert.Equal(t, "cpu1", data.CPU.Cores[1].Label)
	assert.Equal(t, "cpu2", data.CPU.Cores[2].Label)
}

func TestNewRenderDataDoesNotTouchTheSnapshot(t *testing.T) {
	t.Parallel()

	cores := []cpu.CoreUsage{{ID: "cpu0", Usage: 1}, {ID: "cpu1", Usage: 99}}
	snapshot := metrics.Snapshot{CPU: cpu.Stats{Usage: cpu.Usage{Cores: cores}}}

	NewRenderData(snapshot)

	assert.Equal(t, "cpu0", cores[0].ID, "sorting must happen on a copy")
}

func TestNewRenderDataMemoryPercent(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name         string
		info         memory.Info
		expectedRAM  int
		expectedSwap int
	}{
		{
			name:         "half used",
			info:         memory.Info{Total: 16 * gib, Available: 8 * gib, SwapTotal: 8 * gib, SwapAvailable: 4 * gib},
			expectedRAM:  50,
			expectedSwap: 50,
		},
		{
			name:         "swap full reads as full, not as empty",
			info:         memory.Info{Total: 16 * gib, Available: 5 * gib, SwapTotal: 8 * gib, SwapAvailable: 0},
			expectedRAM:  68,
			expectedSwap: 100,
		},
		{
			name:         "no swap configured",
			info:         memory.Info{Total: 16 * gib, Available: 16 * gib},
			expectedRAM:  0,
			expectedSwap: 0,
		},
		{
			name:         "zero total does not divide by zero",
			info:         memory.Info{},
			expectedRAM:  0,
			expectedSwap: 0,
		},
		{
			name:         "available larger than total is ignored",
			info:         memory.Info{Total: 1 * gib, Available: 2 * gib},
			expectedRAM:  0,
			expectedSwap: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			data := NewRenderData(metrics.Snapshot{Memory: memory.Stats{Info: tc.info}})

			assert.Equal(t, tc.expectedRAM, data.Memory.RAM.Percent)
			assert.Equal(t, tc.expectedSwap, data.Memory.Swap.Percent)
			assert.Equal(t, "Memory", data.Memory.RAM.Label)
			assert.Equal(t, "SWAP", data.Memory.Swap.Label)
		})
	}
}

func TestNewRenderDataFormatsMemorySizes(t *testing.T) {
	t.Parallel()

	data := NewRenderData(metrics.Snapshot{Memory: memory.Stats{Info: memory.Info{
		Total: 16 * gib, Available: 5 * gib, SwapTotal: 8 * gib, SwapAvailable: 0,
	}}})

	assert.Equal(t, []Field{
		{
			Name:  "Total",
			Value: "16.00 GB",
		},
		{
			Name:  "Available",
			Value: "5.00 GB",
		},
		{
			Name:  "SWAP Total",
			Value: "8.00 GB",
		},
		{
			Name:  "SWAP Available",
			Value: "0.00 B",
		},
	}, data.Memory.Info)
}

func TestNewRenderDataMounts(t *testing.T) {
	t.Parallel()

	data := NewRenderData(metrics.Snapshot{Disk: disk.Stats{Mounts: []disk.Mount{
		{
			Path:      "/",
			Total:     512 * gib,
			Available: 128 * gib,
		},
		{
			Path:      "/home",
			Total:     1024 * gib,
			Available: 0,
		},
	}}})

	assert.Equal(t, []Mount{
		{Path: "/", Info: []Field{
			{Name: "Total", Value: "512.00 GB"},
			{Name: "Available", Value: "128.00 GB"},
		}},
		{Path: "/home", Info: []Field{
			{Name: "Total", Value: "1.00 TB"},
			{Name: "Available", Value: "0.00 B"},
		}},
	}, data.Disk.Mounts)
}

func TestNewRenderDataHandlesEmptySnapshot(t *testing.T) {
	t.Parallel()

	data := NewRenderData(metrics.Snapshot{})

	assert.Equal(t, "01.01.0001 00:00", data.Time)
	assert.NotZero(t, data.Timestamp)
	assert.Empty(t, data.CPU.Cores)
	assert.Empty(t, data.Disk.Mounts)
	assert.Equal(t, 0, data.CPU.Total.Percent)
	assert.Equal(t, Uptime{Total: "0s", Idle: "0s"}, data.System.Uptime)
	assert.Len(t, data.System.Info, 4)
	assert.Len(t, data.Memory.Info, 4)
}

func TestNewRenderDataTimestampGrowsWithSnapshotTime(t *testing.T) {
	t.Parallel()

	dataOne := NewRenderData(metrics.Snapshot{Time: time.Date(2026, 11, 10, 14, 32, 5, 0, time.UTC)})
	dataTwo := NewRenderData(metrics.Snapshot{Time: time.Date(2026, 11, 10, 14, 32, 5, 200, time.UTC)})

	assert.Greater(t, dataTwo.Timestamp, dataOne.Timestamp)
}
