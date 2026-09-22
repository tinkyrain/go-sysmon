package procfs_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tinkyrain/go-sysmon/internal/procfs"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func procfsDir(t *testing.T, file, filecontent string, perm os.FileMode) string {
	t.Helper()
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, file), []byte(filecontent), perm))
	return dir
}

func getProcfsScanner(path string) procfs.FileScanner {
	return procfs.New(path)
}

func TestProcfsScanSuccess(t *testing.T) {
	scanner := getProcfsScanner("testdata/")
	expected := []string{
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit",
		"sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		"Duis aute irure dolor in reprehenderit in voluptate velit",
	}

	rows, err := scanner.ScanRows("success_file")

	require.NoError(t, err)
	assert.Equal(t, expected, rows)
}

func TestProcfsScanError(t *testing.T) {
	filename := "error_scan"
	path := procfsDir(t, filename, "123123", 0)
	scanner := getProcfsScanner(path)

	_, err := scanner.ScanRows(filename)

	assert.Error(t, err)
}

func TestProcfsScanFileNotFound(t *testing.T) {
	path := procfsDir(t, "file_not_found", "", 0o600)
	scanner := getProcfsScanner(path)

	_, err := scanner.ScanRows("fff_not_fff")

	assert.Error(t, err)
}

func TestProcfsScanBlankFile(t *testing.T) {
	filename := "blank_file"
	path := procfsDir(t, filename, "", 0o600)
	scanner := getProcfsScanner(path)

	rows, err := scanner.ScanRows(filename)

	require.NoError(t, err)
	assert.Equal(t, []string{}, rows)
}
