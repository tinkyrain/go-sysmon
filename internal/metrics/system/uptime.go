package system

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/tinkyrain/go-sysmon/internal/procfs"
)

const uptimeFile = "uptime"

type uptimeReader struct {
	fs procfs.FileScanner
}

type Uptime struct {
	Total float64
	Idle  float64
}

var (
	ErrIncorrectUptimeData = errors.New("uptime data is incorrect")
	ErrEmptyUptimeData     = errors.New("uptime data is empty")
)

func (r uptimeReader) read() (Uptime, error) {
	data, err := r.fs.ScanRow(uptimeFile)
	if err != nil {
		return Uptime{}, fmt.Errorf("scan %q: %w", uptimeFile, err)
	}

	if len(data) == 0 {
		return Uptime{}, fmt.Errorf("read %q: %w", uptimeFile, ErrEmptyUptimeData)
	}

	return parseUptimeLine(data)
}

func parseUptimeLine(line string) (Uptime, error) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return Uptime{}, fmt.Errorf("parse line: %w", ErrIncorrectUptimeData)
	}

	var err error
	uptime := Uptime{}

	if uptime.Total, err = strconv.ParseFloat(fields[0], 64); err != nil {
		return Uptime{}, fmt.Errorf("parsing uptime value %q: %w", fields[0], err)
	}
	if uptime.Idle, err = strconv.ParseFloat(fields[1], 64); err != nil {
		return Uptime{}, fmt.Errorf("parsing uptime value %q: %w", fields[1], err)
	}

	return uptime, nil
}
