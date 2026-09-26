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
	memory := Memory{}
	metrics := map[string]*uint64{
		"MemTotal":     &memory.Total,
		"MemAvailable": &memory.Available,
		"SwapTotal":    &memory.SwapTotal,
		"SwapFree":     &memory.SwapAvailable,
	}

	data, err := r.fs.ScanRows(memoryFilename)
	if err != nil {
		return Memory{}, err
	}

	for _, line := range data {
		name, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		if _, ok := metrics[name]; ok {
			fields := strings.Fields(value)
			if len(fields) == 0 {
				continue
			}
			metricValue, err := strconv.ParseUint(fields[0], 10, 64)
			if err != nil {
				return Memory{}, fmt.Errorf("error parsing metric %q value: %w", name, err)
			}
			*metrics[name] = metricValue * 1024 // Kb in bytes
		}
	}

	return memory, nil
}
