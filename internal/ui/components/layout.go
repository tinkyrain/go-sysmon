package components

import (
	tui "github.com/metaspartan/gotui/v5"
	"github.com/metaspartan/gotui/v5/widgets"
)

func Layout(x1, y1, x2, y2 int) *widgets.Paragraph {
	layout := widgets.NewParagraph()
	layout.Title = "go-sysmon"
	layout.Text = "go-sysmon"
	layout.SetRect(x1, y1, x2, y2)
	layout.Border = true
	layout.TextStyle = tui.NewStyle(tui.ColorWhite, tui.ColorClear, tui.ModifierBold)
	return layout
}
