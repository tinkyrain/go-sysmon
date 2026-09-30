package system

import (
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tinkyrain/go-sysmon/internal/procfs"
)

func loadavgFiles() map[string]string {
	return map[string]string{
		loadAvgFile: "0.15 0.25 0.30 1/456 12345",
	}
}

func TestReadLoadAvgSuccess(t *testing.T) {
	expected := LoadAvg{
		OneMin:       0.15,
		FiveMin:      0.25,
		FifteenMin:   0.30,
		RunningProcs: 1,
		TotalProcs:   456,
		LastPID:      12345,
	}

	reader := loadAvgReader{fs: procfs.New(tempDirWithFiles(t, loadavgFiles(), 0o755, 0o600))}

	result, err := reader.read()

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestReadLoadAvgReadFileError(t *testing.T) {
	dir := tempDirWithFiles(t, map[string]string{}, 0o755, 0o600)
	require.NoError(t, os.Mkdir(filepath.Join(dir, loadAvgFile), 0o755))

	reader := loadAvgReader{fs: procfs.New(dir)}

	result, err := reader.read()

	require.Error(t, err)
	assert.Equal(t, LoadAvg{}, result)
}

func TestReadLoadAvgBlankFileError(t *testing.T) {
	files := loadavgFiles()
	files[loadAvgFile] = ""

	reader := loadAvgReader{fs: procfs.New(tempDirWithFiles(t, files, 0o755, 0o600))}

	result, err := reader.read()

	require.ErrorIs(t, err, ErrInsufficientLoadAvg)
	assert.Equal(t, LoadAvg{}, result)
}

func TestParseLoadAvgLineSuccess(t *testing.T) {
	expected := LoadAvg{
		OneMin:       0.15,
		FiveMin:      0.25,
		FifteenMin:   0.30,
		RunningProcs: 1,
		TotalProcs:   456,
		LastPID:      12345,
	}

	result, err := parseLoadAvgLine("0.15 0.25 0.30 1/456 12345")

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestParseLoadAvgLineIncorrectMetricsCountErr(t *testing.T) {
	result, err := parseLoadAvgLine("123")

	require.ErrorIs(t, err, ErrInsufficientLoadAvg)
	require.Equal(t, LoadAvg{}, result)
}

func TestParseLoadAvgLineIncorrectMetricError(t *testing.T) {
	result, err := parseLoadAvgLine("0.15 0.25 0.30 test 12345")

	require.ErrorIs(t, err, ErrIncorrectLoadAvgMetric)
	require.Equal(t, LoadAvg{}, result)
}

func TestParseLoadAvgLineParseMetricsError(t *testing.T) {
	data := map[string]string{
		"test_1": "test 0.25 0.30 1/456 12345",
		"test_2": "0.15 test 0.30 1/456 12345",
		"test_3": "0.15 0.25 test 1/456 12345",
		"test_4": "0.15 0.25 0.30 1/456 test",
		"test_5": "0.15 0.25 0.30 test/456 12345",
		"test_6": "0.15 0.25 0.30 1/test 12345",
	}

	for key, line := range data {
		t.Run(key, func(t *testing.T) {
			result, err := parseLoadAvgLine(line)

			require.ErrorIs(t, err, strconv.ErrSyntax)
			require.Equal(t, LoadAvg{}, result)
		})
	}
}
