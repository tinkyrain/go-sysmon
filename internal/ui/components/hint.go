package components

import (
	tui "github.com/metaspartan/gotui/v5"
	"github.com/metaspartan/gotui/v5/widgets"
)

func Hint(x1, y1, x2, y2 int) *tui.Grid {
	options := []struct {
		title string
		text  string
	}{
		{
			title: "Quit",
			text:  "Ctrl+C / Q",
		},
		{
			title: "Tab switch",
			text:  "Arrows / Tab",
		},
	}

	spacer := widgets.NewParagraph()
	spacer.Border = false
	spacerCount := len(options) - 1
	spacerRatio := 0.01

	hintRatio := (1.0 - (spacerRatio * float64(spacerCount))) / float64(len(options))

	var hintCols []any

	for key, option := range options {
		hint := widgets.NewParagraph()
		hint.Title = option.title
		hint.Text = option.text
		hint.TextStyle.Fg = tui.ColorBlack
		hint.TextStyle.Bg = tui.ColorWhite
		hint.BorderStyle.Fg = tui.ColorBlack
		hint.BorderStyle.Bg = tui.ColorWhite
		hint.TitleStyle.Fg = tui.ColorBlack
		hint.TitleStyle.Bg = tui.ColorWhite
		hint.BackgroundColor = tui.ColorWhite

		hintCols = append(hintCols, tui.NewCol(hintRatio, hint))

		if key+1 < len(options) {
			hintCols = append(hintCols, tui.NewCol(spacerRatio, spacer))
		}
	}

	grid := tui.NewGrid()
	grid.SetRect(x1, y1, x2, y2)

	grid.Set(
		tui.NewRow(
			1,
			hintCols...,
		),
	)

	return grid
}
