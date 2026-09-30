package cpu

import "github.com/tinkyrain/go-sysmon/internal/procfs"

type Collector struct {
	usageReader *usageReader
}

func New(fs procfs.FileScanner) *Collector {
	return &Collector{
		usageReader: &usageReader{
			fs: fs,
		},
	}
}

func (c *Collector) Collect() (Stats, error) {
	usage, err := c.usageReader.read()
	if err != nil {
		return Stats{}, err
	}

	return Stats{
		Usage: usage,
	}, nil
}
