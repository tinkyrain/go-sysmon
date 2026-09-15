package metrics

import (
	"fmt"
	"strconv"
	"strings"
)

const memoryFilename string = "meminfo"

type Memory struct {
	Total         uint64
	Available     uint64
	SwapTotal     uint64
	SwapAvailable uint64
}

func (c *Collector) readMemory() (Memory, error) {
	collectData := Memory{}
	memoryMetrics := map[string]*uint64{
		"MemTotal":     &collectData.Total,
		"MemAvailable": &collectData.Available,
		"SwapTotal":    &collectData.SwapTotal,
		"SwapFree":     &collectData.SwapAvailable,
	}

	data, err := c.fs.ScanRows(memoryFilename)
	if err != nil {
		return collectData, err
	}

	for _, line := range data {
		name, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		if _, ok := memoryMetrics[name]; ok {
			metricValue, err := strconv.ParseUint(strings.Fields(value)[0], 10, 64)
			if err != nil {
				return collectData, fmt.Errorf("error parsing metric '%q' value: %w", name, err)
			}
			*memoryMetrics[name] = metricValue
		}
	}

	return collectData, nil
}
