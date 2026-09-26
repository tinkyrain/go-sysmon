package metrics

import (
	"testing"

	"github.com/tinkyrain/go-sysmon/internal/procfs"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadCPUSuccessWithoutPrevSample(t *testing.T) {
	expectedCPUUsages := []CPUUsage{
		CPUUsage{ID: "cpu", Usage: 0},
		CPUUsage{ID: "cpu0", Usage: 0},
		CPUUsage{ID: "cpu1", Usage: 0},
	}
	expectedCPUPrevSamples := map[string]CPUSample{
		"cpu": CPUSample{
			Total:   8125909,
			Idle:    7746198,
			NonIdle: 379711,
		},
		"cpu0": {
			Total:   674788,
			Idle:    650794,
			NonIdle: 23994,
		},
		"cpu1": {
			Total:   676907,
			Idle:    661183,
			NonIdle: 15724,
		},
	}
	fs := procfs.New("testdata/")
	cpuReader := CPUReader{
		fs:             fs,
		prevCPUSamples: map[string]CPUSample{},
	}

	readResult, err := cpuReader.Read()

	require.NoError(t, err)
	assert.Equal(t, expectedCPUUsages, readResult)
	assert.Equal(t, expectedCPUPrevSamples, cpuReader.prevCPUSamples)
}

func TestReadCPUSuccessWithPrevSample(t *testing.T) {
	expectedCPUUsages := []CPUUsage{
		CPUUsage{ID: "cpu", Usage: 9.09},
		CPUUsage{ID: "cpu0", Usage: 16.11},
		CPUUsage{ID: "cpu1", Usage: 7.44},
	}
	expectedCPUPrevSamples := map[string]CPUSample{
		"cpu": CPUSample{
			Total:   8125909,
			Idle:    7746198,
			NonIdle: 379711,
		},
		"cpu0": {
			Total:   674788,
			Idle:    650794,
			NonIdle: 23994,
		},
		"cpu1": {
			Total:   676907,
			Idle:    661183,
			NonIdle: 15724,
		},
	}

	fs := procfs.New("testdata/")
	cpuReader := CPUReader{
		fs: fs,
		prevCPUSamples: map[string]CPUSample{
			"cpu": CPUSample{
				Total:   7746198,
				Idle:    6746198,
				NonIdle: 279711,
			},
			"cpu0": CPUSample{
				Total:   650000,
				Idle:    630000,
				NonIdle: 20000,
			},
			"cpu1": CPUSample{
				Total:   600000,
				Idle:    590000,
				NonIdle: 10000,
			},
		},
	}

	readResult, err := cpuReader.Read()

	require.NoError(t, err)
	assert.Equal(t, expectedCPUUsages, readResult)
	assert.Equal(t, expectedCPUPrevSamples, cpuReader.prevCPUSamples)
}

func TestReadCPUWithoutNeedleLines(t *testing.T) {
	expectedCPUUsages := []CPUUsage{}
	tempDir := tempDirWithFile(t, "stat", "", 0o600)
	fs := procfs.New(tempDir)
	cpuReader := CPUReader{
		fs:             fs,
		prevCPUSamples: map[string]CPUSample{},
	}

	readResult, err := cpuReader.Read()

	assert.ErrorIs(t, err, ErrNoCPULines)
	assert.Equal(t, expectedCPUUsages, readResult)
}

func TestReadCPUFileNotFound(t *testing.T) {
	expectedCPUUsages := []CPUUsage{}
	tempDir := tempDirWithFile(t, "stat312", "", 0o600)
	fs := procfs.New(tempDir)
	cpuReader := CPUReader{
		fs:             fs,
		prevCPUSamples: map[string]CPUSample{},
	}

	readResult, err := cpuReader.Read()

	assert.Error(t, err)
	assert.Equal(t, expectedCPUUsages, readResult)
}

func TestReadCPUFilterIncorrectLines(t *testing.T) {
	expectedCPUUsages := []CPUUsage{
		{
			ID:    "cpu0",
			Usage: 0,
		},
	}
	tempDir := tempDirWithFile(t,
		"stat",
		"cpu 28683 104 7543 6\ncpu0 17496 72 6050 650357 437 0 376 0 0 0",
		0o600)
	fs := procfs.New(tempDir)
	cpuReader := CPUReader{
		fs:             fs,
		prevCPUSamples: map[string]CPUSample{},
	}

	readResult, err := cpuReader.Read()

	require.NoError(t, err)
	assert.Equal(t, expectedCPUUsages, readResult)
}

func TestReadCPUErrorConvertMetrics(t *testing.T) {
	expectedCPUUsages := []CPUUsage{}
	tempDir := tempDirWithFile(t,
		"stat",
		"cpu  298830 test 75896 7741115 5083 0 1586 0 0 0",
		0o600)
	fs := procfs.New(tempDir)
	cpuReader := CPUReader{
		fs:             fs,
		prevCPUSamples: map[string]CPUSample{},
	}

	readResult, err := cpuReader.Read()

	require.Error(t, err)
	assert.Equal(t, expectedCPUUsages, readResult)
}

func TestParseCPULineSuccess(t *testing.T) {
	metrics := []string{"ID", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}
	expected := CPUTicks{
		ID:        "ID",
		User:      1,
		Nice:      2,
		System:    3,
		Idle:      4,
		Iowait:    5,
		Irq:       6,
		Softirq:   7,
		Steal:     8,
		Guest:     9,
		GuestNice: 0,
	}

	parseResult, err := parseCPULine(metrics)

	require.NoError(t, err)
	assert.Equal(t, expected, parseResult)
}

func TestParseCPULineCountError(t *testing.T) {
	metrics := []string{"1", "2", "3", "4"}
	expected := CPUTicks{}

	parseResult, err := parseCPULine(metrics)

	assert.ErrorIs(t, err, ErrFewMetricsCountForParsing)
	assert.Equal(t, expected, parseResult)
}

func TestParseCPULineParsingValueError(t *testing.T) {
	metrics := []string{"ID", "1", "2", "3", "test", "1", "2", "3", "test", "1324", "333"}
	expected := CPUTicks{}

	parseResult, err := parseCPULine(metrics)

	assert.ErrorContains(t, err, "parsing cpu line")
	assert.Equal(t, expected, parseResult)
}

func TestCalculateCPUUsageSuccess(t *testing.T) {
	var expected float64 = 50
	curCPUTick := CPUSample{
		Total:   100,
		Idle:    250,
		NonIdle: 350,
	}
	prevCPUTick := CPUSample{
		Total:   80,
		Idle:    200,
		NonIdle: 300,
	}

	usage := calculateCPUUsage(prevCPUTick, curCPUTick)

	assert.Equal(t, expected, usage)
}

func TestCalculateCPUUsagePrevTickIsBlank(t *testing.T) {
	var expected float64 = 0
	curCPUTick := CPUSample{
		Total:   100,
		Idle:    250,
		NonIdle: 350,
	}
	prevCPCUTick := CPUSample{}

	usage := calculateCPUUsage(prevCPCUTick, curCPUTick)

	assert.Equal(t, expected, usage)
}

func TestCalculateCPUUsagePrevTickGreaterCurrent(t *testing.T) {
	var expected float64 = 0
	curCPUTick := CPUSample{
		Total:   80,
		Idle:    200,
		NonIdle: 300,
	}
	prevCPCUTick := CPUSample{
		Total:   100,
		Idle:    250,
		NonIdle: 350,
	}

	usage := calculateCPUUsage(prevCPCUTick, curCPUTick)

	assert.Equal(t, expected, usage)
}

func TestCalculateCPUUsageNonIdleTotalAndDeltaTotalIsZero(t *testing.T) {
	var expected float64 = 0
	curCPUTick := CPUSample{
		Total:   80,
		Idle:    200,
		NonIdle: 300,
	}
	prevCPCUTick := CPUSample{
		Total:   100,
		Idle:    200,
		NonIdle: 300,
	}

	usage := calculateCPUUsage(prevCPCUTick, curCPUTick)

	assert.Equal(t, expected, usage)
}
