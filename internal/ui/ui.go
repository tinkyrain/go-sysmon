package ui

import (
	"sync"

	tui "github.com/metaspartan/gotui/v5"
	"github.com/metaspartan/gotui/v5/widgets"
	"github.com/tinkyrain/go-sysmon/internal/metrics"
	"github.com/tinkyrain/go-sysmon/internal/ui/components"
	"github.com/tinkyrain/go-sysmon/internal/ui/pages"
)

type Renderer interface {
	Render(snapshot metrics.Snapshot)
	Done() <-chan struct{}
	Close()
}

type TuiRender struct {
	once sync.Once
	mu   sync.Mutex

	app    *tui.Application
	doneCh chan struct{} // channel for close ui

	tab      *widgets.TabPane
	lastSnap metrics.Snapshot
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
	go func() {
		events := r.app.Backend.PollEvents()
		for e := range events {
			switch e.ID {
			case "q", "Q", "<C-c>":
				r.doneCh <- struct{}{}
				return
			case "<Left>":
				r.mu.Lock()
				if r.tab != nil {
					r.tab.FocusLeft()
				}
				snap := r.lastSnap
				r.mu.Unlock()

				r.Render(snap)
			case "<Right>":
				r.mu.Lock()
				if r.tab != nil {
					r.tab.FocusRight()
				}
				snap := r.lastSnap
				r.mu.Unlock()

				r.Render(snap)
			case "<Tab>":
				r.mu.Lock()
				if r.tab != nil {
					if r.tab.ActiveTabIndex == len(r.tab.TabNames)-1 {
						r.tab.ActiveTabIndex = 0
					} else {
						r.tab.FocusRight()
					}
				}
				snap := r.lastSnap
				r.mu.Unlock()

				r.Render(snap)
			}
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
	tabHeight := 2
	contentStartX := 2             // take an offset from the main frame
	contentStartY := tabHeight + 1 // take an offset from tabs
	contentEndX := w - contentStartX
	contentEndY := h - tabHeight - 2 // take an offset from the main frame and hints

	r.mu.Lock()

	r.lastSnap = snapshot

	layout := components.Layout(snapshot.Time, 0, 0, w, h-hintsHeight)
	hints := components.Hint(0, h-hintsHeight, w, h)
	tab := components.Tab(0, w, tabHeight)
	if r.tab != nil {
		tab.ActiveTabIndex = r.tab.ActiveTabIndex
	}
	r.tab = tab

	r.mu.Unlock()

	content := pages.SystemPage(snapshot.System, contentStartX, contentStartY, contentEndX, contentEndY)

	r.once.Do(r.eventListener)
	r.app.Backend.Render(layout, tab, hints, content)
}
