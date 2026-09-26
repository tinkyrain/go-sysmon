package metrics

import (
	"strings"
	"syscall"

	"github.com/tinkyrain/go-sysmon/internal/procfs"
)

type DiskReader struct {
	fs         procfs.FileScanner
	statfsFunc func(string) (syscall.Statfs_t, error)
}

type Disk struct {
	Mount     string
	Total     uint64
	Available uint64
}

const mountFilename = "mounts"

var diskPrefix = [3]string{"/dev/sd", "/dev/nvme", "/dev/vd"}

func (r *DiskReader) Read() ([]Disk, error) {
	disks := []Disk{}
	paths, err := r.getMounts(mountFilename)
	if err != nil {
		return disks, err
	}

	for _, path := range paths {
		stat, err := r.statfsFunc(path)
		if err != nil {
			return disks, err
		}
		blockSize := uint64(stat.Bsize)
		disks = append(disks, Disk{
			Mount:     path,
			Total:     stat.Blocks * blockSize,
			Available: stat.Bavail * blockSize,
		})
	}

	return disks, nil
}

func (r *DiskReader) getMounts(filename string) ([]string, error) {
	data, err := r.fs.ScanRows(filename)
	if err != nil {
		return nil, err
	}
	var paths []string

	for _, line := range data {
		fields := strings.Fields(line)

		if len(fields) < 2 {
			continue
		}

		for _, prefix := range diskPrefix {
			if strings.HasPrefix(fields[0], prefix) {
				paths = append(paths, fields[1])
				break
			}
		}
	}

	return paths, nil
}
