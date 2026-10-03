package disk

import (
	"fmt"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tinkyrain/go-sysmon/internal/procfs"
)

func mountsFiles() map[string]string {
	return map[string]string{
		mountFile: `
tmpfs /run tmpfs rw,nosuid,nodev,noexec,relatime,size=1300344k,mode=755,inode64 0 0
/dev/nvme0n1p1 /boot/efi vfat rw,relatime,fmask=0022,dmask=0022,codepage=437,iocharset=iso8859-1,shortname=mixed,errors=remount-ro 0 0
/dev/nvme0n1p5 / ext4 rw,relatime 0 0
`,
	}
}

func TestReadMountsSuccess(t *testing.T) {
	expected := []Mount{
		{
			Path:      "/boot/efi",
			Total:     120,
			Available: 130,
		},
		{
			Path:      "/",
			Total:     120,
			Available: 130,
		},
	}

	reader := mountsReader{
		fs: procfs.New(tempDirWithFiles(t, mountsFiles(), 0o755, 0o600)),
		statfsFunc: func(path string) (syscall.Statfs_t, error) {
			result := syscall.Statfs_t{}
			result.Bsize = 10
			result.Blocks = 12
			result.Bavail = 13
			return result, nil
		},
	}

	result, err := reader.read()

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestReadMountsBlankFile(t *testing.T) {
	files := mountsFiles()
	files[mountFile] = ""

	reader := mountsReader{
		fs: procfs.New(tempDirWithFiles(t, files, 0o755, 0o600)),
		statfsFunc: func(path string) (syscall.Statfs_t, error) {
			result := syscall.Statfs_t{}
			result.Bsize = 10
			result.Blocks = 10
			result.Bavail = 10
			return result, nil
		},
	}

	result, err := reader.read()

	require.ErrorIs(t, err, ErrEmptyMountsFile)
	assert.Equal(t, []Mount{}, result)
}

func TestReadDiskNeedleMetricsNotFound(t *testing.T) {
	files := mountsFiles()
	files[mountFile] = "sysfs /sys sysfs rw,nosuid,nodev,noexec,relatime 0 0\nproc /proc proc rw,nosuid,nodev,noexec,relatime 0 0"

	reader := mountsReader{
		fs: procfs.New(tempDirWithFiles(t, files, 0o755, 0o600)),
		statfsFunc: func(path string) (syscall.Statfs_t, error) {
			result := syscall.Statfs_t{}
			result.Bsize = 10
			result.Blocks = 10
			result.Bavail = 10
			return result, nil
		},
	}

	result, err := reader.read()

	require.ErrorIs(t, err, ErrInsufficientMounts)
	assert.Equal(t, []Mount{}, result)
}

func TestReadDiskFileNotFound(t *testing.T) {
	files := mountsFiles()
	delete(files, mountFile)

	reader := mountsReader{
		fs:         procfs.New(tempDirWithFiles(t, files, 0o755, 0o600)),
		statfsFunc: func(path string) (syscall.Statfs_t, error) { return syscall.Statfs_t{}, nil },
	}

	result, err := reader.read()

	require.Error(t, err)
	assert.Equal(t, []Mount{}, result)
}

func TestReadDiskStatfsError(t *testing.T) {
	reader := mountsReader{
		fs:         procfs.New(tempDirWithFiles(t, mountsFiles(), 0o755, 0o600)),
		statfsFunc: func(path string) (syscall.Statfs_t, error) { return syscall.Statfs_t{}, fmt.Errorf("statfs failed") },
	}

	result, err := reader.read()

	require.Error(t, err)
	assert.Equal(t, []Mount{}, result)
}

func TestReadDiskIncorrectMetricLine(t *testing.T) {
	files := mountsFiles()
	files[mountFile] = "/dev/nvme0n1p1"

	reader := mountsReader{
		fs:         procfs.New(tempDirWithFiles(t, files, 0o755, 0o600)),
		statfsFunc: func(path string) (syscall.Statfs_t, error) { return syscall.Statfs_t{}, nil },
	}

	result, err := reader.read()

	require.ErrorIs(t, err, ErrInsufficientMounts)
	assert.Equal(t, []Mount{}, result)
}
