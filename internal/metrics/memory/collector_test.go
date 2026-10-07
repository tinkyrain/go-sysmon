package memory

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tinkyrain/go-sysmon/internal/procfs"
)

func memoryFiles() map[string]string {
	return map[string]string{
		meminfoFile: `
MemTotal:       13003440 kB
MemFree:          199488 kB
MemAvailable:    5121912 kB
SwapTotal:       4194300 kB
SwapFree:        3383196 kB
Dirty:              1788 kB
		`,
	}
}

func TestCollectSuccess(t *testing.T) {
	expected := Stats{
		Memory: MemInfo{
			Total:         13315522560,
			Available:     5244837888,
			SwapTotal:     4294963200,
			SwapAvailable: 3464392704,
		},
	}
	collector := New(procfs.New(tempDirWithFiles(t, meminfoFiles())))

	result, err := collector.Collect()

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestCollectError(t *testing.T) {
	files := memoryFiles()
	delete(files, meminfoFile)

	collector := New(procfs.New(tempDirWithFiles(t, files)))

	result, err := collector.Collect()

	require.ErrorIs(t, err, os.ErrNotExist)
	assert.Equal(t, Stats{}, result)
}
