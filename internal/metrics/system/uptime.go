package system

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/tinkyrain/go-sysmon/internal/procfs"
)

const uptimeFilename = "uptime"

type UptimeReader struct {
	fs procfs.FileScanner
}

type Uptime struct {
	Total float64
	Idle  float64
}

var ErrInsufficientUptime = errors.New("insufficient uptime data")

func (r *UptimeReader) Read() (Uptime, error) {
	data, err := r.fs.ScanRows(uptimeFilename)
	if err != nil {
		return Uptime{}, err
	}

	if len(data) == 0 {
		return Uptime{}, ErrInsufficientUptime
	}

	// data[0] - because /proc/uptime has one row
	return parseUptimeLine(data[0])
}

func parseUptimeLine(line string) (Uptime, error) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return Uptime{}, ErrInsufficientUptime
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
