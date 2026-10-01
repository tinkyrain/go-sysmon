package memory

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/tinkyrain/go-sysmon/internal/procfs"
)

const meminfoFile string = "meminfo"

type memoryReader struct {
	fs procfs.FileScanner
}

type Memory struct {
	Total         uint64
	Available     uint64
	SwapTotal     uint64
	SwapAvailable uint64
}

func (r *memoryReader) read() (Memory, error) {
	data, err := r.fs.ScanRows(meminfoFile)
	if err != nil {
		return Memory{}, err
	}

	memory := Memory{}
	info := map[string]*uint64{
		"MemTotal":     &memory.Total,
		"MemAvailable": &memory.Available,
		"SwapTotal":    &memory.SwapTotal,
		"SwapFree":     &memory.SwapAvailable,
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
				return Memory{}, fmt.Errorf("error parsing metric %q value: %w", name, err)
			}
			*info[name] = convetredValue * 1024 // Kb in bytes
		}
	}

	return memory, nil
}
