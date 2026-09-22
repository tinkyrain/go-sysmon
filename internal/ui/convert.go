package ui

import (
	"strconv"
)

var sizeUnits = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}

// humanizeSize convert bytes to maximum memory value
func humanizeSize(size float64) string {
	i := 0
	for size >= 1000 && i < len(sizeUnits)-1 {
		size /= 1000
		i++
	}
	return strconv.FormatFloat(size, 'f', 2, 64) + " " + sizeUnits[i]
}

// usedPercent calculate percent value
func usedPercent(total, available uint64) float64 {
	if total == 0 {
		return 0
	}
	return float64(total-available) / float64(total) * 100
}
