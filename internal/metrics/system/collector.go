package system

import "github.com/tinkyrain/go-sysmon/internal/procfs"

type Collector struct {
	infoReader    InfoReader
	loadAvgReader LoadAvgReader
	uptimeReader  UptimeReader
}

func New(fs procfs.FileScanner) Collector {
	return Collector{
		infoReader:    InfoReader{fs: fs},
		loadAvgReader: LoadAvgReader{fs: fs},
		uptimeReader:  UptimeReader{fs: fs},
	}
}

func (c *Collector) Collect() (Stats, error) {
	info, err := c.infoReader.Read()
	if err != nil {
		return Stats{}, err
	}

	loadAvg, err := c.loadAvgReader.Read()
	if err != nil {
		return Stats{}, err
	}

	uptime, err := c.uptimeReader.Read()
	if err != nil {
		return Stats{}, err
	}

	return Stats{
		Info:    info,
		LoadAvg: loadAvg,
		Uptime:  uptime,
	}, nil
}
