package ui

import (
	"sync"

	tui "github.com/metaspartan/gotui/v5"
	"github.com/metaspartan/gotui/v5/widgets"
	"github.com/tinkyrain/go-sysmon/internal/metrics"
	"github.com/tinkyrain/go-sysmon/internal/ui/components"
)

type Renderer interface {
	Render(snapshot metrics.Snapshot)
	Done() <-chan struct{}
	Close()
}

type TuiRender struct {
	app       *tui.Application
	doneCh    chan struct{} // channel for close ui
	tab       *widgets.TabPane
	activeTab int
	once      sync.Once
}

func Init() (Renderer, error) {
	app := tui.NewApp()

	if err := app.Backend.Init(); err != nil {
		return nil, err
	}

	doneCh := make(chan struct{})

	return &TuiRender{
		app:    app,
		doneCh: doneCh,
	}, nil
}

func (r *TuiRender) eventListener() {
	if r.tab == nil {
		panic("Tab is not init")
	}

	go func() {
		events := r.app.Backend.PollEvents()
		for e := range events {
			switch e.ID {
			case "q", "<C-c>":
				r.doneCh <- struct{}{}
				return
			case "<Left>":
				r.tab.FocusLeft()
			case "<Right>", "<Tab>":
				r.tab.FocusRight()
			}

			r.app.Backend.Clear()
		}
	}()
}

func (r *TuiRender) Close() {
	r.app.Backend.Close()
}

func (r *TuiRender) Done() <-chan struct{} {
	return r.doneCh
}

func (r *TuiRender) Render(snapshot metrics.Snapshot) {
	w, h := tui.TerminalDimensions()
	hintsHeight := 3

	layout := components.Layout(snapshot.Time, 0, 0, w, h-hintsHeight)
	hints := components.Hint(0, h-hintsHeight, w, h)

	tab := components.Tab(0, w, 2)
	if r.tab != nil {
		tab.ActiveTabIndex = r.tab.ActiveTabIndex
	}
	r.tab = tab

	r.once.Do(r.eventListener)

	switch r.activeTab {
	case 0:
		r.tab.ActiveTabStyle = tui.NewStyle(tui.ColorWhite, tui.ColorBlue)
	case 1:
		r.tab.ActiveTabStyle = tui.NewStyle(tui.ColorWhite, tui.ColorBlue)
	}

	r.app.Backend.Render(layout, tab, hints)
}
