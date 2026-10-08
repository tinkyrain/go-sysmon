package disk

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tinkyrain/go-sysmon/internal/procfs"
	"github.com/tinkyrain/go-sysmon/internal/statfs"
)

func diskFiles() map[string]string {
	return map[string]string{
		mountFile: `
tmpfs /run tmpfs rw,nosuid,nodev,noexec,relatime,size=1300344k,mode=755,inode64 0 0
/dev/nvme0n1p1 /boot/efi vfat rw,relatime,fmask=0022,dmask=0022,codepage=437,iocharset=iso8859-1,shortname=mixed,errors=remount-ro 0 0
/dev/nvme0n1p5 / ext4 rw,relatime 0 0
`,
	}
}

func TestCollectSuccess(t *testing.T) {
	expected := Stats{
		Mounts: []Mount{
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
		},
	}

	collector := New(
		procfs.New(tempDirWithFiles(t, diskFiles())),
		func(string) (statfs.Stats, error) {
			result := statfs.Stats{}
			result.BlockSize = 10
			result.Blocks = 12
			result.Available = 5
			return result, nil
		},
	)

	result, err := collector.Collect()

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestCollectError(t *testing.T) {
	files := diskFiles()
	delete(files, mountFile)

	collector := New(
		procfs.New(tempDirWithFiles(t, files)),
		func(string) (statfs.Stats, error) {
			result := statfs.Stats{}
			result.BlockSize = 10
			result.Blocks = 12
			result.Available = 13
			return result, nil
		},
	)

	result, err := collector.Collect()

	require.Error(t, err)
	assert.Equal(t, Stats{}, result)
}
