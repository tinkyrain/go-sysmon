package statfs_test

import (
	"go-sysmon/internal/statfs"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDirStatfs(t *testing.T) {
	stat, err := statfs.GetDirStatfs(t.TempDir())

	require.NoError(t, err)
	assert.Positive(t, stat.Bsize)
	assert.Positive(t, stat.Blocks)
	assert.LessOrEqual(t, stat.Bavail, stat.Blocks)
}

func TestGetDirStatfsOnFile(t *testing.T) {
	file := filepath.Join(t.TempDir(), "f")
	require.NoError(t, os.WriteFile(file, nil, 0o600))

	stat, err := statfs.GetDirStatfs(file)

	require.NoError(t, err)
	assert.Positive(t, stat.Blocks)
}

func TestGetDirStatfsNotFound(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nope")

	_, err := statfs.GetDirStatfs(missing)

	require.Error(t, err)
	assert.ErrorIs(t, err, fs.ErrNotExist)
	assert.ErrorContains(t, err, missing)
}
