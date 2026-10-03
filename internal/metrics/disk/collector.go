package disk

import (
	"syscall"

	"github.com/tinkyrain/go-sysmon/internal/procfs"
)

type Collector struct {
	mountsReader mountsReader
}

func New(
	fs procfs.FileScanner,
	statfsFunc func(string) (syscall.Statfs_t, error),
) Collector {
	return Collector{
		mountsReader: mountsReader{
			fs:         fs,
			statfsFunc: statfsFunc,
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
