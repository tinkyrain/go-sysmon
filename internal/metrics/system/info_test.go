package system

//
// import (
// 	"os"
// 	"path/filepath"
// 	"runtime"
// 	"testing"
//
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// 	"github.com/tinkyrain/go-sysmon/internal/procfs"
// )
//
// func tempProcWithFiles(t *testing.T, files map[string]string) string {
// 	t.Helper()
// 	dir := t.TempDir()
// 	for name, content := range files {
// 		path := filepath.Join(dir, name)
// 		require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o755))
// 		require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
// 	}
// 	return dir
// }
//
// func TestReadInfoSuccess(t *testing.T) {
// 	expected := Info{
// 		Hostname: "mak-pc",
// 		OS:       "Linux",
// 		Kernel:   "6.8.0-45-generic",
// 		Arch:     runtime.GOARCH,
// 	}
//
// 	fs := procfs.New("testdata/")
// 	reader := InfoReader{fs: fs}
//
// 	result, err := reader.Read()
//
// 	require.NoError(t, err)
// 	assert.Equal(t, expected, result)
// }
//
// func TestReadInfoFileNotFound(t *testing.T) {
// 	for _, missing := range []string{osTypeFilepath, osReleaseFilepath, hostNameFilepath} {
// 		t.Run(missing, func(t *testing.T) {
// 			files := infoFiles()
// 			delete(files, missing)
//
// 			fs := procfs.New(tempProcWithFiles(t, files))
// 			reader := InfoReader{fs: fs}
//
// 			result, err := reader.Read()
//
// 			require.Error(t, err)
// 			assert.Equal(t, Info{}, result)
// 		})
// 	}
// }
//
// func TestReadInfoReadFileError(t *testing.T) {
// 	files := infoFiles()
// 	delete(files, hostNameFilepath)
// 	dir := tempProcWithFiles(t, files)
//
// 	require.NoError(t, os.Mkdir(filepath.Join(dir, hostNameFilepath), 0o755))
//
// 	reader := InfoReader{fs: procfs.New(dir)}
//
// 	result, err := reader.Read()
//
// 	require.Error(t, err)
// 	assert.Equal(t, Info{}, result)
// }
//
// func TestReadInfoBlankFile(t *testing.T) {
// 	files := infoFiles()
// 	files[hostNameFilepath] = ""
//
// 	reader := InfoReader{fs: procfs.New(tempProcWithFiles(t, files))}
//
// 	result, err := reader.Read()
//
// 	require.NoError(t, err)
// 	assert.Equal(t, "not found", result.Hostname)
// 	assert.Equal(t, "Linux", result.OS)
// }
