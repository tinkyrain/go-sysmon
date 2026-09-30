package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tinkyrain/go-sysmon/internal/procfs"
)

func cpuFiles() map[string]string {
	return map[string]string{
		statFile: `
cpu  298830 3399 75896 7741115 5083 0 1586 0 0 0
cpu0 17496 72 6050 650357 437 0 376 0 0 0
cpu1 10757 57 3837 661024 159 0 1073 0 0 0
intr 30062225 138 0 0 0 0 0 0 1 0 7 0 0 0 0 0 0 0 0 0
`,
	}
}

func TestCollectSuccess(t *testing.T) {
	expected := Stats{
		Usage: Usage{
			Total: 0,
			Cores: []CoreUsage{
				{ID: "cpu0", Usage: 0},
				{ID: "cpu1", Usage: 0},
			},
		},
	}
	collector := New(procfs.New(tempDirWithFiles(t, cpuFiles(), 0o755, 0o600)))

	result, err := collector.Collect()

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestCollectError(t *testing.T) {
	for file := range cpuFiles() {
		t.Run(file, func(t *testing.T) {
			files := cpuFiles()
			delete(files, file)

			collector := New(procfs.New(tempDirWithFiles(t, files, 0o755, 0o600)))

			result, err := collector.Collect()

			require.Error(t, err)
			assert.Equal(t, Stats{}, result)
		})
	}
}
