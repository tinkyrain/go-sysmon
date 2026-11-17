package components

import (
	"github.com/metaspartan/gotui/v5/widgets"
)

func UsageBlock(title string, percent int) *widgets.Gauge {
	gauge := widgets.NewGauge()
	gauge.Title = title
	gauge.Percent = percent
	return gauge
}
