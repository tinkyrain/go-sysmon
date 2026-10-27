package components

import (
	tui "github.com/metaspartan/gotui/v5"
)

func Layout(x1, y1, x2, y2 int) *tui.Block {
	layout := &tui.Block{
		Border:       true,
		BorderStyle:  tui.NewStyle(tui.ColorWhite),
		BorderLeft:   true,
		BorderRight:  true,
		BorderTop:    true,
		BorderBottom: true,
		BorderType:   tui.BorderLine,
		Title:        "go-sysmon",
		TitleStyle:   tui.NewStyle(tui.ColorWhite, tui.ColorClear, tui.ModifierBold),
	}
	layout.SetRect(x1, y1, x2, y2)
	return layout
}
