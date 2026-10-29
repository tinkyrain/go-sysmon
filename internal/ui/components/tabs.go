package components

import (
	tui "github.com/metaspartan/gotui/v5"
	"github.com/metaspartan/gotui/v5/widgets"
)

func Tabs(x1, y1, x2, y2 int) *widgets.TabPane {
	tabpane := widgets.NewTabPane("System", "CPU", "Memory", "Disk")
	tabpane.SetRect(x1, y1, x2, y2)

	tabpane.Border = false
	tabpane.ActiveTabStyle, tabpane.InactiveTabStyle =
		tui.NewStyle(tui.ColorBlack, tui.ColorWhite),
		tui.NewStyle(tui.ColorClear)
	tabpane.PadRight, tabpane.PadLeft, tabpane.TabGap = 1, 1, 1
	tabpane.Separator = ""

	return tabpane
}
