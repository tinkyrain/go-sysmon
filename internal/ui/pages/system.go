package pages

import (
	"fmt"

	tui "github.com/metaspartan/gotui/v5"
	"github.com/metaspartan/gotui/v5/widgets"
	"github.com/tinkyrain/go-sysmon/internal/ui/view"
)

func SystemPage(data view.System, x1, y1, x2, y2 int) tui.Drawable {
	info := widgets.NewParagraph()
	info.Title = "Info"
	info.Text = fieldsText(data.Info)

	loadAvg := widgets.NewParagraph()
	loadAvg.Title = "LoadAvg"
	loadAvg.Text = fieldsText(data.LoadAvg)

	uptime := widgets.NewParagraph()
	uptime.Title = "Uptime"
	uptime.TitleAlignment = tui.AlignCenter
	uptime.TextAlignment = tui.AlignCenter
	uptime.VerticalAlignment = tui.AlignMiddle
	uptime.Text = fmt.Sprintf("Total: %s    Idle: %s", data.Uptime.Total, data.Uptime.Idle)

	grid := tui.NewGrid()
	grid.SetRect(x1, y1, x2, y2)
	grid.Set(
		tui.NewRow(
			0.85,
			tui.NewCol(0.5, info),
			tui.NewCol(0.5, loadAvg),
		),
		tui.NewRow(
			0.15,
			tui.NewCol(1, uptime),
		),
	)

	return grid
}
