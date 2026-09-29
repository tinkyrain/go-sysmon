package system

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/tinkyrain/go-sysmon/internal/procfs"
)

const loadAvgFile = "loadavg"

type loadAvgReader struct {
	fs procfs.FileScanner
}

type LoadAvg struct {
	OneMin       float64
	FiveMin      float64
	FifteenMin   float64
	RunningProcs uint64
	TotalProcs   uint64
	LastPID      uint64
}

var ErrInsufficientLoadAvg = errors.New("insufficient loadavg data")
var ErrIncorrectMetricLoadAvg = errors.New("incorrect loadavg data")

func (r loadAvgReader) read() (LoadAvg, error) {
	data, err := r.fs.ScanRows(loadAvgFile)
	if err != nil {
		return LoadAvg{}, err
	}

	if len(data) == 0 {
		return LoadAvg{}, ErrInsufficientLoadAvg
	}

	// data[0] - because /proc/loadavg has one row
	return parseLoadAvgLine(data[0])
}

func parseLoadAvgLine(line string) (LoadAvg, error) {
	fields := strings.Fields(line)
	if len(fields) < 5 {
		return LoadAvg{}, ErrInsufficientLoadAvg
	}

	var err error
	loadAvg := LoadAvg{}

	if loadAvg.OneMin, err = strconv.ParseFloat(fields[0], 64); err != nil {
		return LoadAvg{}, fmt.Errorf("parsing loadavg value %q: %w", fields[0], err)
	}
	if loadAvg.FiveMin, err = strconv.ParseFloat(fields[1], 64); err != nil {
		return LoadAvg{}, fmt.Errorf("parsing loadavg value %q: %w", fields[1], err)
	}
	if loadAvg.FifteenMin, err = strconv.ParseFloat(fields[2], 64); err != nil {
		return LoadAvg{}, fmt.Errorf("parsing loadavg value %q: %w", fields[2], err)
	}
	if loadAvg.LastPID, err = strconv.ParseUint(fields[4], 10, 64); err != nil {
		return LoadAvg{}, fmt.Errorf("parsing loadavg value %q: %w", fields[4], err)
	}

	running, total, ok := strings.Cut(fields[3], "/")
	if !ok {
		return LoadAvg{}, fmt.Errorf("parsing loadavg value %q: %w", fields[3], ErrIncorrectMetricLoadAvg)
	}

	if loadAvg.RunningProcs, err = strconv.ParseUint(running, 10, 64); err != nil {
		return LoadAvg{}, fmt.Errorf("parsing loadavg value %q: %w", running, err)
	}
	if loadAvg.TotalProcs, err = strconv.ParseUint(total, 10, 64); err != nil {
		return LoadAvg{}, fmt.Errorf("parsing loadavg value %q: %w", total, err)
	}

	return loadAvg, nil
}
