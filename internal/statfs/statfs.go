package statfs

import (
	"fmt"
	"syscall"
)

type Stats struct {
	BlockSize uint64
	Blocks    uint64
	Available uint64
}

type Func func(dir string) (Stats, error)

func GetDirStatfs(dir string) (Stats, error) {
	var s syscall.Statfs_t
	err := syscall.Statfs(dir, &s)
	if err != nil {
		return Stats{}, fmt.Errorf("statfs %q: %w", dir, err)
	}
	return Stats{
		BlockSize: uint64(s.Bsize),
		Blocks:    s.Blocks,
		Available: s.Bavail,
	}, nil
}
