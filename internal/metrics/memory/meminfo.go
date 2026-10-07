package memory

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/tinkyrain/go-sysmon/internal/procfs"
)

const meminfoFile string = "meminfo"

type meminfoReader struct {
	fs procfs.FileScanner
}

type MemInfo struct {
	Total         uint64
	Available     uint64
	SwapTotal     uint64
	SwapAvailable uint64
}

var (
	ErrEmptyMemInfoData     = errors.New("meminfo data is empty")
	ErrIncorrectMemInfoData = errors.New("meminfo data is incorrect")
)

func (r meminfoReader) read() (MemInfo, error) {
	data, err := r.fs.ScanRows(meminfoFile)
	if err != nil {
		return MemInfo{}, fmt.Errorf("scan %q: %w", meminfoFile, err)
	}

	if len(data) == 0 {
		return MemInfo{}, fmt.Errorf("read %q: %w", meminfoFile, ErrEmptyMemInfoData)
	}

	dataHasMetric := map[string]struct{}{}

	meminfo := MemInfo{}
	info := map[string]*uint64{
		"MemTotal":     &meminfo.Total,
		"MemAvailable": &meminfo.Available,
		"SwapTotal":    &meminfo.SwapTotal,
		"SwapFree":     &meminfo.SwapAvailable,
	}

	for _, line := range data {
		name, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		if _, ok := info[name]; ok {
			fields := strings.Fields(value)
			if len(fields) == 0 {
				continue
			}
			convertedValue, err := strconv.ParseUint(fields[0], 10, 64)
			if err != nil {
				return MemInfo{}, fmt.Errorf("parsing metric %q value: %w", name, err)
			}
			dataHasMetric[name] = struct{}{}
			*info[name] = convertedValue * 1024 // Kb in bytes
		}
	}

	for metric := range info {
		if _, ok := dataHasMetric[metric]; !ok {
			return MemInfo{}, fmt.Errorf(
				"read %q metric %q: %w",
				meminfoFile,
				metric,
				ErrIncorrectMemInfoData,
			)
		}
	}

	return meminfo, nil
}
