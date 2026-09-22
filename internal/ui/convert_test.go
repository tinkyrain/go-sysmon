package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHumanizeSizeSuccess(t *testing.T) {
	expected := "1.07 GB"
	result := humanizeSize(float64(1073741824))
	assert.Equal(t, expected, result)
}

func TestHumanizeSizeInputSizeLess1000(t *testing.T) {
	expected := "1.00 GB"
	result := humanizeSize(float64(1000000000))
	assert.Equal(t, expected, result)
}

func TestHumanizeSizeInputSizeEqual1000(t *testing.T) {
	expected := "1.00 KB"
	result := humanizeSize(float64(1000))
	assert.Equal(t, expected, result)
}

func TestUserPercentSuccess(t *testing.T) {
	var expected float64 = 30
	result := usedPercent(100, 70)
	assert.Equal(t, expected, result)
}

func TestUserPercentTotalIsZero(t *testing.T) {
	var expected float64 = 0
	result := usedPercent(0, 70)
	assert.Equal(t, expected, result)
}

func TestUserPercentAvailableIsZero(t *testing.T) {
	var expected float64 = 100
	result := usedPercent(100, 0)
	assert.Equal(t, expected, result)
}
