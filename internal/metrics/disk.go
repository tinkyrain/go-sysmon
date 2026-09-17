package metrics

import (
	"go-sysmon/internal/procfs"
	"strings"
	"syscall"
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
	var disks []Disk
	mountPaths, err := r.getMounts(mountFilename)
	if err != nil {
		return disks, err
	}

	for _, mountPath := range mountPaths {
		stat, err := r.statfsFunc(mountPath)
		if err != nil {
			return disks, err
		}
		blockSize := uint64(stat.Bsize)
		disks = append(disks, Disk{
			Mount:     mountPath,
			Total:     stat.Blocks * blockSize,
			Available: stat.Bavail * blockSize,
		})
	}

	return disks, nil
}

func (r *DiskReader) getMounts(filename string) ([]string, error) {
	mountData, err := r.fs.ScanRows(filename)
	if err != nil {
		return nil, err
	}
	var mountPaths []string

	for _, line := range mountData {
		fields := strings.Fields(line)

		if len(fields) < 2 {
			continue
		}

		for _, prefix := range diskPrefix {
			if strings.HasPrefix(fields[0], prefix) {
				mountPaths = append(mountPaths, fields[1])
				break
			}
		}
	}

	return mountPaths, nil
}
