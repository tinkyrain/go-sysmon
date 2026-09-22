package ui

import (
	"fmt"
	"os"

	"github.com/tinkyrain/go-sysmon/internal/metrics"
)

const (
	altScreenOn  = "\x1b[?1049h"
	altScreenOff = "\x1b[?1049l"
	cursorOff    = "\x1b[?25l"
	cursorOn     = "\x1b[?25h"
	clearHome    = "\x1b[2J\x1b[H"
	home         = "\x1b[H"
)

type Dashboard struct {
	out *os.File
}

func NewDashboard() *Dashboard {
	_, _ = fmt.Fprint(os.Stdout, altScreenOn+cursorOff+clearHome)
	return &Dashboard{out: os.Stdout}
}

func (d *Dashboard) Render(s metrics.Snapshot) {
	_, _ = fmt.Fprint(d.out, home+render(s))
}

func (d *Dashboard) Stop() {
	_, _ = fmt.Fprint(d.out, altScreenOff+cursorOn)
}
