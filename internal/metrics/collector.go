package metrics

import (
	"go-sysmon/internal/procfs"
	"time"
)

type Collector struct {
	fs           procfs.FileScanner
	memoryReader MemoryReader
	diskReader   DiskReader
	cpuReader    CPUReader
}

func New(procRoot string) *Collector {
	fs := procfs.New(procRoot)
	return &Collector{
		fs: fs,
		memoryReader: MemoryReader{
			fs: fs,
		},
		diskReader: DiskReader{
			fs: fs,
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

	cpuTick, err := c.cpuReader.Read()
	if err != nil {
		return Snapshot{}, err
	}

	return Snapshot{
		Time:     time.Now(),
		Disks:    disks,
		Memory:   memory,
		CPUUsage: cpuTick,
	}, nil
}
