package pages

import (
	"sort"

	tui "github.com/metaspartan/gotui/v5"
	"github.com/tinkyrain/go-sysmon/internal/metrics/cpu"
	"github.com/tinkyrain/go-sysmon/internal/ui/components"
)

func CPUPage(stats cpu.Stats, x1, y1, x2, y2 int) tui.Drawable {
	totalUsage := components.UsageLine("Total", int(stats.Usage.Total))
	coreLines := []any{}

	// get top 7 usage cores for display
	cores := make([]cpu.CoreUsage, len(stats.Usage.Cores))
	copy(cores, stats.Usage.Cores)
	sort.Slice(cores, func(i, j int) bool {
		return cores[i].Usage > cores[j].Usage
	})
	if len(cores) > 7 {
		cores = cores[:7]
	}

	for _, core := range cores {
		rowHeight := 1.0 / float64(len(cores))
		coreLines = append(
			coreLines,
			tui.NewRow(rowHeight, components.UsageLine(core.ID, int(core.Usage))),
		)
	}

	grid := tui.NewGrid()
	grid.SetRect(x1, y1, x2, y2)
	grid.Set(
		tui.NewRow(0.1, totalUsage),
		tui.NewRow(0.9,
			tui.NewCol(0.5, coreLines...),
		),
	)

	return grid
}
