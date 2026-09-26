package metrics

import (
	"time"

	"github.com/tinkyrain/go-sysmon/internal/procfs"
	"github.com/tinkyrain/go-sysmon/internal/statfs"
)

type Collector struct {
	memoryReader MemoryReader
	diskReader   DiskReader
	cpuReader    CPUReader
}

func New(procRoot string) *Collector {
	fs := procfs.New(procRoot)
	return &Collector{
		memoryReader: MemoryReader{
			fs: fs,
		},
		diskReader: DiskReader{
			fs:         fs,
			statfsFunc: statfs.GetDirStatfs,
		},
		cpuReader: CPUReader{
			fs: fs,
		},
	}
}

func (c *Collector) Collect() (Snapshot, error) {
	memory, err := c.memoryReader.Read()
	if err != nil {
		return Snapshot{}, err
	}

	disks, err := c.diskReader.Read()
	if err != nil {
		return Snapshot{}, err
	}

	cpuUsages, err := c.cpuReader.Read()
	if err != nil {
		return Snapshot{}, err
	}

	return Snapshot{
		Time:      time.Now(),
		Disks:     disks,
		Memory:    memory,
		CPUUsages: cpuUsages,
	}, nil
}
