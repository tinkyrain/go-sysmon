package statfs

import (
	"fmt"
	"syscall"
)

func GetDirStatfs(dir string) (syscall.Statfs_t, error) {
	var result syscall.Statfs_t
	err := syscall.Statfs(dir, &result)
	if err != nil {
		return result, fmt.Errorf("statfs %s: %w", dir, err)
	}
	return result, nil
}
