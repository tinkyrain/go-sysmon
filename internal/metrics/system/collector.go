package system

import "github.com/tinkyrain/go-sysmon/internal/procfs"

type Collector struct {
	infoReader    infoReader
	loadAvgReader loadAvgReader
	uptimeReader  uptimeReader
}

func New(fs procfs.FileScanner) Collector {
	return Collector{
		infoReader:    infoReader{fs: fs},
		loadAvgReader: loadAvgReader{fs: fs},
		uptimeReader:  uptimeReader{fs: fs},
	}
}

func (c Collector) Collect() (Stats, error) {
	info, err := c.infoReader.read()
	if err != nil {
		return Stats{}, err
	}

	loadAvg, err := c.loadAvgReader.read()
	if err != nil {
		return Stats{}, err
	}

	uptime, err := c.uptimeReader.read()
	if err != nil {
		return Stats{}, err
	}

	return Stats{
		Info:    info,
		LoadAvg: loadAvg,
		Uptime:  uptime,
	}, nil
}
