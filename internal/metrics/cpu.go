package metrics

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type CPUTick struct {
	Total   uint64
	Idle    uint64
	NonIdle uint64
}

type CPUTickMetrics struct {
	User      uint64
	Nice      uint64
	System    uint64
	Idle      uint64
	Iowait    uint64
	Irq       uint64
	Softirq   uint64
	Steal     uint64
	Guest     uint64
	GuestNice uint64
}

const statFilename = "stat"
const cpuParamName = "cpu"

func (c *Collector) readCPUTick() (CPUTick, error) {
	var tick CPUTick
	data, err := c.fs.ScanRows(statFilename)
	if err != nil {
		return tick, err
	}

	var cpuTickMetricsRows []string

	for _, line := range data {
		fields := strings.Fields(line)
		// 11 because cpu has 11 metrics
		if len(fields) < 11 {
			continue
		}
		if fields[0] != cpuParamName {
			continue
		}
		cpuTickMetricsRows = fields[1:]
		break
	}

	if len(cpuTickMetricsRows) == 0 {
		return tick, errors.New("not found cpu row")
	}

	tickMetrics, err := parseCPUTick(cpuTickMetricsRows)
	if err != nil {
		return tick, err
	}

	idle := tickMetrics.Idle + tickMetrics.Iowait
	nonIdle := tickMetrics.User + tickMetrics.Nice + tickMetrics.System + tickMetrics.Irq +
		tickMetrics.Softirq + tickMetrics.Steal
	total := idle + nonIdle

	tick.Total = total
	tick.Idle = idle
	tick.NonIdle = nonIdle

	return tick, nil
}

func calculateCPUUsage(prev, cur CPUTick) float64 {
	var percentUsage float64

	if prev.Total == 0 || cur.Idle < prev.Idle || cur.NonIdle < prev.NonIdle {
		return 0
	}

	idleDelta := cur.Idle - prev.Idle
	nonIdleDelta := cur.NonIdle - prev.NonIdle
	totalDelta := idleDelta + nonIdleDelta

	if nonIdleDelta > 0 && totalDelta > 0 {
		percentUsage = float64(nonIdleDelta) / float64(totalDelta) * 100
	}

	return percentUsage
}

func parseCPUTick(metrics []string) (CPUTickMetrics, error) {
	if len(metrics) < 10 {
		return CPUTickMetrics{}, errors.New("parsing CPU ticks need 10 metrics")
	}
	converted := make([]uint64, 0, len(metrics))
	for _, value := range metrics {
		convertedValue, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return CPUTickMetrics{}, fmt.Errorf("parsing cpu metric %q: %w", value, err)
		}
		converted = append(converted, convertedValue)
	}
	return CPUTickMetrics{
		User:      converted[0],
		Nice:      converted[1],
		System:    converted[2],
		Idle:      converted[3],
		Iowait:    converted[4],
		Irq:       converted[5],
		Softirq:   converted[6],
		Steal:     converted[7],
		Guest:     converted[8],
		GuestNice: converted[9],
	}, nil
}
