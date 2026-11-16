package components

import (
	"github.com/metaspartan/gotui/v5/widgets"
)

func UsageBlock(title string, percent int) *widgets.Gauge {
	line := widgets.NewGauge()
	line.Title = title
	line.Percent = percent
	return line
}
