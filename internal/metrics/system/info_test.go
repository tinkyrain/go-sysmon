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
		hostNameFilepath:  "hostname",
		osReleaseFilepath: "osrelease",
		osTypeFilepath:    "ostype",
	}
}

func TestReadInfoSuccess(t *testing.T) {
	files := infoFiles()
	expected := Info{
		Hostname: files[hostNameFilepath],
		OS:       files[osTypeFilepath],
		Kernel:   files[osReleaseFilepath],
		Arch:     runtime.GOARCH,
	}

	fs := procfs.New(TempDirWithFiles(t, files, 0o755, 0o600))
	reader := InfoReader{fs: fs}

	result, err := reader.Read()

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestReadInfoFileNotFound(t *testing.T) {
	for _, missing := range []string{osTypeFilepath, osReleaseFilepath, hostNameFilepath} {
		t.Run(missing, func(t *testing.T) {
			files := infoFiles()
			delete(files, missing)

			fs := procfs.New(TempDirWithFiles(t, files, 0o755, 0o600))
			reader := InfoReader{fs: fs}

			result, err := reader.Read()

			require.Error(t, err)
			assert.Equal(t, Info{}, result)
		})
	}
}

func TestReadInfoReadFileError(t *testing.T) {
	files := infoFiles()
	delete(files, hostNameFilepath)
	dir := TempDirWithFiles(t, files, 0o755, 0o600)

	require.NoError(t, os.Mkdir(filepath.Join(dir, hostNameFilepath), 0o755))

	reader := InfoReader{fs: procfs.New(dir)}

	result, err := reader.Read()

	require.Error(t, err)
	assert.Equal(t, Info{}, result)
}

func TestReadInfoBlankFile(t *testing.T) {
	files := infoFiles()
	files[hostNameFilepath] = ""

	reader := InfoReader{fs: procfs.New(TempDirWithFiles(t, files, 0o755, 0o600))}

	result, err := reader.Read()

	require.NoError(t, err)
	assert.Equal(t, "not found", result.Hostname)
	assert.Equal(t, "ostype", result.OS)
	assert.Equal(t, "osrelease", result.Kernel)
}
