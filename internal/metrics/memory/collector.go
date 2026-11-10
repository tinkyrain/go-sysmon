package memory

import "github.com/tinkyrain/go-sysmon/internal/procfs"

type Collector struct {
	meminfoReader meminfoReader
}

func New(fs procfs.Reader) Collector {
	return Collector{
		meminfoReader: meminfoReader{
			fs: fs,
		},
	}
}

func (c Collector) Collect() (Stats, error) {
	info, err := c.meminfoReader.read()
	if err != nil {
		return Stats{}, err
	}

	return Stats{
		Info: info,
	}, nil
}
