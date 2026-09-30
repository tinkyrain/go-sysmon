package cpu

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func tempDirWithFiles(
	t *testing.T,
	files map[string]string,
	dirPerm os.FileMode,
	filePerm os.FileMode,
) string {
	t.Helper()
	dir := t.TempDir()
	for name, content := range files {
		path := filepath.Join(dir, name)
		require.NoError(t, os.MkdirAll(filepath.Dir(path), dirPerm))
		require.NoError(t, os.WriteFile(path, []byte(content), filePerm))
	}
	return dir
}
