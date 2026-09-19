package metrics

import (
	"errors"
	"go-sysmon/internal/procfs"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	validMeminfo = "MemTotal:       13003440 kB\n" +
		"MemAvailable:    5121912 kB\n" +
		"SwapTotal:       4194300 kB\n" +
		"SwapFree:        3383196 kB\n"

	validMounts = "/dev/nvme0n1p5 / ext4 rw,relatime 0 0\n" +
		"tmpfs /run tmpfs rw 0 0\n"

	validStat = "cpu  298830 3399 75896 7741115 5083 0 1586 0 0 0\n" +
		"cpu0 17496 72 6050 650357 437 0 376 0 0 0\n"
)

var errStatfsFail = errors.New("statfs failed")

func okStatfs(string) (syscall.Statfs_t, error) {
	return syscall.Statfs_t{Bsize: 1024, Blocks: 100, Bavail: 40}, nil
}

func errStatfs(string) (syscall.Statfs_t, error) {
	return syscall.Statfs_t{}, errStatfsFail
}

func procDir(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, content := range files {
		require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600))
	}
	return dir
}

func newTestCollector(dir string, statfsFunc func(string) (syscall.Statfs_t, error)) *Collector {
	fs := procfs.New(dir)
	return &Collector{
		memoryReader: MemoryReader{fs: fs},
		diskReader:   DiskReader{fs: fs, statfsFunc: statfsFunc},
		cpuReader:    CPUReader{fs: fs},
	}
}

func TestNew(t *testing.T) {
	c := New("/proc")

	require.NotNil(t, c)
	// все ридеры смотрят в один и тот же procRoot
	assert.Equal(t, procfs.New("/proc"), c.memoryReader.fs)
	assert.Equal(t, procfs.New("/proc"), c.diskReader.fs)
	assert.Equal(t, procfs.New("/proc"), c.cpuReader.fs)
	assert.NotNil(t, c.diskReader.statfsFunc)
}

func TestCollectSuccess(t *testing.T) {
	dir := procDir(t, map[string]string{
		"meminfo": validMeminfo,
		"mounts":  validMounts,
		"stat":    validStat,
	})
	c := newTestCollector(dir, okStatfs)

	snap, err := c.Collect()

	require.NoError(t, err)
	assert.Equal(t, Memory{
		Total:         13003440,
		Available:     5121912,
		SwapTotal:     4194300,
		SwapAvailable: 3383196,
	}, snap.Memory)
	assert.Equal(t, []Disk{
		{Mount: "/", Total: 100 * 1024, Available: 40 * 1024},
	}, snap.Disks)
	assert.Zero(t, snap.CPUUsage)
	assert.WithinDuration(t, time.Now(), snap.Time, time.Minute)
}

func TestCollectMemoryError(t *testing.T) {
	dir := procDir(t, map[string]string{
		"mounts": validMounts,
		"stat":   validStat,
	})
	c := newTestCollector(dir, okStatfs)

	_, err := c.Collect()

	require.Error(t, err)
	assert.ErrorIs(t, err, fs.ErrNotExist)
}

func TestCollectDiskError(t *testing.T) {
	dir := procDir(t, map[string]string{
		"meminfo": validMeminfo,
		"mounts":  validMounts,
		"stat":    validStat,
	})
	c := newTestCollector(dir, errStatfs)

	_, err := c.Collect()

	require.Error(t, err)
	assert.ErrorIs(t, err, errStatfsFail)
}

func TestCollectCPUError(t *testing.T) {
	dir := procDir(t, map[string]string{
		"meminfo": validMeminfo,
		"mounts":  validMounts,
	})
	c := newTestCollector(dir, okStatfs)

	_, err := c.Collect()

	require.Error(t, err)
	assert.ErrorIs(t, err, fs.ErrNotExist)
}
