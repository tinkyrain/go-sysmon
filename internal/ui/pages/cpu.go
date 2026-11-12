package pages

import (
	tui "github.com/metaspartan/gotui/v5"
	"github.com/metaspartan/gotui/v5/widgets"
	"github.com/tinkyrain/go-sysmon/internal/ui/components"
	"github.com/tinkyrain/go-sysmon/internal/ui/view"
)

const (
	coreHeight       = 3
	coresTitleHeight = 3
)

func CPUPage(data view.CPU, x1, y1, x2, y2 int) tui.Drawable {
	column := widgets.NewFlex()
	column.Border = false
	column.Direction = widgets.FlexColumn
	column.SetRect(x1-1, y1-1, x2+1, y2+1)

	if y2-y1 < coreHeight {
		return column
	}

	column.AddItem(components.UsageLine(data.Total.Label, data.Total.Percent), coreHeight, 0, false)

	coresTitle := widgets.NewParagraph()
	coresTitle.Border = false
	coresTitle.TextAlignment = tui.AlignCenter
	coresTitle.Text = "Cores"
	column.AddItem(coresTitle, coresTitleHeight, 0, false)

	headerHeight := coreHeight + coresTitleHeight
	coreCount := (y2 - y1 - headerHeight) / coreHeight
	for _, core := range visibleCores(data.Cores, coreCount) {
		column.AddItem(components.UsageLine(core.Label, core.Percent), coreHeight, 0, false)
	}

	return column
}

func visibleCores(cores []view.Gauge, count int) []view.Gauge {
	if count < 1 {
		return nil
	}

	if len(cores) > count {
		return cores[:count]
	}

	return cores
}
