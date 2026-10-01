package memory

import (
	"strconv"
	"testing"

	"github.com/tinkyrain/go-sysmon/internal/procfs"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func meminfoFiles() map[string]string {
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

func TestReadMemorySuccess(t *testing.T) {
	expected := Memory{
		Total:         13315522560,
		Available:     5244837888,
		SwapTotal:     4294963200,
		SwapAvailable: 3464392704,
	}

	reader := memoryReader{procfs.New(tempDirWithFiles(t, meminfoFiles(), 0o755, 0o600))}

	result, err := reader.read()

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestReadMemoryBlankFile(t *testing.T) {
	files := meminfoFiles()
	files[meminfoFile] = ""

	reader := memoryReader{fs: procfs.New(tempDirWithFiles(t, files, 0o755, 0o600))}

	result, err := reader.read()

	require.NoError(t, err)
	assert.Equal(t, Memory{}, result)
}

func TestReadMemoryNeedleMetricsNotFound(t *testing.T) {
	files := meminfoFiles()
	files[meminfoFile] = "Buffers: 338020 kB\nCached: 1234 kB"

	reader := memoryReader{fs: procfs.New(tempDirWithFiles(t, files, 0o755, 0o600))}

	result, err := reader.read()

	require.NoError(t, err)
	assert.Equal(t, Memory{}, result)
}

func TestReadMemoryFileNotFound(t *testing.T) {
	files := meminfoFiles()
	delete(files, meminfoFile)

	reader := memoryReader{fs: procfs.New(tempDirWithFiles(t, files, 0o755, 0o600))}

	result, err := reader.read()

	require.Error(t, err)
	assert.Equal(t, Memory{}, result)
}

func TestReadMemoryErrorConvertMetrics(t *testing.T) {
	files := meminfoFiles()
	files[meminfoFile] = "Buffers: 338020 kB\nMemAvailable: is_not_converted_string kB"
	reader := memoryReader{fs: procfs.New(tempDirWithFiles(t, files, 0o755, 0o600))}

	result, err := reader.read()

	require.ErrorIs(t, err, strconv.ErrSyntax)
	assert.Equal(t, Memory{}, result)
}

func TestReadMemoryIncorrectMetricLine(t *testing.T) {
	files := meminfoFiles()
	files[meminfoFile] = "Buffers 338020 kB\nMemAvailable 1321223 kB"

	reader := memoryReader{fs: procfs.New(tempDirWithFiles(t, files, 0o755, 0o600))}

	result, err := reader.read()

	require.NoError(t, err)
	assert.Equal(t, Memory{}, result)
}

func TestReadMemoryBlankMetricRow(t *testing.T) {
	files := meminfoFiles()
	files[meminfoFile] = "Buffers: 338020 kB\nMemAvailable:"

	reader := memoryReader{fs: procfs.New(tempDirWithFiles(t, files, 0o755, 0o600))}

	result, err := reader.read()

	require.NoError(t, err)
	assert.Equal(t, Memory{}, result)
}
