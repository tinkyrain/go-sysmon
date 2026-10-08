package memory

import "github.com/tinkyrain/go-sysmon/internal/procfs"

type Collector struct {
	memoryReader meminfoReader
}

func New(fs procfs.FileScanner) Collector {
	return Collector{
		memoryReader: meminfoReader{
			fs: fs,
		},
	}
}

func (c Collector) Collect() (Stats, error) {
	memory, err := c.memoryReader.read()
	if err != nil {
		return Stats{}, err
	}

	return Stats{
		Memory: memory,
	}, nil
}
