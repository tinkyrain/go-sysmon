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

func TestProcfsReadLinesSuccess(t *testing.T) {
	cases := []struct {
		name     string
		filename string
		content  string
		expected []string
	}{
		{
			name:     "a single line comes back as one element",
			filename: "file",
			content:  "Lorem ipsum dolor sit amet, consectetur adipiscing elit",
			expected: []string{"Lorem ipsum dolor sit amet, consectetur adipiscing elit"},
		},
		{
			name:     "every line becomes its own element",
			filename: "file",
			content:  "Lorem ipsum\ndolor sit amet",
			expected: []string{
				"Lorem ipsum",
				"dolor sit amet",
			},
		},
		{
			name:     "an empty file reads as an empty slice",
			filename: "file",
			content:  "",
			expected: []string{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			require.NoError(t, os.WriteFile(filepath.Join(dir, tc.filename), []byte(tc.content), 0o600))

			reader := procfs.New(dir)
			lines, err := reader.ReadLines(tc.filename)

			require.NoError(t, err)
			assert.Equal(t, tc.expected, lines)
		})
	}
}

func TestProcfsReadLinesError(t *testing.T) {
	cases := []struct {
		name           string
		targetFilename string
		filename       string
		content        string
		expectedResult []string
		expectedError  error
	}{
		{
			name:           "a missing file is an error",
			filename:       "file",
			targetFilename: "file_not_found",
			expectedResult: nil,
			expectedError:  os.ErrNotExist,
		},
		{
			name:           "a line over 64 KB is an error",
			filename:       "file",
			targetFilename: "file",
			content:        strings.Repeat("A", 65*1024),
			expectedResult: nil,
			expectedError:  bufio.ErrTooLong,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			require.NoError(t, os.WriteFile(filepath.Join(dir, tc.filename), []byte(tc.content), 0o600))

			reader := procfs.New(dir)
			result, err := reader.ReadLines(tc.targetFilename)

			require.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestProcfsReadLineSuccess(t *testing.T) {
	cases := []struct {
		name     string
		filename string
		content  string
		expected string
	}{
		{
			name:     "a single line comes back whole",
			filename: "file",
			content:  "Lorem ipsum dolor sit amet, consectetur adipiscing elit",
			expected: "Lorem ipsum dolor sit amet, consectetur adipiscing elit",
		},
		{
			name:     "only the first line is returned",
			filename: "file",
			content:  "Lorem ipsum\ndolor sit amet",
			expected: "Lorem ipsum",
		},
		{
			name:     "an empty file reads as an empty string",
			filename: "file",
			content:  "",
			expected: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			require.NoError(t, os.WriteFile(filepath.Join(dir, tc.filename), []byte(tc.content), 0o600))

			reader := procfs.New(dir)
			line, err := reader.ReadLine(tc.filename)

			require.NoError(t, err)
			assert.Equal(t, tc.expected, line)
		})
	}
}

func TestProcfsReadLineError(t *testing.T) {
	cases := []struct {
		name           string
		targetFilename string
		filename       string
		content        string
		expectedResult string
		expectedError  error
	}{
		{
			name:           "a missing file is an error",
			filename:       "file",
			targetFilename: "file_not_found",
			expectedResult: "",
			expectedError:  os.ErrNotExist,
		},
		{
			name:           "a line over 64 KB is an error",
			filename:       "file",
			targetFilename: "file",
			content:        strings.Repeat("A", 65*1024),
			expectedResult: "",
			expectedError:  bufio.ErrTooLong,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			require.NoError(t, os.WriteFile(filepath.Join(dir, tc.filename), []byte(tc.content), 0o600))

			reader := procfs.New(dir)
			result, err := reader.ReadLine(tc.targetFilename)

			require.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}
