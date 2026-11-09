package components

import (
	tui "github.com/metaspartan/gotui/v5"
)

func Layout(timeText string, x1, y1, x2, y2 int) *tui.Block {
	layout := tui.NewBlock()
	layout.BorderStyle = tui.NewStyle(tui.ColorWhite)
	layout.TitleLeft = "go-sysmon"
	layout.TitleRight = timeText
	layout.TitleStyle = tui.NewStyle(tui.ColorWhite, tui.ColorClear, tui.ModifierBold)
	layout.SetRect(x1, y1, x2, y2)
	return layout
}
