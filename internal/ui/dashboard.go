package ui

import (
	"go-sysmon/internal/metrics"

	"github.com/pterm/pterm"
)

type Dashboard struct {
	area *pterm.AreaPrinter
}

func NewDashboard() (*Dashboard, error) {
	area, err := pterm.DefaultArea.WithFullscreen().Start()
	if err != nil {
		return nil, err
	}
	return &Dashboard{area: area}, nil
}

func (d *Dashboard) Render(s metrics.Snapshot) {
	d.area.Update(render(s))
}

func (d *Dashboard) Stop() { _ = d.area.Stop() }
