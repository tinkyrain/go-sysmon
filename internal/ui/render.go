package ui

import (
	"fmt"
	"strings"

	"go-sysmon/internal/metrics"

	"github.com/charmbracelet/lipgloss"
	"github.com/pterm/pterm"
)

var (
	frameStyle = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1)
	boxStyle   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1)

	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	sectionStyle = lipgloss.NewStyle().Bold(true).Underline(true)
	dimStyle     = lipgloss.NewStyle().Faint(true)
)

func gauge(percent float64, width int) string {
	percent = clamp(percent, 0, 100)
	if width < 1 {
		width = 1
	}
	filled := int(percent / 100 * float64(width))
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)

	color := lipgloss.Color("2") // green
	switch {
	case percent >= 85:
		color = lipgloss.Color("1") // red
	case percent >= 60:
		color = lipgloss.Color("3") // yellow
	}
	return lipgloss.NewStyle().Foreground(color).Render(bar) + fmt.Sprintf(" %5.1f%%", percent)
}

func meter(label string, total, available uint64, gaugeW int) string {
	return strings.Join([]string{
		label,
		gauge(usedPercent(total, available), gaugeW),
		dimStyle.Render(humanizeSize(float64(total-available)) + " / " + humanizeSize(float64(total))),
	}, "\n")
}

func render(s metrics.Snapshot) string {
	w, h := terminalSize()

	const cols = 3
	bodyW := w - 4                      // inside frame: border(2) + padding(2)
	boxW := (bodyW - (cols - 1)) / cols // per-box outer width, minus 1-col gaps
	boxW = min(max(boxW, 16), 44)       // keep boxes sane on tiny/huge terminals
	textW := boxW - 4                   // inside box: border(2) + padding(2)
	gaugeW := max(textW-7, 6)           // reserve " 100.0%"

	bs := boxStyle.Width(boxW - 2)

	cpu := bs.Render(sectionStyle.Render("CPU") + "\n\n" + gauge(s.CPUUsage, gaugeW))

	mem := bs.Render(sectionStyle.Render("Memory") + "\n\n" +
		meter("RAM", s.Memory.Total, s.Memory.Available, gaugeW) + "\n\n" +
		meter("SWAP", s.Memory.SwapTotal, s.Memory.SwapAvailable, gaugeW))

	var d strings.Builder
	d.WriteString(sectionStyle.Render("Disks"))
	for i, disk := range s.Disks {
		if i > 0 {
			d.WriteString("\n" + dimStyle.Render(strings.Repeat("─", textW)))
		}
		d.WriteString("\n\n" + meter(disk.Mount, disk.Total, disk.Available, gaugeW))
	}
	disks := bs.Render(d.String())

	row := lipgloss.JoinHorizontal(lipgloss.Top, cpu, " ", mem, " ", disks)

	header := titleStyle.Render("go-sysmon") + "  " +
		dimStyle.Render(s.Time.Format("2006-01-02 15:04:05"))

	body := lipgloss.JoinVertical(lipgloss.Left, header, "", row)

	return frameStyle.Width(w - 2).Height(max(h-4, 3)).Render(body)
}

func terminalSize() (int, int) {
	w, h := pterm.GetTerminalWidth(), pterm.GetTerminalHeight()
	if w <= 0 {
		w = 80
	}
	if h <= 0 {
		h = 24
	}
	return w, h
}

func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
