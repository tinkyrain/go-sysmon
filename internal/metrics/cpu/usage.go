package cpu

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/tinkyrain/go-sysmon/internal/procfs"
)

const statFile = "stat"
const cpuPrefix = "cpu"

type usageReader struct {
	fs               procfs.FileScanner
	prevUsageSamples map[string]usageSample
}

type ticks struct {
	id        string
	user      uint64
	nice      uint64
	system    uint64
	idle      uint64
	iowait    uint64
	irq       uint64
	softirq   uint64
	steal     uint64
	guest     uint64
	guestNice uint64
}

type usageSample struct {
	total   uint64
	idle    uint64
	nonIdle uint64
}

type CoreUsage struct {
	ID    string
	Usage float64
}

type Usage struct {
	Total float64
	Cores []CoreUsage
}

var (
	ErrIncorrectStatData = errors.New("usage is incorrect")
	ErrEmptyStatData     = errors.New("usage is empty")
)

func (r *usageReader) read() (Usage, error) {
	data, err := r.fs.ScanRows(statFile)
	if err != nil {
		return Usage{}, fmt.Errorf("scan %q: %w", statFile, err)
	}

	var lines [][]string

	for _, line := range data {
		fields := strings.Fields(line)
		// 11 because cpu has 11 ticks value
		if len(fields) < 11 {
			continue
		}
		if !strings.HasPrefix(fields[0], cpuPrefix) {
			continue
		}
		lines = append(lines, fields)
	}

	if len(lines) == 0 {
		return Usage{}, fmt.Errorf("read %q: %w", statFile, ErrEmptyStatData)
	}

	usage := Usage{}
	samples := map[string]usageSample{}

	for _, line := range lines {
		t, err := parseTicksLine(line)
		if err != nil {
			return Usage{}, err
		}

		idle := t.idle + t.iowait
		nonIdle := t.user + t.nice + t.system + t.irq +
			t.softirq + t.steal
		total := idle + nonIdle

		sample := usageSample{
			total:   total,
			idle:    idle,
			nonIdle: nonIdle,
		}

		samples[t.id] = sample

		prevSample, ok := r.prevUsageSamples[t.id]
		var calculatedUsage float64

		if ok {
			calculatedUsage = calculateUsage(prevSample, sample)
		}

		if t.id == cpuPrefix {
			usage.Total = calculatedUsage
		} else {
			usage.Cores = append(usage.Cores, CoreUsage{
				ID:    t.id,
				Usage: calculatedUsage,
			})
		}
	}

	r.prevUsageSamples = samples

	return usage, nil
}

func calculateUsage(prev, cur usageSample) float64 {
	var percentUsage float64

	if prev.total == 0 || cur.idle < prev.idle || cur.nonIdle < prev.nonIdle {
		return 0
	}

	idleDelta := cur.idle - prev.idle
	nonIdleDelta := cur.nonIdle - prev.nonIdle
	totalDelta := idleDelta + nonIdleDelta

	if nonIdleDelta > 0 && totalDelta > 0 {
		percentUsage = float64(nonIdleDelta) / float64(totalDelta) * 100
	}

	return percentUsage
}

func parseTicksLine(line []string) (ticks, error) {
	if len(line) < 11 {
		return ticks{}, fmt.Errorf("parse line %q: %w", strings.Join(line, " "), ErrIncorrectStatData)
	}
	id := line[0]
	line = line[1:]
	converted := make([]uint64, 0, len(line))
	for _, value := range line {
		convertedValue, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return ticks{}, fmt.Errorf("parse value in line %q: %w", value, err)
		}
		converted = append(converted, convertedValue)
	}

	return ticks{
		id:        id,
		user:      converted[0],
		nice:      converted[1],
		system:    converted[2],
		idle:      converted[3],
		iowait:    converted[4],
		irq:       converted[5],
		softirq:   converted[6],
		steal:     converted[7],
		guest:     converted[8],
		guestNice: converted[9],
	}, nil
}
