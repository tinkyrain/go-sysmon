package pages

import (
	"sort"

	tui "github.com/metaspartan/gotui/v5"
	"github.com/metaspartan/gotui/v5/widgets"
	"github.com/tinkyrain/go-sysmon/internal/metrics/cpu"
	"github.com/tinkyrain/go-sysmon/internal/ui/components"
)

const usageLineHeight = 3

func CPUPage(stats cpu.Stats, x1, y1, x2, y2 int) tui.Drawable {
	column := widgets.NewFlex()
	column.Border = false
	column.Direction = widgets.FlexColumn
	column.SetRect(x1-1, y1-1, x2+1, y2+1)

	if y2-y1 < usageLineHeight {
		return column
	}

	column.AddItem(components.UsageLine("Total", int(stats.Usage.Total)), usageLineHeight, 0, false)

	coresTitle := widgets.NewParagraph()
	coresTitle.Border = false
	coresTitle.TextAlignment = tui.AlignCenter
	coresTitle.Text = "Cores"
	column.AddItem(coresTitle, 3, 1, false)

	coreCount := (y2-y1)/usageLineHeight - 1
	for _, core := range topCores(stats.Usage.Cores, coreCount) {
		column.AddItem(components.UsageLine(core.ID, int(core.Usage)), usageLineHeight, 0, false)
	}

	return column
}

func topCores(all []cpu.CoreUsage, count int) []cpu.CoreUsage {
	if count < 1 {
		return nil
	}

	cores := make([]cpu.CoreUsage, len(all))
	copy(cores, all)

	sort.SliceStable(cores, func(i, j int) bool {
		if cores[i].Usage != cores[j].Usage {
			return cores[i].Usage > cores[j].Usage
		}
		return cores[i].ID < cores[j].ID
	})

	if len(cores) > count {
		cores = cores[:count]
	}

	return cores
}
