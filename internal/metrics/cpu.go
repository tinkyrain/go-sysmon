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

type CPU struct {
	Title        string
	UsagePercent float64
}

const filenameStat = "stat"
const cpuParamName = "cpu"

func (c *Collector) readCPU() (map[string]CPUTick, error) {
	var result map[string]CPUTick = make(map[string]CPUTick)
	data, err := c.fs.ScanRows(filenameStat)
	if err != nil {
		return result, err
	}

	for _, line := range data {
		fields := strings.Fields(line)
		// 11 because cpu has 11 metrics
		if len(fields) < 11 {
			continue
		}

		if !strings.HasPrefix(fields[0], cpuParamName) {
			continue
		}
		key := cpuParamName
		if len(fields[0]) > len(cpuParamName) {
			key = fields[0][len(cpuParamName):]
		}

		cpuMetrics, err := convertStringMetricsToUint(map[string]string{
			"user":       fields[1],
			"nice":       fields[2],
			"system":     fields[3],
			"idle":       fields[4],
			"iowait":     fields[5],
			"irq":        fields[6],
			"softirq":    fields[7],
			"steal":      fields[8],
			"guest":      fields[9],
			"guest_nice": fields[10],
		})

		if err != nil {
			return result, err
		}

		Idle := cpuMetrics["idle"] + cpuMetrics["iowait"]
		NonIdle := cpuMetrics["user"] + cpuMetrics["nice"] + cpuMetrics["system"] + cpuMetrics["irq"] +
			cpuMetrics["softirq"] + cpuMetrics["steal"]
		Total := Idle + NonIdle

		result[key] = CPUTick{
			Total:   Total,
			Idle:    Idle,
			NonIdle: NonIdle,
		}
	}

	return result, nil
}

func calculateCPUUsage(prev, cur map[string]CPUTick) []CPU {
	var result []CPU
	for core, curTick := range cur {
		cpu := CPU{
			Title:        core,
			UsagePercent: 0,
		}
		if prevTick, ok := prev[core]; ok {
			idleDelta := curTick.Idle - prevTick.Idle
			nonIdleDelta := curTick.NonIdle - prevTick.NonIdle
			totalDelta := idleDelta + nonIdleDelta

			if nonIdleDelta > 0 && totalDelta > 0 {
				cpu.UsagePercent = float64(nonIdleDelta) / float64(totalDelta) * 100
			}
		}
		result = append(result, cpu)
	}
	return result
}

func convertStringMetricsToUint(metrics map[string]string) (map[string]uint, error) {
	converted := make(map[string]uint, len(metrics))
	for key, value := range metrics {
		convertedValue, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return converted, fmt.Errorf("could not convert %s to uint", value)
		}
		converted[key] = uint(convertedValue)
	}
	return converted, nil
}
