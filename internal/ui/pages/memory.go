package pages

import (
	"fmt"
	"strconv"
	"strings"

	tui "github.com/metaspartan/gotui/v5"
	"github.com/metaspartan/gotui/v5/widgets"
	"github.com/tinkyrain/go-sysmon/internal/metrics/memory"
	"github.com/tinkyrain/go-sysmon/internal/ui/components"
)

func MemoryPage(stats memory.Stats, x1, y1, x2, y2 int) tui.Drawable {
	memoryPercent := int((1.0 - (float64(stats.Memory.Available) / float64(stats.Memory.Total))) * 100)

	var swapPercent int
	if stats.Memory.SwapTotal > 0 {
		swapPercent = int((1.0 - (float64(stats.Memory.SwapAvailable) / float64(stats.Memory.SwapTotal))) * 100)
	}

	info := widgets.NewParagraph()
	info.Title = "Info"
	info.Text = strings.Join([]string{
		fmt.Sprintf("Total: %s", HumanizeSize(float64(stats.Memory.Total))),
		fmt.Sprintf("Available: %s", HumanizeSize(float64(stats.Memory.Available))),
		fmt.Sprintf("SWAP Total: %s", HumanizeSize(float64(stats.Memory.SwapTotal))),
		fmt.Sprintf("SWAP Available: %s", HumanizeSize(float64(stats.Memory.SwapAvailable))),
	}, "\n")

	memoryUsageLine := components.UsageLine("Memory", memoryPercent)
	swapUsageLine := components.UsageLine("SWAP", swapPercent)

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

var sizeUnits = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}

func HumanizeSize(size float64) string {
	i := 0
	for size >= 1024 && i < len(sizeUnits)-1 {
		size /= 1024
		i++
	}
	return strconv.FormatFloat(size, 'f', 2, 64) + " " + sizeUnits[i]
}
