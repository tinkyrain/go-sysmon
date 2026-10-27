package components

import (
	tui "github.com/metaspartan/gotui/v5"
	"github.com/metaspartan/gotui/v5/widgets"
)

func Hints(x1, y1, x2, y2 int) *tui.Grid {

	quit := widgets.NewParagraph()
	quit.Title = "Quit"
	quit.Text = "Ctrl+C / Q"

	quit.TextStyle.Fg = tui.ColorBlack
	quit.TextStyle.Bg = tui.ColorWhite
	quit.BorderStyle.Fg = tui.ColorBlack
	quit.BorderStyle.Bg = tui.ColorWhite
	quit.TitleStyle.Fg = tui.ColorBlack
	quit.TitleStyle.Bg = tui.ColorWhite
	quit.BackgroundColor = tui.ColorWhite

	quit.SetRect(0, 0, 1, 1)

	tab := widgets.NewParagraph()
	tab.Title = "Tab switch"
	tab.Text = "Arrows / Tab"

	tab.TextStyle.Fg = tui.ColorBlack
	tab.TextStyle.Bg = tui.ColorWhite
	tab.BorderStyle.Fg = tui.ColorBlack
	tab.BorderStyle.Bg = tui.ColorWhite
	tab.TitleStyle.Fg = tui.ColorBlack
	tab.TitleStyle.Bg = tui.ColorWhite
	tab.BackgroundColor = tui.ColorWhite

	tab.SetRect(1, 0, 1, 1)

	grid := tui.NewGrid()
	grid.SetRect(x1, y1, x2, y2)

	spacer := widgets.NewParagraph()
	spacer.Border = false

	grid.Set(
		tui.NewRow(
			1.0/2,
			tui.NewCol(0.4995, quit),
			tui.NewCol(0.001, spacer),
			tui.NewCol(0.4995, tab),
		),
	)

	return grid
}
