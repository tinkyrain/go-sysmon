package pages

import (
	"sort"

	tui "github.com/metaspartan/gotui/v5"
	"github.com/metaspartan/gotui/v5/widgets"
	"github.com/tinkyrain/go-sysmon/internal/metrics/cpu"
	"github.com/tinkyrain/go-sysmon/internal/ui/components"
)

func CPUPage(stats cpu.Stats, x1, y1, x2, y2 int) tui.Drawable {
	usageRowSize := 3

	totalUsage := components.UsageLine("Total", int(stats.Usage.Total))

	// sorted cpu cores by usages
	cores := make([]cpu.CoreUsage, len(stats.Usage.Cores))
	copy(cores, stats.Usage.Cores)
	sort.Slice(cores, func(i, j int) bool {
		return cores[i].Usage > cores[j].Usage
	})

	coresColumn := widgets.NewFlex()
	coresColumn.Border = false
	coresColumn.Direction = widgets.FlexColumn
	for _, core := range cores {
		coresColumn.AddItem(components.UsageLine(core.ID, int(core.Usage)), usageRowSize, 0, false)
	}

	emptyRight := widgets.NewFlex()
	emptyRight.Border = false
	emptyRight.Direction = widgets.FlexColumn

	screenSplit := widgets.NewFlex()
	screenSplit.Border = false
	screenSplit.Direction = widgets.FlexRow

	screenSplit.AddItem(coresColumn, 0, 1, false)
	screenSplit.AddItem(emptyRight, 0, 1, false)

	mainFlex := widgets.NewFlex()
	mainFlex.Border = false
	mainFlex.Direction = widgets.FlexColumn
	mainFlex.AddItem(totalUsage, usageRowSize, 0, false)
	mainFlex.AddItem(screenSplit, y2-y1-usageRowSize, 1, false)

	mainFlex.SetRect(x1, y1, x2, y2)

	return mainFlex
}
