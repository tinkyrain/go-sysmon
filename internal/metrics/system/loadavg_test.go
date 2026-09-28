package system

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tinkyrain/go-sysmon/internal/procfs"
)

func loadavgFiles() map[string]string {
	return map[string]string{
		"loadavg": "0.15 0.25 0.30 1/456 12345",
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
	fs := procfs.New(TempDirWithFiles(t, loadavgFiles(), 0o755, 0o600))
	reader := LoadAvgReader{fs: fs}

	result, err := reader.Read()

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestReadLoadAvgReadFileError(t *testing.T) {
	expected := LoadAvg{}

	filepath := TempDirWithFiles(t, loadavgFiles(), 0, 0)
	fs := procfs.New(filepath)
	reader := LoadAvgReader{fs: fs}

	result, err := reader.Read()

	require.Error(t, err)
	assert.Equal(t, expected, result)
}

func TestReadLoadAvgBlankFileError(t *testing.T) {
	expected := LoadAvg{}

	files := loadavgFiles()
	files["loadavg"] = ""

	filepath := TempDirWithFiles(t, files, 0o755, 0o600)
	fs := procfs.New(filepath)
	reader := LoadAvgReader{fs: fs}

	result, err := reader.Read()

	require.ErrorIs(t, err, ErrInsufficientLoadAvgData)
	assert.Equal(t, expected, result)
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

	require.ErrorIs(t, err, ErrInsufficientLoadAvgData)
	require.Equal(t, LoadAvg{}, result)
}

func TestParseLoadAvgLineParseMetricsError(t *testing.T) {
	data := []string{
		"test 0.25 0.30 1/456 12345",
		"0.15 test 0.30 1/456 12345",
		"0.15 0.25 test 1/456 12345",
		"0.15 0.25 0.30 1/456 test",
		"0.15 0.25 0.30 test/456 12345",
		"0.15 0.25 0.30 1/test 12345",
		"0.15 0.25 0.30 test 12345",
	}

	for _, line := range data {
		t.Run(line, func(t *testing.T) {
			result, err := parseLoadAvgLine(line)

			require.Error(t, err)
			require.Equal(t, LoadAvg{}, result)
		})
	}
}
