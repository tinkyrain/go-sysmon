package statfs_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/tinkyrain/go-sysmon/internal/statfs"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {
	stat, err := statfs.Read(t.TempDir())

	require.NoError(t, err)
	assert.Positive(t, stat.BlockSize)
	assert.Positive(t, stat.Blocks)
	assert.LessOrEqual(t, stat.Available, stat.Blocks)
}

func TestReadOnFile(t *testing.T) {
	file := filepath.Join(t.TempDir(), "f")
	require.NoError(t, os.WriteFile(file, nil, 0o600))

	stat, err := statfs.Read(file)

	require.NoError(t, err)
	assert.Positive(t, stat.Blocks)
}

func TestReadNotFound(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nope")
	_, err := statfs.Read(missing)
	require.ErrorIs(t, err, fs.ErrNotExist)
}
