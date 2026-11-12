package metrics

import (
	"fmt"
	"time"

	"github.com/tinkyrain/go-sysmon/internal/metrics/cpu"
	"github.com/tinkyrain/go-sysmon/internal/metrics/disk"
	"github.com/tinkyrain/go-sysmon/internal/metrics/memory"
	"github.com/tinkyrain/go-sysmon/internal/metrics/system"
	"github.com/tinkyrain/go-sysmon/internal/procfs"
	"github.com/tinkyrain/go-sysmon/internal/statfs"
)

type Collector struct {
	system system.Collector
	cpu    *cpu.Collector
	memory memory.Collector
	disk   disk.Collector
}

func New(procRoot string, statfsFunc statfs.Func) *Collector {
	fs := procfs.New(procRoot)

	return &Collector{
		system: system.New(fs),
		cpu:    cpu.New(fs),
		memory: memory.New(fs),
		disk:   disk.New(fs, statfsFunc),
	}
}

func (c *Collector) Collect() (Snapshot, error) {
	systemStats, err := c.system.Collect()
	if err != nil {
		return Snapshot{}, fmt.Errorf("collecting system metrics: %w", err)
	}

	cpuStats, err := c.cpu.Collect()
	if err != nil {
		return Snapshot{}, fmt.Errorf("collecting cpu metrics: %w", err)
	}

	memoryStats, err := c.memory.Collect()
	if err != nil {
		return Snapshot{}, fmt.Errorf("collecting memory metrics: %w", err)
	}

	diskStats, err := c.disk.Collect()
	if err != nil {
		return Snapshot{}, fmt.Errorf("collecting disk metrics: %w", err)
	}

	return Snapshot{
		Time:   time.Now(),
		System: systemStats,
		CPU:    cpuStats,
		Memory: memoryStats,
		Disk:   diskStats,
	}, nil
}
