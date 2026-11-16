package view

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHumanizeDuration(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		seconds  float64
		expected string
	}{
		{
			name:     "zero",
			seconds:  0,
			expected: "0s",
		},
		{
			name:     "seconds only",
			seconds:  42.9,
			expected: "42s",
		},
		{
			name:     "minutes and seconds",
			seconds:  125,
			expected: "2m 5s",
		},
		{
			name:     "hours drop the seconds",
			seconds:  3*3600 + 4*60 + 59,
			expected: "3h 4m",
		},
		{
			name:     "days",
			seconds:  5*86400 + 21*3600 + 25*60,
			expected: "5d 21h 25m",
		},
		{
			name:     "day boundary keeps zeroes",
			seconds:  86400,
			expected: "1d 0h 0m",
		},
		{
			name:     "negative is not a duration",
			seconds:  -1,
			expected: "n/a",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, HumanizeDuration(tc.seconds))
		})
	}
}

func TestHumanizeSize(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		size     float64
		expected string
	}{
		{
			name:     "zero",
			size:     0,
			expected: "0.00 B",
		},
		{
			name:     "bytes stay bytes",
			size:     1023,
			expected: "1023.00 B",
		},
		{
			name:     "exact kilobyte",
			size:     kib,
			expected: "1.00 KB",
		},
		{
			name:     "fraction is kept",
			size:     1536,
			expected: "1.50 KB",
		},
		{
			name:     "megabytes",
			size:     512 * mib,
			expected: "512.00 MB",
		},
		{
			name:     "gigabytes match free -h",
			size:     16 * gib,
			expected: "16.00 GB",
		},
		{
			name:     "fraction below the next unit",
			size:     2*gib - 1,
			expected: "2.00 GB",
		},
		{
			name:     "exabytes",
			size:     1 << 60,
			expected: "1.00 EB",
		},
		{
			name:     "beyond the largest unit it stops dividing",
			size:     1 << 70,
			expected: "1024.00 EB",
		},
		{
			name:     "negative is not a size",
			size:     -1,
			expected: "n/a",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, HumanizeSize(tc.size))
		})
	}
}

func TestUsedPercent(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		total     uint64
		available uint64
		expected  int
	}{
		{
			name:      "an unknown total reads as zero",
			total:     0,
			available: 0,
			expected:  0,
		},
		{
			name:      "available without a total reads as zero",
			total:     0,
			available: 50,
			expected:  0,
		},
		{
			name:      "nothing used reads as zero",
			total:     50,
			available: 50,
			expected:  0,
		},
		{
			name:      "used is the share that is not available",
			total:     120,
			available: 90,
			expected:  25,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, usedPercent(tc.total, tc.available))
		})
	}
}
