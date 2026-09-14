package metrics

import (
	"fmt"
	"go-sysmon/internal/procfs"
	"strings"
	"syscall"
)

type Disk struct {
	Mount     string
	Total     uint64
	Available uint64
}

const mountFilename = "mounts"

var diskPrefix = [3]string{"/dev/sd", "/dev/nvme", "/dev/vd"}

func (c *Collector) readDisks() ([]Disk, error) {
	var disks []Disk
	mountPaths, err := getMounts(c.fs, mountFilename)
	if err != nil {
		return disks, err
	}

	for _, mountPath := range mountPaths {
		var stat syscall.Statfs_t
		err := syscall.Statfs(mountPath, &stat)
		if err != nil {
			return disks, fmt.Errorf("error getting disk stats: %s", mountPath)
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

func getMounts(scanner procfs.FileScanner, filename string) ([]string, error) {
	mountData, err := scanner.ScanRows(filename)
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
