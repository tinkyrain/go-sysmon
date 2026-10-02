package memory

import (
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

func (r *meminfoReader) read() (MemInfo, error) {
	data, err := r.fs.ScanRows(meminfoFile)
	if err != nil {
		return MemInfo{}, err
	}

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
			convetredValue, err := strconv.ParseUint(fields[0], 10, 64)
			if err != nil {
				return MemInfo{}, fmt.Errorf("error parsing metric %q value: %w", name, err)
			}
			*info[name] = convetredValue * 1024 // Kb in bytes
		}
	}

	return meminfo, nil
}
