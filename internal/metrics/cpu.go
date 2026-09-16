package metrics

import (
	"errors"
	"fmt"
	"go-sysmon/internal/procfs"
	"strconv"
	"strings"
)

type CPUReader struct {
	fs          procfs.FileScanner
	prevCPUTick CPUTick
}

type CPUTick struct {
	Total   uint64
	Idle    uint64
	NonIdle uint64
}

const statFilename = "stat"
const cpuParamName = "cpu"

func (r *CPUReader) Read() (float64, error) {
	data, err := r.fs.ScanRows(statFilename)
	if err != nil {
		return 0, err
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
		return 0, errors.New("not found cpu row")
	}

	tickMetrics, err := parseCPUTick(cpuTickMetricsRows)
	if err != nil {
		return 0, err
	}

	idle := tickMetrics["idle"] + tickMetrics["iowait"]
	nonIdle := tickMetrics["user"] + tickMetrics["nice"] + tickMetrics["system"] + tickMetrics["irq"] +
		tickMetrics["softirq"] + tickMetrics["steal"]
	total := idle + nonIdle

	var tick CPUTick

	tick.Total = total
	tick.Idle = idle
	tick.NonIdle = nonIdle

	cpuUsage := calculateCPUUsage(r.prevCPUTick, tick)
	r.prevCPUTick = tick

	return cpuUsage, nil
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

func parseCPUTick(metrics []string) (map[string]uint64, error) {
	if len(metrics) < 10 {
		return map[string]uint64{}, errors.New("parsing CPU ticks need 10 metrics")
	}
	converted := make([]uint64, 0, len(metrics))
	for _, value := range metrics {
		convertedValue, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return map[string]uint64{}, fmt.Errorf("parsing cpu metric %q: %w", value, err)
		}
		converted = append(converted, convertedValue)
	}

	return map[string]uint64{
		"user":      converted[0],
		"nice":      converted[1],
		"system":    converted[2],
		"idle":      converted[3],
		"iowait":    converted[4],
		"irq":       converted[5],
		"softirq":   converted[6],
		"steal":     converted[7],
		"guest":     converted[8],
		"guestNice": converted[9],
	}, nil
}
