package view_test

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
	"github.com/tinkyrain/go-sysmon/internal/ui/view"
)

const (
	kib = 1024
	mib = 1024 * kib
	gib = 1024 * mib
)

func TestHumanizeSize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		size float64
		want string
	}{
		{name: "zero", size: 0, want: "0.00 B"},
		{name: "bytes stay bytes", size: 1023, want: "1023.00 B"},
		{name: "exact kilobyte", size: kib, want: "1.00 KB"},
		{name: "fraction is kept", size: 1536, want: "1.50 KB"},
		{name: "megabytes", size: 512 * mib, want: "512.00 MB"},
		{name: "gigabytes match free -h", size: 16 * gib, want: "16.00 GB"},
		{name: "fraction below the next unit", size: 2*gib - 1, want: "2.00 GB"},
		{name: "exabytes", size: 1 << 60, want: "1.00 EB"},
		{name: "beyond the largest unit it stops dividing", size: 1 << 70, want: "1024.00 EB"},
		{name: "negative is not a size", size: -1, want: "n/a"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, view.HumanizeSize(tt.size))
		})
	}
}

func TestHumanizeDuration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		seconds float64
		want    string
	}{
		{name: "zero", seconds: 0, want: "0s"},
		{name: "seconds only", seconds: 42.9, want: "42s"},
		{name: "minutes and seconds", seconds: 125, want: "2m 5s"},
		{name: "hours drop the seconds", seconds: 3*3600 + 4*60 + 59, want: "3h 4m"},
		{name: "days", seconds: 5*86400 + 21*3600 + 25*60, want: "5d 21h 25m"},
		{name: "day boundary keeps zeroes", seconds: 86400, want: "1d 0h 0m"},
		{name: "negative is not a duration", seconds: -1, want: "n/a"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, view.HumanizeDuration(tt.seconds))
		})
	}
}

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

	data := view.NewRenderData(snapshot)

	assert.Equal(t, "10.11.2026 14:32", data.Time)
	assert.Equal(t, []view.Field{
		{Name: "Hostname", Value: "thinkpad"},
		{Name: "OS", Value: "Arch Linux"},
		{Name: "Kernel", Value: "6.12.1"},
		{Name: "Arch", Value: "x86_64"},
	}, data.System.Info)
	assert.Equal(t, []view.Field{
		{Name: "One min", Value: "1.50"},
		{Name: "Five min", Value: "0.75"},
		{Name: "Fifteen min", Value: "0.20"},
		{Name: "Running procs", Value: "3"},
		{Name: "Total procs", Value: "412"},
		{Name: "Last run PID", Value: "90210"},
	}, data.System.LoadAvg)
	assert.Equal(t, view.Uptime{Total: "5d 21h 25m", Idle: "1h 0m"}, data.System.Uptime)
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

	data := view.NewRenderData(snapshot)

	assert.Equal(t, view.Gauge{Label: "Total", Percent: 37}, data.CPU.Total)
	assert.Equal(t, []view.Gauge{
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

	data := view.NewRenderData(snapshot)

	require.Len(t, data.CPU.Cores, 3)
	assert.Equal(t, "cpu0", data.CPU.Cores[0].Label)
	assert.Equal(t, "cpu1", data.CPU.Cores[1].Label)
	assert.Equal(t, "cpu2", data.CPU.Cores[2].Label)
}

func TestNewRenderDataDoesNotTouchTheSnapshot(t *testing.T) {
	t.Parallel()

	cores := []cpu.CoreUsage{{ID: "cpu0", Usage: 1}, {ID: "cpu1", Usage: 99}}
	snapshot := metrics.Snapshot{CPU: cpu.Stats{Usage: cpu.Usage{Cores: cores}}}

	view.NewRenderData(snapshot)

	assert.Equal(t, "cpu0", cores[0].ID, "sorting must happen on a copy")
}

func TestNewRenderDataMemoryPercent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		info     memory.Info
		wantRAM  int
		wantSwap int
	}{
		{
			name:     "half used",
			info:     memory.Info{Total: 16 * gib, Available: 8 * gib, SwapTotal: 8 * gib, SwapAvailable: 4 * gib},
			wantRAM:  50,
			wantSwap: 50,
		},
		{
			name:     "swap full reads as full, not as empty",
			info:     memory.Info{Total: 16 * gib, Available: 5 * gib, SwapTotal: 8 * gib, SwapAvailable: 0},
			wantRAM:  68,
			wantSwap: 100,
		},
		{
			name:     "no swap configured",
			info:     memory.Info{Total: 16 * gib, Available: 16 * gib},
			wantRAM:  0,
			wantSwap: 0,
		},
		{
			name:     "zero total does not divide by zero",
			info:     memory.Info{},
			wantRAM:  0,
			wantSwap: 0,
		},
		{
			name:     "available larger than total is ignored",
			info:     memory.Info{Total: 1 * gib, Available: 2 * gib},
			wantRAM:  0,
			wantSwap: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			data := view.NewRenderData(metrics.Snapshot{Memory: memory.Stats{Info: tt.info}})

			assert.Equal(t, tt.wantRAM, data.Memory.RAM.Percent)
			assert.Equal(t, tt.wantSwap, data.Memory.Swap.Percent)
			assert.Equal(t, "Memory", data.Memory.RAM.Label)
			assert.Equal(t, "SWAP", data.Memory.Swap.Label)
		})
	}
}

func TestNewRenderDataFormatsMemorySizes(t *testing.T) {
	t.Parallel()

	data := view.NewRenderData(metrics.Snapshot{Memory: memory.Stats{Info: memory.Info{
		Total: 16 * gib, Available: 5 * gib, SwapTotal: 8 * gib, SwapAvailable: 0,
	}}})

	assert.Equal(t, []view.Field{
		{Name: "Total", Value: "16.00 GB"},
		{Name: "Available", Value: "5.00 GB"},
		{Name: "SWAP Total", Value: "8.00 GB"},
		{Name: "SWAP Available", Value: "0.00 B"},
	}, data.Memory.Info)
}

func TestNewRenderDataMounts(t *testing.T) {
	t.Parallel()

	data := view.NewRenderData(metrics.Snapshot{Disk: disk.Stats{Mounts: []disk.Mount{
		{Path: "/", Total: 512 * gib, Available: 128 * gib},
		{Path: "/home", Total: 1024 * gib, Available: 0},
	}}})

	assert.Equal(t, []view.Mount{
		{Path: "/", Info: []view.Field{
			{Name: "Total", Value: "512.00 GB"},
			{Name: "Available", Value: "128.00 GB"},
		}},
		{Path: "/home", Info: []view.Field{
			{Name: "Total", Value: "1.00 TB"},
			{Name: "Available", Value: "0.00 B"},
		}},
	}, data.Disk.Mounts)
}

func TestNewRenderDataHandlesEmptySnapshot(t *testing.T) {
	t.Parallel()

	data := view.NewRenderData(metrics.Snapshot{})

	assert.Equal(t, "01.01.0001 00:00", data.Time)
	assert.Empty(t, data.CPU.Cores)
	assert.Empty(t, data.Disk.Mounts)
	assert.Equal(t, 0, data.CPU.Total.Percent)
	assert.Equal(t, view.Uptime{Total: "0s", Idle: "0s"}, data.System.Uptime)
	assert.Len(t, data.System.Info, 4)
	assert.Len(t, data.Memory.Info, 4)
}
