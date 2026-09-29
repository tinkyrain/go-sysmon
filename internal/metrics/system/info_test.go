package system

import (
	"os"
	"path/filepath"
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

	fs := procfs.New(tempDirWithFiles(t, files, 0o755, 0o600))
	reader := infoReader{fs: fs}

	result, err := reader.read()

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestReadInfoFileNotFound(t *testing.T) {
	for _, missing := range []string{ostypeFile, osreleaseFile, hostnameFile} {
		t.Run(missing, func(t *testing.T) {
			files := infoFiles()
			delete(files, missing)

			fs := procfs.New(tempDirWithFiles(t, files, 0o755, 0o600))
			reader := infoReader{fs: fs}

			result, err := reader.read()

			require.Error(t, err)
			assert.Equal(t, Info{}, result)
		})
	}
}

func TestReadInfoReadFileError(t *testing.T) {
	files := infoFiles()
	delete(files, hostnameFile)
	dir := tempDirWithFiles(t, files, 0o755, 0o600)

	require.NoError(t, os.Mkdir(filepath.Join(dir, hostnameFile), 0o755))

	reader := infoReader{fs: procfs.New(dir)}

	result, err := reader.read()

	require.Error(t, err)
	assert.Equal(t, Info{}, result)
}

func TestReadInfoBlankFile(t *testing.T) {
	files := infoFiles()
	files[hostnameFile] = ""

	reader := infoReader{fs: procfs.New(tempDirWithFiles(t, files, 0o755, 0o600))}

	result, err := reader.read()

	require.NoError(t, err)
	assert.Equal(t, "not found", result.Hostname)
	assert.Equal(t, "ostype", result.OS)
	assert.Equal(t, "osrelease", result.Kernel)
}
