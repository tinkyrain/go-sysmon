package disk

import (
	"errors"
	"strings"
	"syscall"

	"github.com/tinkyrain/go-sysmon/internal/procfs"
)

type mountsReader struct {
	fs         procfs.FileScanner
	statfsFunc func(string) (syscall.Statfs_t, error)
}

type Mount struct {
	Path      string
	Total     uint64
	Available uint64
}

const mountFile = "mounts"

var pathPrefix = [3]string{"/dev/sd", "/dev/nvme", "/dev/vd"}

var (
	ErrEmptyMountsFile    = errors.New("empty mounts file")
	ErrInsufficientMounts = errors.New("insufficient mounts data")
)

func (r *mountsReader) read() ([]Mount, error) {
	data, err := r.fs.ScanRows(mountFile)
	if err != nil {
		return []Mount{}, err
	}

	if len(data) == 0 {
		return []Mount{}, ErrEmptyMountsFile
	}

	var paths []string

	for _, line := range data {
		fields := strings.Fields(line)

		if len(fields) < 2 {
			continue
		}

		for _, prefix := range pathPrefix {
			if strings.HasPrefix(fields[0], prefix) {
				paths = append(paths, fields[1])
				break
			}
		}
	}

	if len(paths) == 0 {
		return []Mount{}, ErrInsufficientMounts
	}

	mounts := []Mount{}

	for _, path := range paths {
		stat, err := r.statfsFunc(path)
		if err != nil {
			return []Mount{}, err
		}
		blockSize := uint64(stat.Bsize)
		mounts = append(mounts, Mount{
			Path:      path,
			Total:     stat.Blocks * blockSize,
			Available: stat.Bavail * blockSize,
		})
	}

	return mounts, nil
}
