package pages

import (
	"fmt"

	tui "github.com/metaspartan/gotui/v5"
	"github.com/metaspartan/gotui/v5/widgets"
	"github.com/tinkyrain/go-sysmon/internal/metrics/system"
)

func SystemPage(stats system.Stats, x1, y1, x2, y2 int) *tui.Grid {
	systemInfo := widgets.NewParagraph()
	systemInfo.Title = "Info"
	systemInfo.Text = fmt.Sprintf(`
		Hostname: %s
		OS: %s
		Kernel: %s
		Arch: %s
	`,
		stats.Info.Hostname,
		stats.Info.OS,
		stats.Info.Kernel,
		stats.Info.Arch,
	)

	loadAvg := widgets.NewParagraph()
	loadAvg.Title = "LoadAvg"
	loadAvg.Text = fmt.Sprintf(`
		One min: %.2f
		Five min: %.2f
		Fifteen min: %.2f
		Running procs: %d
		Total procs: %d
		Last run PID: %d
	`,
		stats.LoadAvg.OneMin,
		stats.LoadAvg.FiveMin,
		stats.LoadAvg.FifteenMin,
		stats.LoadAvg.RunningProcs,
		stats.LoadAvg.TotalProcs,
		stats.LoadAvg.LastPID,
	)

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
			tui.NewCol(0.5, systemInfo),
			tui.NewCol(0.5, loadAvg),
		),
		tui.NewRow(
			0.15,
			tui.NewCol(1, uptime),
		),
	)

	return grid
}
