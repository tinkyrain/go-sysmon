package disk

import (
	"io/fs"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tinkyrain/go-sysmon/internal/procfs"
)

func mountsFiles() map[string]string {
	return map[string]string{
		mountFile: `tmpfs /run tmpfs rw,nosuid,nodev,noexec,relatime,size=1300344k,mode=755,inode64 0 0
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
			Available: 50,
		},
		{
			Path:      "/",
			Total:     120,
			Available: 50,
		},
	}

	reader := mountsReader{
		fs: procfs.New(tempDirWithFiles(t, mountsFiles())),
		statfs: func(path string) (syscall.Statfs_t, error) {
			result := syscall.Statfs_t{}
			result.Bsize = 10
			result.Blocks = 12
			result.Bavail = 5
			return result, nil
		},
	}

	result, err := reader.read()

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestReadMountsError(t *testing.T) {
	cases := []struct {
		name               string
		fileContent        string
		statfsFunc         statfsFunc
		expectedError      error
		expectedReadResult []Mount
	}{
		{
			name:               "Empty file",
			fileContent:        "",
			statfsFunc:         func(string) (syscall.Statfs_t, error) { return syscall.Statfs_t{}, nil },
			expectedError:      ErrEmptyMountsData,
			expectedReadResult: nil,
		},
		{
			name:               "Not found suitable mounts",
			fileContent:        "/dev/nvme0n1p1 /boot/efi vfat123 rw,relatime,fmask=0022,dmask=0022",
			statfsFunc:         func(string) (syscall.Statfs_t, error) { return syscall.Statfs_t{}, nil },
			expectedError:      ErrNoSuitableMountsData,
			expectedReadResult: nil,
		},
		{
			name:               "Statfs error",
			fileContent:        "/dev/nvme0n1p5 / ext4 rw,relatime 0 0",
			statfsFunc:         func(string) (syscall.Statfs_t, error) { return syscall.Statfs_t{}, fs.ErrNotExist },
			expectedError:      fs.ErrNotExist,
			expectedReadResult: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			files := mountsFiles()
			files[mountFile] = tc.fileContent

			reader := mountsReader{
				fs:     procfs.New(tempDirWithFiles(t, files)),
				statfs: tc.statfsFunc,
			}

			result, err := reader.read()

			require.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, tc.expectedReadResult, result)
		})
	}
}
