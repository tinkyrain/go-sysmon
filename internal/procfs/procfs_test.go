package procfs_test

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tinkyrain/go-sysmon/internal/procfs"
)

func TestProcfsScanSuccess(t *testing.T) {
	cases := []struct {
		name     string
		filename string
		content  string
		expected []string
	}{
		{
			name:     "Test case: successfully complete reading the file",
			filename: "file",
			content:  "Lorem ipsum dolor sit amet, consectetur adipiscing elit",
			expected: []string{"Lorem ipsum dolor sit amet, consectetur adipiscing elit"},
		},
		{
			name:     "Test case: successfully complete reading the file with many rows",
			filename: "file",
			content:  "Lorem ipsum\ndolor sit amet",
			expected: []string{
				"Lorem ipsum",
				"dolor sit amet",
			},
		},
		{
			name:     "Test case: successfully complete reading blank the file",
			filename: "file",
			content:  "",
			expected: []string{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			require.NoError(t, os.WriteFile(filepath.Join(dir, tc.filename), []byte(tc.content), 0o600))

			scanner := procfs.New(dir)
			rows, err := scanner.ScanRows(tc.filename)

			require.NoError(t, err)
			assert.Equal(t, tc.expected, rows)
		})
	}
}

func TestProcfsScanError(t *testing.T) {
	cases := []struct {
		name               string
		targetFilename     string
		filename           string
		content            string
		expectedScanResult []string
		expectedError      error
	}{
		{
			name:               "Test case: file not found",
			filename:           "file",
			targetFilename:     "file_not_found",
			expectedScanResult: nil,
			expectedError:      os.ErrNotExist,
		},
		{
			name:               "Test case: file with row more 64 KB",
			filename:           "file",
			targetFilename:     "file",
			content:            strings.Repeat("A", 65*1024),
			expectedScanResult: nil,
			expectedError:      bufio.ErrTooLong,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			require.NoError(t, os.WriteFile(filepath.Join(dir, tc.filename), []byte(tc.content), 0o600))

			scanner := procfs.New(dir)
			result, err := scanner.ScanRows(tc.targetFilename)

			require.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, tc.expectedScanResult, result)
		})
	}
}
