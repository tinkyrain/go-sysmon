package cpu

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tinkyrain/go-sysmon/internal/procfs"
)

func usageFiles() map[string]string {
	return map[string]string{
		statFile: `
cpu  298830 3399 75896 7741115 5083 0 1586 0 0 0
cpu0 17496 72 6050 650357 437 0 376 0 0 0
cpu1 10757 57 3837 661024 159 0 1073 0 0 0
intr 30062225 138 0 0 0 0 0 0 1 0 7 0 0 0 0 0 0 0 0 0
`,
	}
}

func TestReadUsageSuccess(t *testing.T) {
	cases := []struct {
		name                     string
		prevUsageSamples         map[string]usageSample
		expectedUsage            Usage
		expectedPrevUsageSamples map[string]usageSample
	}{
		{
			name:             "Success read without prev samples",
			prevUsageSamples: map[string]usageSample{},
			expectedUsage: Usage{
				Total: 0,
				Cores: []CoreUsage{
					{ID: "cpu0", Usage: 0},
					{ID: "cpu1", Usage: 0},
				},
			},
			expectedPrevUsageSamples: map[string]usageSample{
				"cpu": {
					total:   8125909,
					idle:    7746198,
					nonIdle: 379711,
				},
				"cpu0": {
					total:   674788,
					idle:    650794,
					nonIdle: 23994,
				},
				"cpu1": {
					total:   676907,
					idle:    661183,
					nonIdle: 15724,
				},
			},
		},
		{
			name: "Success read with prev samples",
			prevUsageSamples: map[string]usageSample{
				"cpu": {
					total:   7746198,
					idle:    6746198,
					nonIdle: 279711,
				},
				"cpu0": {
					total:   650000,
					idle:    630000,
					nonIdle: 20000,
				},
				"cpu1": {
					total:   600000,
					idle:    590000,
					nonIdle: 10000,
				},
			},
			expectedUsage: Usage{
				Total: 9.090909090909092,
				Cores: []CoreUsage{
					{ID: "cpu0", Usage: 16.112635146038407},
					{ID: "cpu1", Usage: 7.442755535907004},
				},
			},
			expectedPrevUsageSamples: map[string]usageSample{
				"cpu": {
					total:   8125909,
					idle:    7746198,
					nonIdle: 379711,
				},
				"cpu0": {
					total:   674788,
					idle:    650794,
					nonIdle: 23994,
				},
				"cpu1": {
					total:   676907,
					idle:    661183,
					nonIdle: 15724,
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reader := usageReader{
				fs:               procfs.New(tempDirWithFiles(t, usageFiles())),
				prevUsageSamples: tc.prevUsageSamples,
			}

			result, err := reader.read()

			require.NoError(t, err)
			assert.Equal(t, tc.expectedUsage, result)
			assert.Equal(t, tc.expectedPrevUsageSamples, reader.prevUsageSamples)
		})
	}
}

func TestReadUsageScanError(t *testing.T) {
	files := usageFiles()
	delete(files, statFile)

	reader := usageReader{
		fs:               procfs.New(tempDirWithFiles(t, files)),
		prevUsageSamples: map[string]usageSample{},
	}

	result, err := reader.read()

	require.ErrorIs(t, err, os.ErrNotExist)
	assert.Equal(t, Usage{}, result)
}

func TestReadUsageEmptyDataError(t *testing.T) {
	cases := []struct {
		name        string
		fileContent string
	}{
		{
			name:        "Empty file",
			fileContent: "",
		},
		{
			name:        "Without cpu lines",
			fileContent: "intr 30062225 138 0 0 0 0 0 0 1 0 7 0 0 0 0 0 0 0 0 0",
		},
		{
			name:        "Cpu lines with insufficient ticks",
			fileContent: "cpu  298830 3399 75896\ncpu0 17496 72 6050",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			files := usageFiles()
			files[statFile] = tc.fileContent

			reader := usageReader{
				fs:               procfs.New(tempDirWithFiles(t, files)),
				prevUsageSamples: map[string]usageSample{},
			}

			result, err := reader.read()

			require.ErrorIs(t, err, ErrEmptyStatData)
			assert.Equal(t, Usage{}, result)
		})
	}
}

func TestReadUsageConvertTicksError(t *testing.T) {
	files := usageFiles()
	files[statFile] = "cpu  298830 test 75896 7741115 5083 0 1586 0 0 0"

	reader := usageReader{
		fs:               procfs.New(tempDirWithFiles(t, files)),
		prevUsageSamples: map[string]usageSample{},
	}

	result, err := reader.read()

	require.ErrorIs(t, err, strconv.ErrSyntax)
	assert.Equal(t, Usage{}, result)
}

func TestReadUsageSkipsUnsuitableLines(t *testing.T) {
	expected := Usage{
		Total: 0,
		Cores: []CoreUsage{
			{ID: "cpu0", Usage: 0},
		},
	}

	files := usageFiles()
	files[statFile] = `cpu 28683 104 7543 6
cpu0 17496 72 6050 650357 437 0 376 0 0 0
intr 30062225 138 0 0 0 0 0 0 1 0 7 0 0 0 0 0 0 0 0 0`

	reader := usageReader{
		fs:               procfs.New(tempDirWithFiles(t, files)),
		prevUsageSamples: map[string]usageSample{},
	}

	result, err := reader.read()

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestParseTicksLineSuccess(t *testing.T) {
	expected := ticks{
		id:        "cpu0",
		user:      17496,
		nice:      72,
		system:    6050,
		idle:      650357,
		iowait:    437,
		irq:       0,
		softirq:   376,
		steal:     0,
		guest:     0,
		guestNice: 0,
	}

	result, err := parseTicksLine([]string{
		"cpu0", "17496", "72", "6050", "650357", "437", "0", "376", "0", "0", "0",
	})

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestParseTicksLineError(t *testing.T) {
	cases := []struct {
		name          string
		line          []string
		expectedError error
	}{
		{
			name:          "Insufficient ticks in the line",
			line:          []string{"cpu0", "17496", "72", "6050"},
			expectedError: ErrIncorrectStatData,
		},
		{
			name: "Convert tick value to uint error",
			line: []string{
				"cpu0", "17496", "test", "6050", "650357", "437", "0", "376", "0", "0", "0",
			},
			expectedError: strconv.ErrSyntax,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := parseTicksLine(tc.line)

			require.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, ticks{}, result)
		})
	}
}

func TestCalculateUsage(t *testing.T) {
	cases := []struct {
		name     string
		prev     usageSample
		cur      usageSample
		expected float64
	}{
		{
			name:     "Half of the ticks is busy",
			prev:     usageSample{total: 500, idle: 200, nonIdle: 300},
			cur:      usageSample{total: 600, idle: 250, nonIdle: 350},
			expected: 50,
		},
		{
			name:     "Prev sample is blank",
			prev:     usageSample{},
			cur:      usageSample{total: 600, idle: 250, nonIdle: 350},
			expected: 0,
		},
		{
			name:     "Idle ticks decreased",
			prev:     usageSample{total: 550, idle: 250, nonIdle: 300},
			cur:      usageSample{total: 550, idle: 200, nonIdle: 350},
			expected: 0,
		},
		{
			name:     "Non idle ticks decreased",
			prev:     usageSample{total: 550, idle: 200, nonIdle: 350},
			cur:      usageSample{total: 550, idle: 250, nonIdle: 300},
			expected: 0,
		},
		{
			name:     "Only idle ticks grew",
			prev:     usageSample{total: 500, idle: 200, nonIdle: 300},
			cur:      usageSample{total: 550, idle: 250, nonIdle: 300},
			expected: 0,
		},
		{
			name:     "Ticks did not change",
			prev:     usageSample{total: 500, idle: 200, nonIdle: 300},
			cur:      usageSample{total: 500, idle: 200, nonIdle: 300},
			expected: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := calculateUsage(tc.prev, tc.cur)

			assert.Equal(t, tc.expected, result)
		})
	}
}
