package procfs_test

import (
	"go-sysmon/internal/procfs"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getProcfsScanner(path string) procfs.FileScanner {
	return procfs.New(path)
}

func TestProcfsSuccessScan(t *testing.T) {
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

func TestProcfsErrorScan(t *testing.T) {
	scanner := getProcfsScanner("testdata/")
	_, err := scanner.ScanRows("error_scan")
	assert.Error(t, err)
}
