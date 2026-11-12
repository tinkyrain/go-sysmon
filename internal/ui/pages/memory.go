package pages

import (
	tui "github.com/metaspartan/gotui/v5"
	"github.com/metaspartan/gotui/v5/widgets"
	"github.com/tinkyrain/go-sysmon/internal/ui/components"
	"github.com/tinkyrain/go-sysmon/internal/ui/view"
)

func MemoryPage(data view.Memory, x1, y1, x2, y2 int) tui.Drawable {
	info := widgets.NewParagraph()
	info.Title = "Info"
	info.Text = fieldsText(data.Info)

	memoryUsageLine := components.UsageLine(data.RAM.Label, data.RAM.Percent)
	swapUsageLine := components.UsageLine(data.Swap.Label, data.Swap.Percent)

	grid := tui.NewGrid()
	grid.Set(
		tui.NewRow(1,
			tui.NewCol(0.5, info),
			tui.NewCol(0.5,
				tui.NewRow(0.5, memoryUsageLine),
				tui.NewRow(0.5, swapUsageLine),
			),
		),
	)
	grid.SetRect(x1, y1, x2, y2)

	return grid
}
