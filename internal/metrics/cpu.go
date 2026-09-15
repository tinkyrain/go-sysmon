package metrics

import (
	"fmt"
	"strconv"
	"strings"
)

type CPUTick struct {
	Total   uint
	Idle    uint
	NonIdle uint
}

type CPUTickMetrics struct {
	User      uint
	Nice      uint
	System    uint
	Idle      uint
	Iowait    uint
	Irq       uint
	Softirq   uint
	Steal     uint
	Guest     uint
	GuestNice uint
}

const statFilename = "stat"
const cpuParamName = "cpu"

func (c *Collector) readCPUTick() (CPUTick, error) {
	var tick CPUTick
	data, err := c.fs.ScanRows(statFilename)
	if err != nil {
		return tick, err
	}

	for _, line := range data {
		fields := strings.Fields(line)
		// 11 because cpu has 11 metrics
		if len(fields) < 11 {
			continue
		}

		if fields[0] != cpuParamName {
			continue
		}

		tickMetrics, err := convertArrayMetricsToStruct(fields[1:])
		if err != nil {
			return tick, err
		}

		Idle := tickMetrics.Idle + tickMetrics.Iowait
		NonIdle := tickMetrics.User + tickMetrics.Nice + tickMetrics.System + tickMetrics.Irq +
			tickMetrics.Softirq + tickMetrics.Steal
		Total := Idle + NonIdle

		tick = CPUTick{
			Total:   Total,
			Idle:    Idle,
			NonIdle: NonIdle,
		}
		break
	}

	return tick, nil
}

func calculateCPUUsage(prev, cur CPUTick) float64 {
	var percentUsage float64

	idleDelta := cur.Idle - prev.Idle
	nonIdleDelta := cur.NonIdle - prev.NonIdle
	totalDelta := idleDelta + nonIdleDelta

	if nonIdleDelta > 0 && totalDelta > 0 {
		percentUsage = float64(nonIdleDelta) / float64(totalDelta) * 100
	}

	return percentUsage
}

func convertArrayMetricsToStruct(metrics []string) (CPUTickMetrics, error) {
	converted := make([]uint, len(metrics))
	for _, value := range metrics {
		convertedValue, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return CPUTickMetrics{}, fmt.Errorf("could not convert %s to uint", value)
		}
		converted = append(converted, uint(convertedValue))
	}
	return CPUTickMetrics{
		User:      converted[1],
		Nice:      converted[2],
		System:    converted[3],
		Idle:      converted[4],
		Iowait:    converted[5],
		Irq:       converted[6],
		Softirq:   converted[7],
		Steal:     converted[8],
		Guest:     converted[9],
		GuestNice: converted[10],
	}, nil
}
