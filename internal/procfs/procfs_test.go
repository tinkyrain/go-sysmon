package procfs_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tinkyrain/go-sysmon/internal/procfs"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newScanner(t *testing.T, file, filecontent string) procfs.FileScanner {
	t.Helper()
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, file), []byte(filecontent), 0o600))
	return procfs.New(dir)
}

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
			scanner := newScanner(t, tc.filename, tc.content)
			rows, err := scanner.ScanRows(tc.filename)

			require.NoError(t, err)
			assert.Equal(t, tc.expected, rows)
		})
	}
}

// func TestProcfsScanError(t *testing.T) {
// 	filename := "error_scan"
// 	path := procfsDir(t, filename, "123123")
// 	scanner := getProcfsScanner(path)
//
// 	_, err := scanner.ScanRows(filename)
//
// 	assert.Error(t, err)
// }
//
// func TestProcfsScanFileNotFound(t *testing.T) {
// 	path := procfsDir(t, "file_not_found", "")
// 	scanner := getProcfsScanner(path)
//
// 	_, err := scanner.ScanRows("fff_not_fff")
//
// 	assert.Error(t, err)
// }
