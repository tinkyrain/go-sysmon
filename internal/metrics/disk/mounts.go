package disk

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tinkyrain/go-sysmon/internal/procfs"
	"github.com/tinkyrain/go-sysmon/internal/statfs"
)

type mountsReader struct {
	fs     procfs.FileScanner
	statfs statfs.Func
}

type Mount struct {
	Path      string
	Total     uint64
	Available uint64
}

const mountFile = "mounts"

var fileSystemTypes = map[string]struct{}{
	"ext4":    {},
	"xfs":     {},
	"btrfs":   {},
	"vfat":    {},
	"ntfs3":   {},
	"zfs":     {},
	"fuseblk": {},
	"ext2":    {},
	"ext3":    {},
	"exfat":   {},
	"f2fs":    {},
}

var (
	ErrEmptyMountsData      = errors.New("mounts data is empty")
	ErrNoSuitableMountsData = errors.New("suitable mounts not found")
)

func (r mountsReader) read() ([]Mount, error) {
	data, err := r.fs.ScanRows(mountFile)
	if err != nil {
		return nil, fmt.Errorf("scan %q: %w", mountFile, err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("read %q: %w", mountFile, ErrEmptyMountsData)
	}

	var paths []string

	for _, line := range data {
		fields := strings.Fields(line)

		if len(fields) < 3 {
			continue
		}

		if _, ok := fileSystemTypes[fields[2]]; ok {
			paths = append(paths, fields[1])
		}
	}

	if len(paths) == 0 {
		return nil, fmt.Errorf("read %q: %w", mountFile, ErrNoSuitableMountsData)
	}

	mounts := []Mount{}

	for _, path := range paths {
		stat, err := r.statfs(path)
		if err != nil {
			return nil, fmt.Errorf("read %q: %w", path, err)
		}
		mounts = append(mounts, Mount{
			Path:      path,
			Total:     stat.Blocks * stat.BlockSize,
			Available: stat.Available * stat.BlockSize,
		})
	}

	return mounts, nil
}
