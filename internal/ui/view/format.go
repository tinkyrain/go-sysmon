package view

import (
	"fmt"
	"math"
	"strconv"
)

const timeLayout = "02.01.2006 15:04"

const unknown = "n/a"

var sizeUnits = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}

func HumanizeSize(size float64) string {
	if math.IsNaN(size) || math.IsInf(size, 0) || size < 0 {
		return unknown
	}

	i := 0
	for size >= 1024 && i < len(sizeUnits)-1 {
		size /= 1024
		i++
	}

	return strconv.FormatFloat(size, 'f', 2, 64) + " " + sizeUnits[i]
}

func HumanizeDuration(seconds float64) string {
	if math.IsNaN(seconds) || math.IsInf(seconds, 0) || seconds < 0 {
		return unknown
	}

	total := int64(seconds)
	days := total / 86400
	hours := total % 86400 / 3600
	minutes := total % 3600 / 60

	switch {
	case days > 0:
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	case hours > 0:
		return fmt.Sprintf("%dh %dm", hours, minutes)
	case minutes > 0:
		return fmt.Sprintf("%dm %ds", minutes, total%60)
	default:
		return fmt.Sprintf("%ds", total%60)
	}
}

func usedPercent(total, available uint64) int {
	if total == 0 || available >= total {
		return 0
	}

	return int(float64(total-available) / float64(total) * 100)
}
