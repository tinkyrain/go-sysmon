package metrics

import (
	"go-sysmon/internal/procfs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func procDir(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, content := range files {
		require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(content), 0o600))
	}
	return dir
}

func TestReadMemorySuccess(t *testing.T) {
	expected := Memory{
		Total:         13003440,
		Available:     5121912,
		SwapTotal:     4194300,
		SwapAvailable: 3383196,
	}
	fs := procfs.New("testdata/")
	reader := MemoryReader{fs}

	readResult, err := reader.Read()

	require.NoError(t, err)
	assert.Equal(t, expected, readResult)
}

func TestReadMemoryBlankFile(t *testing.T) {
	expected := Memory{
		Total:         0,
		Available:     0,
		SwapTotal:     0,
		SwapAvailable: 0,
	}
	tempDir := procDir(t, map[string]string{
		"meminfo": "",
	})
	fs := procfs.New(tempDir)
	reader := MemoryReader{fs}

	readResult, err := reader.Read()

	require.NoError(t, err)
	assert.Equal(t, expected, readResult)
}

func TestReadMemoryNeedleMetricsNotFound(t *testing.T) {
	expected := Memory{
		Total:         0,
		Available:     0,
		SwapTotal:     0,
		SwapAvailable: 0,
	}
	tempDir := procDir(t, map[string]string{
		"meminfo": "Buffers: 338020 kB\nCached: 1234 kB",
	})
	fs := procfs.New(tempDir)
	reader := MemoryReader{fs}

	readResult, err := reader.Read()

	require.NoError(t, err)
	assert.Equal(t, expected, readResult)
}

func TestReadMemoryFileNotFound(t *testing.T) {
	tempDir := procDir(t, map[string]string{})
	fs := procfs.New(tempDir)
	reader := MemoryReader{fs}

	_, err := reader.Read()

	require.Error(t, err)
}

func TestReadMemoryErrorConvertMetrics(t *testing.T) {
	tempDir := procDir(t, map[string]string{
		"meminfo": "Buffers: 338020 kB\nMemAvailable: is_not_converted_string kB",
	})
	fs := procfs.New(tempDir)
	reader := MemoryReader{fs}

	_, err := reader.Read()

	require.ErrorContains(t, err, "error parsing metric")
}

func TestReadMemoryIncorrectMetricLine(t *testing.T) {
	expected := Memory{
		Total:         0,
		Available:     0,
		SwapTotal:     0,
		SwapAvailable: 0,
	}
	tempDir := procDir(t, map[string]string{
		"meminfo": "Buffers 338020 kB\nMemAvailable 1321223 kB",
	})
	fs := procfs.New(tempDir)
	reader := MemoryReader{fs}

	readResult, err := reader.Read()

	require.NoError(t, err)
	assert.Equal(t, expected, readResult)
}

func TestReadMemoryBlankMetricRow(t *testing.T) {
	expected := Memory{
		Total:         0,
		Available:     0,
		SwapTotal:     0,
		SwapAvailable: 0,
	}
	tempDir := procDir(t, map[string]string{
		"meminfo": "Buffers: 338020 kB\nMemAvailable:",
	})
	fs := procfs.New(tempDir)
	reader := MemoryReader{fs}

	readResult, err := reader.Read()

	require.NoError(t, err)
	assert.Equal(t, expected, readResult)
}
