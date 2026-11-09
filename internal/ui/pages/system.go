package pages

import (
	"fmt"
	"strings"

	tui "github.com/metaspartan/gotui/v5"
	"github.com/metaspartan/gotui/v5/widgets"
	"github.com/tinkyrain/go-sysmon/internal/metrics/system"
)

func SystemPage(stats system.Stats, x1, y1, x2, y2 int) tui.Drawable {
	info := widgets.NewParagraph()
	info.Title = "Info"
	info.Text = strings.Join([]string{
		fmt.Sprintf("Hostname: %s", stats.Info.Hostname),
		fmt.Sprintf("OS: %s", stats.Info.OS),
		fmt.Sprintf("Kernel: %s", stats.Info.Kernel),
		fmt.Sprintf("Arch: %s", stats.Info.Arch),
	}, "\n")

	loadAvg := widgets.NewParagraph()
	loadAvg.Title = "LoadAvg"
	loadAvg.Text = strings.Join([]string{
		fmt.Sprintf("One min: %.2f", stats.LoadAvg.OneMin),
		fmt.Sprintf("Five min: %.2f", stats.LoadAvg.FiveMin),
		fmt.Sprintf("Fifteen min: %.2f", stats.LoadAvg.FifteenMin),
		fmt.Sprintf("Running procs: %d", stats.LoadAvg.RunningProcs),
		fmt.Sprintf("Total procs: %d", stats.LoadAvg.TotalProcs),
		fmt.Sprintf("Last run PID: %d", stats.LoadAvg.LastPID),
	}, "\n")

	uptime := widgets.NewParagraph()
	uptime.Title = "Uptime"
	uptime.TitleAlignment = tui.AlignCenter
	uptime.TextAlignment = tui.AlignCenter
	uptime.VerticalAlignment = tui.AlignMiddle

	uptime.Text = fmt.Sprintf(
		"Total: %.2f    Idle: %.2f",
		stats.Uptime.Total,
		stats.Uptime.Idle,
	)

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
