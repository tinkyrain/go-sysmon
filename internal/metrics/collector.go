package metrics

import (
	"go-sysmon/internal/procfs"
	"time"
)

type Collector struct {
	interval    time.Duration
	fs          procfs.FileScanner
	prevCPUTick CPUTick
}

func New(procRoot string, interval time.Duration) *Collector {
	fs := procfs.New(procRoot)
	return &Collector{
		interval: interval,
		fs:       fs,
	}
}

func (c *Collector) Collect() (Snapshot, error) {
	memory, err := c.readMemory()
	if err != nil {
		return Snapshot{}, err
	}

	disks, err := c.readDisks()
	if err != nil {
		return Snapshot{}, err
	}

	cpuTick, err := c.readCPUTick()
	if err != nil {
		return Snapshot{}, err
	}

	c.prevCPUTick = cpuTick

	return Snapshot{
		Time:     time.Now(),
		Disks:    disks,
		Memory:   memory,
		CPUUsage: calculateCPUUsage(c.prevCPUTick, cpuTick),
	}, nil
}
