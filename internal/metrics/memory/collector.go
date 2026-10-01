package memory

import "github.com/tinkyrain/go-sysmon/internal/procfs"

type Collector struct {
	memoryReader memoryReader
}

func New(fs procfs.FileScanner) Collector {
	return Collector{
		memoryReader: memoryReader{
			fs: fs,
		},
	}
}

func (c *Collector) Collect() (Stats, error) {
	memory, err := c.memoryReader.read()
	if err != nil {
		return Stats{}, err
	}

	return Stats{
		Memory: memory,
	}, nil
}
