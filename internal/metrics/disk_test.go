package metrics

import (
	"fmt"
	"syscall"
	"testing"

	"github.com/tinkyrain/go-sysmon/internal/procfs"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadDiskSuccess(t *testing.T) {
	expected := []Disk{
		{
			Mount:     "/",
			Total:     120,
			Available: 130,
		},
		{
			Mount:     "/boot/efi",
			Total:     120,
			Available: 130,
		},
	}
	fs := procfs.New("testdata/")
	reader := DiskReader{
		fs: fs,
		statfsFunc: func(path string) (syscall.Statfs_t, error) {
			result := syscall.Statfs_t{}
			result.Bsize = 10
			result.Blocks = 12
			result.Bavail = 13
			return result, nil
		},
	}

	readResult, err := reader.Read()

	require.NoError(t, err)
	assert.Equal(t, expected, readResult)
}

func TestReadDiskBlankFile(t *testing.T) {
	expected := []Disk{}
	tempDir := tempDirWithFile(t, "mounts", "", 0o600)
	fs := procfs.New(tempDir)
	reader := DiskReader{
		fs: fs,
		statfsFunc: func(path string) (syscall.Statfs_t, error) {
			result := syscall.Statfs_t{}
			result.Bsize = 10
			result.Blocks = 10
			result.Bavail = 10
			return result, nil
		},
	}

	readResult, err := reader.Read()

	require.NoError(t, err)
	assert.Equal(t, expected, readResult)
}

func TestReadDiskNeedleMetricsNotFound(t *testing.T) {
	expected := []Disk{}
	dir := tempDirWithFile(
		t,
		"mounts",
		"sysfs /sys sysfs rw,nosuid,nodev,noexec,relatime 0 0\nproc /proc proc rw,nosuid,nodev,noexec,relatime 0 0",
		0o600)
	fs := procfs.New(dir)
	reader := DiskReader{
		fs: fs,
		statfsFunc: func(path string) (syscall.Statfs_t, error) {
			result := syscall.Statfs_t{}
			result.Bsize = 10
			result.Blocks = 10
			result.Bavail = 10
			return result, nil
		},
	}

	readResult, err := reader.Read()

	require.NoError(t, err)
	assert.Equal(t, expected, readResult)
}

func TestReadDiskFileNotFound(t *testing.T) {
	tempDir := tempDirWithFile(t, "mounts_not_found", "", 0o600)
	fs := procfs.New(tempDir)
	reader := DiskReader{
		fs:         fs,
		statfsFunc: func(path string) (syscall.Statfs_t, error) { return syscall.Statfs_t{}, nil },
	}

	_, err := reader.Read()

	assert.Error(t, err)
}

func TestReadDiskStatfsError(t *testing.T) {
	fs := procfs.New("testdata/")
	reader := DiskReader{
		fs:         fs,
		statfsFunc: func(path string) (syscall.Statfs_t, error) { return syscall.Statfs_t{}, fmt.Errorf("statfs failed") },
	}

	_, err := reader.Read()

	assert.Error(t, err)
}

func TestReadDiskIncorrectMetricLine(t *testing.T) {
	expected := []Disk{}
	tempDir := tempDirWithFile(t, "mounts", "/dev/nvme0n1p1", 0o600)
	fs := procfs.New(tempDir)
	reader := DiskReader{
		fs:         fs,
		statfsFunc: func(path string) (syscall.Statfs_t, error) { return syscall.Statfs_t{}, nil },
	}

	readResult, err := reader.Read()

	require.NoError(t, err)
	assert.Equal(t, expected, readResult)
}
