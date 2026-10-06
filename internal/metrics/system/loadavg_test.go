package system

import (
	"os"
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

func TestReadLoadAvgScanError(t *testing.T) {
	files := loadavgFiles()
	delete(files, loadAvgFile)

	reader := loadAvgReader{fs: procfs.New(tempDirWithFiles(t, files, 0o755, 0o600))}

	result, err := reader.read()

	require.ErrorIs(t, err, os.ErrNotExist)
	assert.Equal(t, LoadAvg{}, result)
}

func TestReadLoadAvgEmptyDataError(t *testing.T) {
	files := loadavgFiles()
	files[loadAvgFile] = ""

	reader := loadAvgReader{fs: procfs.New(tempDirWithFiles(t, files, 0o755, 0o600))}

	result, err := reader.read()

	require.ErrorIs(t, err, ErrEmptyLoadAvgData)
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

func TestParseLoadAvgLineError(t *testing.T) {
	cases := []struct {
		name          string
		line          string
		expectedError error
	}{
		{
			name:          "Incorrect data in the line",
			line:          "0.25 0.30 1/456",
			expectedError: ErrIncorrectLoadAvgData,
		},
		{
			name:          "One min metric is incorrect",
			line:          "test 0.25 0.30 1/456 12345",
			expectedError: strconv.ErrSyntax,
		},
		{
			name:          "Five min metric is incorrect",
			line:          "0.15 test 0.30 1/456 12345",
			expectedError: strconv.ErrSyntax,
		},
		{
			name:          "Fifteen min metric is incorrect",
			line:          "0.15 0.25 test 1/456 12345",
			expectedError: strconv.ErrSyntax,
		},
		{
			name:          "Lastpid metric is incorrect",
			line:          "0.15 0.25 0.30 1/456 test",
			expectedError: strconv.ErrSyntax,
		},
		{
			name:          "Running procs metric is incorrect",
			line:          "0.15 0.25 0.30 test/456 12345",
			expectedError: strconv.ErrSyntax,
		},
		{
			name:          "Total procs metric is incorrect",
			line:          "0.15 0.25 0.30 1/test 12345",
			expectedError: strconv.ErrSyntax,
		},
		{
			name:          "Procs metric string is incorrect",
			line:          "0.15 0.25 0.30 test 12345",
			expectedError: ErrIncorrectLoadAvgData,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := parseLoadAvgLine(tc.line)

			require.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, LoadAvg{}, result)
		})
	}
}
