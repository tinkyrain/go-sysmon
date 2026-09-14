package metrics

import (
	"go-sysmon/internal/procfs"
	"time"
)

type Collector struct {
	interval time.Duration
	fs       procfs.FileScanner
}

func New(procRoot string, interval time.Duration) *Collector {
	fs := procfs.New(procRoot)
	return &Collector{
		interval: interval,
		fs:       fs,
	}
}

func (c *Collector) Collect() (Snapshot, error) {
	var snapshot = Snapshot{}
	snapshot.Time = time.Now()

	memory, err := c.readMemory()
	if err != nil {
		return snapshot, err
	}
	disks, err := c.readDisks()
	if err != nil {
		return snapshot, err
	}

	snapshot.Disks = disks
	snapshot.Memory = memory

	return snapshot, nil
}
