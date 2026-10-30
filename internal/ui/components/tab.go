package components

import (
	"unicode/utf8"

	tui "github.com/metaspartan/gotui/v5"
	"github.com/metaspartan/gotui/v5/widgets"
)

// Tabs - render tab row.
// x1,x2 - allocated width coordinate
// y - allocated row coordinate
func Tab(x1, x2, y int) *widgets.TabPane {
	options := []string{"System", "CPU", "Memory", "Disk"}
	gap, pl, pr := 1, 1, 1
	tabWidth := 0

	for key, option := range options {
		if key != 0 {
			tabWidth += gap
		}
		tabWidth += pl
		tabWidth += pr
		tabWidth += utf8.RuneCountInString(option)
	}
	start := x1 + (x2-x1-tabWidth)/2

	tabpane := widgets.NewTabPane(options...)
	tabpane.Border = false
	tabpane.ActiveTabStyle, tabpane.InactiveTabStyle =
		tui.NewStyle(tui.ColorBlack, tui.ColorWhite),
		tui.NewStyle(tui.ColorClear)
	tabpane.PadRight, tabpane.PadLeft, tabpane.TabGap = pr, pl, gap
	tabpane.Separator = ""
	tabpane.SetRect(start-1, y, start+tabWidth+1, y)

	return tabpane
}
