package metrics

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/tinkyrain/go-sysmon/internal/procfs"
)

type CPUReader struct {
	fs             procfs.FileScanner
	prevCPUSamples map[string]CPUSample
}

type CPUTicks struct {
	ID        string
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

type CPUSample struct {
	Total   uint64
	Idle    uint64
	NonIdle uint64
}

type CPUUsage struct {
	ID    string
	Usage float64
}

const cpuFilename = "stat"
const cpuMetricPrefix = "cpu"

var ErrFewMetricsCountForParsing = errors.New("parsing CPU ticks need 10 metrics")
var ErrNoCPULines = errors.New("no cpu lines in file")

func (r *CPUReader) Read() ([]CPUUsage, error) {
	result := []CPUUsage{}

	data, err := r.fs.ScanRows(cpuFilename)
	if err != nil {
		return []CPUUsage{}, err
	}

	var lines [][]string

	for _, line := range data {
		fields := strings.Fields(line)
		// 11 because cpu has 11 ticks value
		if len(fields) < 11 {
			continue
		}
		if !strings.HasPrefix(fields[0], cpuMetricPrefix) {
			continue
		}
		lines = append(lines, fields)
	}

	if len(lines) == 0 {
		return []CPUUsage{}, ErrNoCPULines
	}

	samples := map[string]CPUSample{}

	for _, line := range lines {
		t, err := parseCPULine(line)
		if err != nil {
			return []CPUUsage{}, err
		}

		idle := t.Idle + t.Iowait
		nonIdle := t.User + t.Nice + t.System + t.Irq +
			t.Softirq + t.Steal
		total := idle + nonIdle

		sample := CPUSample{
			Total:   total,
			Idle:    idle,
			NonIdle: nonIdle,
		}

		samples[t.ID] = sample

		prevSample, ok := r.prevCPUSamples[t.ID]
		var usage float64

		if ok {
			usage = calculateCPUUsage(prevSample, sample)
		}

		result = append(result, CPUUsage{
			ID:    t.ID,
			Usage: usage,
		})
	}

	r.prevCPUSamples = samples

	return result, nil
}

func calculateCPUUsage(prev, cur CPUSample) float64 {
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

	return math.Round(percentUsage*100) / 100
}

func parseCPULine(line []string) (CPUTicks, error) {
	if len(line) < 11 {
		return CPUTicks{}, ErrFewMetricsCountForParsing
	}
	coreName := line[0]
	line = line[1:]
	converted := make([]uint64, 0, len(line))
	for _, value := range line {
		convertedValue, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return CPUTicks{}, fmt.Errorf("parsing cpu line %q: %w", value, err)
		}
		converted = append(converted, convertedValue)
	}

	return CPUTicks{
		ID:        coreName,
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
