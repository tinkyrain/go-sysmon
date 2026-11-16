package memory

import (
	"os"
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

func TestReadMemInfoSuccess(t *testing.T) {
	expected := Info{
		Total:         13315522560,
		Available:     5244837888,
		SwapTotal:     4294963200,
		SwapAvailable: 3464392704,
	}

	reader := meminfoReader{procfs.New(tempDirWithFiles(t, meminfoFiles()))}

	result, err := reader.read()

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestReadMemInfoScanError(t *testing.T) {
	files := meminfoFiles()
	delete(files, meminfoFile)

	reader := meminfoReader{fs: procfs.New(tempDirWithFiles(t, files))}

	result, err := reader.read()

	require.ErrorIs(t, err, os.ErrNotExist)
	assert.Equal(t, Info{}, result)
}

func TestReadMemInfoError(t *testing.T) {
	cases := []struct {
		name          string
		fileContent   string
		expectedError error
	}{
		{
			name:          "an empty file has no memory data",
			fileContent:   "",
			expectedError: ErrEmptyMemInfoData,
		},
		{
			name:          "a value that is not a number is rejected",
			fileContent:   "MemTotal:	testme",
			expectedError: strconv.ErrSyntax,
		},
		{
			name:          "MemTotal missing is rejected",
			fileContent:   "MemAvailable: 200 kB\nSwapTotal: 10 kB\nSwapFree: 3 kB",
			expectedError: ErrIncorrectMemInfoData,
		},
		{
			name:          "MemAvailable missing is rejected",
			fileContent:   "MemTotal: 210 kB\nSwapTotal: 10 kB\nSwapFree: 3 kB",
			expectedError: ErrIncorrectMemInfoData,
		},
		{
			name:          "SwapTotal missing is rejected",
			fileContent:   "MemTotal: 210 kB\nMemAvailable: 200 kB\nSwapFree: 3 kB",
			expectedError: ErrIncorrectMemInfoData,
		},
		{
			name:          "SwapFree missing is rejected",
			fileContent:   "MemTotal: 210 kB\nMemAvailable: 200 kB\nSwapTotal: 10 kB",
			expectedError: ErrIncorrectMemInfoData,
		},
		{
			name:          "keys without values are rejected",
			fileContent:   "MemTotal:\nMemAvailable:\nSwapTotal:\nSwapFree:",
			expectedError: ErrIncorrectMemInfoData,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			files := meminfoFiles()
			files[meminfoFile] = tc.fileContent

			reader := meminfoReader{fs: procfs.New(tempDirWithFiles(t, files))}

			result, err := reader.read()

			require.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, Info{}, result)
		})
	}
}
