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

func (c *Collector) Collect() (*Snapshot, error) {
	snapshot := &Snapshot{}

	memory, err := c.readMemory()
	if err != nil {
		return nil, err
	}

	disks, err := c.readDisks()
	if err != nil {
		return nil, err
	}

	cpuTick, err := c.readCPUTick()
	if err != nil {
		return nil, err
	}

	snapshot.Time = time.Now()
	snapshot.Disks = disks
	snapshot.Memory = memory
	snapshot.CPUUsage = calculateCPUUsage(c.prevCPUTick, cpuTick)

	c.prevCPUTick = cpuTick

	return snapshot, nil
}
