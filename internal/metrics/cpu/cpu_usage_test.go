package cpu

import (
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

func TestReadUsageReaderSuccessWithoutPrevSample(t *testing.T) {
	expectedUsage := Usage{
		Total: 0,
		Cores: []CoreUsage{
			{ID: "cpu0", Usage: 0},
			{ID: "cpu1", Usage: 0},
		},
	}
	expectedPrevUsageSamples := map[string]usageSample{
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
	}

	reader := usageReader{
		fs:               procfs.New(tempDirWithFiles(t, usageFiles(), 0o755, 0o600)),
		prevUsageSamples: map[string]usageSample{},
	}

	result, err := reader.read()

	require.NoError(t, err)
	assert.Equal(t, expectedUsage, result)
	assert.Equal(t, expectedPrevUsageSamples, reader.prevUsageSamples)
}

func TestUsageReaderSuccessWithPrevSample(t *testing.T) {
	expectedUsage := Usage{
		Total: 9.090909090909092,
		Cores: []CoreUsage{
			{ID: "cpu0", Usage: 16.112635146038407},
			{ID: "cpu1", Usage: 7.442755535907004},
		},
	}
	expectedPrevUsageSamples := map[string]usageSample{
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
	}

	reader := usageReader{
		fs: procfs.New(tempDirWithFiles(t, usageFiles(), 0o755, 0o600)),
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
	}

	result, err := reader.read()

	require.NoError(t, err)
	assert.Equal(t, expectedUsage, result)
	assert.Equal(t, expectedPrevUsageSamples, reader.prevUsageSamples)
}

func TestUsageReaderWithoutNeedleLines(t *testing.T) {
	files := usageFiles()
	files[statFile] = ""

	reader := usageReader{
		fs:               procfs.New(tempDirWithFiles(t, files, 0o755, 0o600)),
		prevUsageSamples: map[string]usageSample{},
	}

	result, err := reader.read()

	assert.ErrorIs(t, err, ErrNoTickLines)
	assert.Equal(t, Usage{}, result)
}

func TestUsageReaderFileNotFound(t *testing.T) {
	files := usageFiles()
	delete(files, statFile)

	reader := usageReader{
		fs:               procfs.New(tempDirWithFiles(t, files, 0o755, 0o600)),
		prevUsageSamples: map[string]usageSample{},
	}

	result, err := reader.read()

	require.Error(t, err)
	assert.Equal(t, Usage{}, result)
}

func TestUsageReaderFilterIncorrectLines(t *testing.T) {
	expectedUsage := Usage{
		Total: 0,
		Cores: []CoreUsage{
			{
				ID:    "cpu0",
				Usage: 0,
			},
		},
	}

	files := usageFiles()
	files[statFile] = "cpu 28683 104 7543 6\ncpu0 17496 72 6050 650357 437 0 376 0 0 0"

	reader := usageReader{
		fs:               procfs.New(tempDirWithFiles(t, files, 0o755, 0o600)),
		prevUsageSamples: map[string]usageSample{},
	}

	result, err := reader.read()

	require.NoError(t, err)
	assert.Equal(t, expectedUsage, result)
}

func TestUsageReaderErrorConvertMetrics(t *testing.T) {
	files := usageFiles()
	files[statFile] = "cpu  298830 test 75896 7741115 5083 0 1586 0 0 0"

	reader := usageReader{
		fs:               procfs.New(tempDirWithFiles(t, files, 0o755, 0o600)),
		prevUsageSamples: map[string]usageSample{},
	}

	result, err := reader.read()

	require.Error(t, err)
	assert.Equal(t, Usage{}, result)
}

func TestParseTicksLineSuccess(t *testing.T) {
	expected := ticks{
		id:        "ID",
		user:      1,
		nice:      2,
		system:    3,
		idle:      4,
		iowait:    5,
		irq:       6,
		softirq:   7,
		steal:     8,
		guest:     9,
		guestNice: 0,
	}

	result, err := parseTicksLine([]string{"ID", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0"})

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestParseTicksLineCountError(t *testing.T) {
	result, err := parseTicksLine([]string{"1", "2", "3", "4"})

	assert.ErrorIs(t, err, ErrFewMetricsCountForParsing)
	assert.Equal(t, ticks{}, result)
}

func TestParseTicksLineParsingValueError(t *testing.T) {
	result, err := parseTicksLine([]string{"ID", "1", "2", "3", "test", "1", "2", "3", "test", "1324", "333"})

	assert.ErrorIs(t, err, strconv.ErrSyntax)
	assert.Equal(t, ticks{}, result)
}

func TestCalculateUsageSuccess(t *testing.T) {
	var expected float64 = 50
	curCPUSample := usageSample{
		total:   100,
		idle:    250,
		nonIdle: 350,
	}
	prevCPUSample := usageSample{
		total:   80,
		idle:    200,
		nonIdle: 300,
	}

	result := calculateUsage(prevCPUSample, curCPUSample)

	assert.Equal(t, expected, result)
}

func TestCalculateUsagePrevTickIsBlank(t *testing.T) {
	var expected float64 = 0
	curCPUSample := usageSample{
		total:   100,
		idle:    250,
		nonIdle: 350,
	}
	prevCPUSample := usageSample{}

	result := calculateUsage(prevCPUSample, curCPUSample)

	assert.Equal(t, expected, result)
}

func TestCalculateUsagePrevTickGreaterCurrent(t *testing.T) {
	var expected float64 = 0
	curCPUSample := usageSample{
		total:   80,
		idle:    200,
		nonIdle: 300,
	}
	prevCPUSample := usageSample{
		total:   100,
		idle:    250,
		nonIdle: 350,
	}

	result := calculateUsage(prevCPUSample, curCPUSample)

	assert.Equal(t, expected, result)
}

func TestCalculateUsageNonIdleTotalAndDeltaTotalIsZero(t *testing.T) {
	var expected float64 = 0
	curCPUSample := usageSample{
		total:   80,
		idle:    200,
		nonIdle: 300,
	}
	prevCPUSample := usageSample{
		total:   100,
		idle:    200,
		nonIdle: 300,
	}

	result := calculateUsage(prevCPUSample, curCPUSample)

	assert.Equal(t, expected, result)
}
