package components

import (
	tui "github.com/metaspartan/gotui/v5"
	"github.com/metaspartan/gotui/v5/widgets"
)

func UsageLine(title string, percent int) *widgets.LineGauge {
	line := widgets.NewLineGauge()
	line.Title = title
	line.Percent = percent
	line.BarRune = '█'
	line.BarRuneEmpty = '░'
	line.LineColor = tui.ColorWhite
	return line
}
