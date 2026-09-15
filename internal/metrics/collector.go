package metrics

import (
	"go-sysmon/internal/procfs"
	"time"
)

type Collector struct {
	interval     time.Duration
	fs           procfs.FileScanner
	prevCPUTicks map[string]CPUTick
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

	cpu, err := c.readCPU()
	if err != nil {
		return nil, err
	}
	var cpuUsagePercent []CPU
	if len(c.prevCPUTicks) == 0 {
		for core, _ := range cpu {
			cpuUsagePercent = append(cpuUsagePercent, CPU{Title: core, UsagePercent: 0})
		}
	} else {
		cpuUsagePercent = calculateCPUUsage(c.prevCPUTicks, cpu)
	}
	c.prevCPUTicks = cpu

	snapshot.Time = time.Now()
	snapshot.Disks = disks
	snapshot.Memory = memory
	snapshot.CPUUsage = cpuUsagePercent

	return snapshot, nil
}
