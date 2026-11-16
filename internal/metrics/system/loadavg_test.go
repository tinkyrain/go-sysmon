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

	reader := loadAvgReader{fs: procfs.New(tempDirWithFiles(t, loadavgFiles()))}

	result, err := reader.read()

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestReadLoadAvgScanError(t *testing.T) {
	files := loadavgFiles()
	delete(files, loadAvgFile)

	reader := loadAvgReader{fs: procfs.New(tempDirWithFiles(t, files))}

	result, err := reader.read()

	require.ErrorIs(t, err, os.ErrNotExist)
	assert.Equal(t, LoadAvg{}, result)
}

func TestReadLoadAvgEmptyDataError(t *testing.T) {
	files := loadavgFiles()
	files[loadAvgFile] = ""

	reader := loadAvgReader{fs: procfs.New(tempDirWithFiles(t, files))}

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
			name:          "a line with too few fields is rejected",
			line:          "0.25 0.30 1/456",
			expectedError: ErrIncorrectLoadAvgData,
		},
		{
			name:          "a one min value that is not a number is rejected",
			line:          "test 0.25 0.30 1/456 12345",
			expectedError: strconv.ErrSyntax,
		},
		{
			name:          "a five min value that is not a number is rejected",
			line:          "0.15 test 0.30 1/456 12345",
			expectedError: strconv.ErrSyntax,
		},
		{
			name:          "a fifteen min value that is not a number is rejected",
			line:          "0.15 0.25 test 1/456 12345",
			expectedError: strconv.ErrSyntax,
		},
		{
			name:          "a last PID that is not a number is rejected",
			line:          "0.15 0.25 0.30 1/456 test",
			expectedError: strconv.ErrSyntax,
		},
		{
			name:          "a running procs value that is not a number is rejected",
			line:          "0.15 0.25 0.30 test/456 12345",
			expectedError: strconv.ErrSyntax,
		},
		{
			name:          "a total procs value that is not a number is rejected",
			line:          "0.15 0.25 0.30 1/test 12345",
			expectedError: strconv.ErrSyntax,
		},
		{
			name:          "a procs field without a slash is rejected",
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
