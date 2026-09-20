package ui

import (
	"strconv"
)

var sizeUnits = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}

func humanizeSize(size float64) string {
	i := 0
	for size >= 1000 && i < len(sizeUnits)-1 {
		size /= 1000
		i++
	}
	return strconv.FormatFloat(size, 'f', 2, 64) + " " + sizeUnits[i]
}
