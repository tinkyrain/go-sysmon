package metrics

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tinkyrain/go-sysmon/internal/procfs"
)

const memoryFilename string = "meminfo"

type MemoryReader struct {
	fs procfs.FileScanner
}

type Memory struct {
	Total         uint64
	Available     uint64
	SwapTotal     uint64
	SwapAvailable uint64
}

func (r *MemoryReader) Read() (Memory, error) {
	collectData := Memory{}
	memoryMetrics := map[string]*uint64{
		"MemTotal":     &collectData.Total,
		"MemAvailable": &collectData.Available,
		"SwapTotal":    &collectData.SwapTotal,
		"SwapFree":     &collectData.SwapAvailable,
	}

	data, err := r.fs.ScanRows(memoryFilename)
	if err != nil {
		return collectData, err
	}

	for _, line := range data {
		name, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		if _, ok := memoryMetrics[name]; ok {
			fields := strings.Fields(value)
			if len(fields) == 0 {
				continue
			}
			metricValue, err := strconv.ParseUint(fields[0], 10, 64)
			if err != nil {
				return collectData, fmt.Errorf("error parsing metric %q value: %w", name, err)
			}
			*memoryMetrics[name] = metricValue * 1024 // Kb in bytes
		}
	}

	return collectData, nil
}
