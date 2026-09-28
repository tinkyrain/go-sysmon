package system

import (
	"errors"
	"strconv"
	"strings"

	"github.com/tinkyrain/go-sysmon/internal/procfs"
)

const loadAvgFilename = "loadavg"

type LoadAvgReader struct {
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

var (
	ErrEmptyLoadAvg            = errors.New("loadavg data is empty")
	ErrInsufficientLoadAvgData = errors.New("insufficient loadavg data")
)

func (r *LoadAvgReader) Read() (LoadAvg, error) {
	data, err := r.fs.ScanRows(loadAvgFilename)
	if err != nil {
		return LoadAvg{}, err
	}

	if len(data) == 0 {
		return LoadAvg{}, ErrEmptyLoadAvg
	}

	// data[0] - because /proc/loadavg has one row
	return parseLoadAvgLine(data[0])
}

func parseLoadAvgLine(line string) (LoadAvg, error) {
	fields := strings.Fields(line)
	if len(fields) < 5 {
		return LoadAvg{}, ErrInsufficientLoadAvgData
	}

	var err error
	loadAvg := LoadAvg{}

	if loadAvg.OneMin, err = strconv.ParseFloat(fields[0], 64); err != nil {
		return LoadAvg{}, err
	}
	if loadAvg.FiveMin, err = strconv.ParseFloat(fields[1], 64); err != nil {
		return LoadAvg{}, err
	}
	if loadAvg.FifteenMin, err = strconv.ParseFloat(fields[2], 64); err != nil {
		return LoadAvg{}, err
	}
	if loadAvg.LastPID, err = strconv.ParseUint(fields[4], 10, 64); err != nil {
		return LoadAvg{}, err
	}

	running, total, ok := strings.Cut(fields[3], "/")
	if !ok {
		return LoadAvg{}, ErrInsufficientLoadAvgData
	}

	if loadAvg.RunningProcs, err = strconv.ParseUint(running, 10, 64); err != nil {
		return LoadAvg{}, err
	}
	if loadAvg.TotalProcs, err = strconv.ParseUint(total, 10, 64); err != nil {
		return LoadAvg{}, err
	}

	return loadAvg, nil
}
