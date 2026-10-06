package system

import (
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tinkyrain/go-sysmon/internal/procfs"
)

func infoFiles() map[string]string {
	return map[string]string{
		hostnameFile:  "hostname",
		osreleaseFile: "osrelease",
		ostypeFile:    "ostype",
	}
}

func TestReadInfoSuccess(t *testing.T) {
	files := infoFiles()
	expected := Info{
		Hostname: files[hostnameFile],
		OS:       files[ostypeFile],
		Kernel:   files[osreleaseFile],
		Arch:     runtime.GOARCH,
	}

	reader := infoReader{fs: procfs.New(tempDirWithFiles(t, files))}

	result, err := reader.read()

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestReadInfoFileNotExistError(t *testing.T) {
	cases := []struct {
		name       string
		targetFile string
	}{
		{
			name:       "File hostname not found",
			targetFile: hostnameFile,
		},
		{
			name:       "File osrelease not found",
			targetFile: osreleaseFile,
		},
		{
			name:       "File ostype not found",
			targetFile: ostypeFile,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			files := infoFiles()
			delete(files, tc.targetFile)

			reader := infoReader{fs: procfs.New(tempDirWithFiles(t, files))}
			result, err := reader.read()

			require.ErrorIs(t, err, os.ErrNotExist)
			assert.Equal(t, Info{}, result)
		})
	}
}

func TestReadInfoFileIsEmptyError(t *testing.T) {
	cases := []struct {
		name       string
		targetFile string
	}{
		{
			name:       "File hostname is empty",
			targetFile: hostnameFile,
		},
		{
			name:       "File osrelease is empty",
			targetFile: osreleaseFile,
		},
		{
			name:       "File ostype is empty",
			targetFile: ostypeFile,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			files := infoFiles()
			files[tc.targetFile] = ""

			reader := infoReader{fs: procfs.New(tempDirWithFiles(t, files))}
			result, err := reader.read()

			require.ErrorIs(t, err, ErrEmptyInfoData)
			assert.Equal(t, Info{}, result)
		})
	}
}
