package disk

import (
	"github.com/tinkyrain/go-sysmon/internal/procfs"
)

type Collector struct {
	mountsReader mountsReader
}

func New(
	fs procfs.FileScanner,
	statfs statfsFunc,
) Collector {
	return Collector{
		mountsReader: mountsReader{
			fs:     fs,
			statfs: statfs,
		},
	}
}

func (c Collector) Collect() (Stats, error) {
	mounts, err := c.mountsReader.read()
	if err != nil {
		return Stats{}, err
	}

	return Stats{
		Mounts: mounts,
	}, nil
}
