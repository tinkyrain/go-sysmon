package metrics

import (
	"go-sysmon/internal/procfs"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadCPUSuccessWithoutPrevTick(t *testing.T) {
	var expectedCPUUsage float64 = 0
	expectedCPUPrevTick := CPUTick{
		Total:   8125909,
		Idle:    7746198,
		NonIdle: 379711,
	}
	fs := procfs.New("testdata/")
	cpuReader := CPUReader{
		fs:          fs,
		prevCPUTick: CPUTick{},
	}

	readResult, err := cpuReader.Read()

	require.NoError(t, err)
	assert.Equal(t, readResult, expectedCPUUsage)
	assert.Equal(t, expectedCPUPrevTick, cpuReader.prevCPUTick)
}

func TestReadCPUSuccessWithPrevTick(t *testing.T) {
	expectedCPUUsage := 9.09
	expectedCPUPrevTick := CPUTick{
		Total:   8125909,
		Idle:    7746198,
		NonIdle: 379711,
	}
	fs := procfs.New("testdata/")
	cpuReader := CPUReader{
		fs: fs,
		prevCPUTick: CPUTick{
			Total:   7746198,
			Idle:    6746198,
			NonIdle: 279711,
		},
	}

	readResult, err := cpuReader.Read()

	require.NoError(t, err)
	assert.Equal(t, readResult, expectedCPUUsage)
	assert.Equal(t, expectedCPUPrevTick, cpuReader.prevCPUTick)
}

func TestReadCPUBlankFile(t *testing.T) {
	var expected float64 = 0
	tempDir := tempDirWithFile(t, "stat", "", 0o600)
	fs := procfs.New(tempDir)
	cpuReader := CPUReader{
		fs:          fs,
		prevCPUTick: CPUTick{},
	}

	readResult, err := cpuReader.Read()

	assert.NoError(t, err)
	assert.Equal(t, expected, readResult)
}

func TestReadCPUNeedleMetricsNotFound(t *testing.T) {
	var expectedCPUUsage float64 = 0
	expectedCPUPrevTick := CPUTick{
		Total:   0,
		Idle:    0,
		NonIdle: 0,
	}
	tempDir := tempDirWithFile(t,
		"stat",
		"cpu10 28683 104 7543 640187 301 0 10 0 0 0\ncpu11 17748 68 4300 655703 656 0 6 0 0 0",
		0o600)
	fs := procfs.New(tempDir)
	cpuReader := CPUReader{
		fs:          fs,
		prevCPUTick: CPUTick{},
	}

	readResult, err := cpuReader.Read()

	require.NoError(t, err)
	assert.Equal(t, readResult, expectedCPUUsage)
	assert.Equal(t, expectedCPUPrevTick, cpuReader.prevCPUTick)
}

func TestReadCPUFileNotFound(t *testing.T) {
	var expectedCPUUsage float64 = 0
	tempDir := tempDirWithFile(t, "stat312", "", 0o600)
	fs := procfs.New(tempDir)
	cpuReader := CPUReader{
		fs:          fs,
		prevCPUTick: CPUTick{},
	}

	readResult, err := cpuReader.Read()

	assert.Error(t, err)
	assert.Equal(t, readResult, expectedCPUUsage)
}

func TestReadCPUIncorrectMetricLine(t *testing.T) {
	var expectedCPUUsage float64 = 0
	expectedCPUPrevTick := CPUTick{
		Total:   0,
		Idle:    0,
		NonIdle: 0,
	}
	tempDir := tempDirWithFile(t,
		"stat",
		"cpu 28683 104 7543 6",
		0o600)
	fs := procfs.New(tempDir)
	cpuReader := CPUReader{
		fs:          fs,
		prevCPUTick: CPUTick{},
	}

	readResult, err := cpuReader.Read()

	require.NoError(t, err)
	assert.Equal(t, readResult, expectedCPUUsage)
	assert.Equal(t, expectedCPUPrevTick, cpuReader.prevCPUTick)
}

func TestReadCPUErrorConvertMetrics(t *testing.T) {
	var expectedCPUUsage float64 = 0
	tempDir := tempDirWithFile(t,
		"stat",
		"cpu  298830 test 75896 7741115 5083 0 1586 0 0 0",
		0o600)
	fs := procfs.New(tempDir)
	cpuReader := CPUReader{
		fs:          fs,
		prevCPUTick: CPUTick{},
	}

	readResult, err := cpuReader.Read()

	require.Error(t, err)
	assert.Equal(t, readResult, expectedCPUUsage)
}

func TestParseCPUTickSuccess(t *testing.T) {
	metrics := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}
	expected := map[string]uint64{
		"user":      1,
		"nice":      2,
		"system":    3,
		"idle":      4,
		"iowait":    5,
		"irq":       6,
		"softirq":   7,
		"steal":     8,
		"guest":     9,
		"guestNice": 0,
	}

	parseResult, err := parseCPUTick(metrics)

	require.NoError(t, err)
	assert.InDeltaMapValues(t, expected, parseResult, 0)
}

func TestParseCPUTickMetricsCountError(t *testing.T) {
	metrics := []string{"1", "2", "3", "4"}
	expected := map[string]uint64{}

	parseResult, err := parseCPUTick(metrics)

	assert.ErrorIs(t, err, ErrFewMetricsCountForParsing)
	assert.Equal(t, expected, parseResult)
}

func TestParseCPUTickParsingValueError(t *testing.T) {
	metrics := []string{"1", "2", "3", "test", "1", "2", "3", "test", "1324", "333"}
	expected := map[string]uint64{}

	parseResult, err := parseCPUTick(metrics)

	assert.ErrorContains(t, err, "parsing cpu metric")
	assert.Equal(t, expected, parseResult)
}

func TestCalculateCPUUsageSuccess(t *testing.T) {
	var expected float64 = 50
	curCPUTick := CPUTick{
		Total:   100,
		Idle:    250,
		NonIdle: 350,
	}
	prevCPCUTick := CPUTick{
		Total:   80,
		Idle:    200,
		NonIdle: 300,
	}

	usage := calculateCPUUsage(prevCPCUTick, curCPUTick)

	assert.Equal(t, expected, usage)
}

func TestCalculateCPUUsagePrevTickIsBlank(t *testing.T) {
	var expected float64 = 0
	curCPUTick := CPUTick{
		Total:   100,
		Idle:    250,
		NonIdle: 350,
	}
	prevCPCUTick := CPUTick{}

	usage := calculateCPUUsage(prevCPCUTick, curCPUTick)

	assert.Equal(t, expected, usage)
}

func TestCalculateCPUUsagePrevTickGreaterCurrent(t *testing.T) {
	var expected float64 = 0
	curCPUTick := CPUTick{
		Total:   80,
		Idle:    200,
		NonIdle: 300,
	}
	prevCPCUTick := CPUTick{
		Total:   100,
		Idle:    250,
		NonIdle: 350,
	}

	usage := calculateCPUUsage(prevCPCUTick, curCPUTick)

	assert.Equal(t, expected, usage)
}

func TestCalculateCPUUsageNonIdleTotalAndDeltaTotalIsZero(t *testing.T) {
	var expected float64 = 0
	curCPUTick := CPUTick{
		Total:   80,
		Idle:    200,
		NonIdle: 300,
	}
	prevCPCUTick := CPUTick{
		Total:   100,
		Idle:    200,
		NonIdle: 300,
	}

	usage := calculateCPUUsage(prevCPCUTick, curCPUTick)

	assert.Equal(t, expected, usage)
}
