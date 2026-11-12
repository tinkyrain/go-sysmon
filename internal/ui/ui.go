package ui

import (
	"sync"

	tui "github.com/metaspartan/gotui/v5"
	"github.com/metaspartan/gotui/v5/widgets"
	"github.com/tinkyrain/go-sysmon/internal/metrics"
	"github.com/tinkyrain/go-sysmon/internal/ui/components"
	"github.com/tinkyrain/go-sysmon/internal/ui/pages"
	"github.com/tinkyrain/go-sysmon/internal/ui/view"
)

type Renderer interface {
	Render(snapshot metrics.Snapshot)
	Done() <-chan struct{}
	Close()
}

type tuiRenderer struct {
	once sync.Once
	mu   sync.Mutex

	app    *tui.Application
	doneCh chan struct{} // channel for close ui

	tab      *widgets.TabPane
	lastData view.RenderData
}

func New() (Renderer, error) {
	app := tui.NewApp()

	if err := app.Backend.Init(); err != nil {
		return nil, err
	}

	doneCh := make(chan struct{})

	return &tuiRenderer{
		app:    app,
		doneCh: doneCh,
	}, nil
}

func (r *tuiRenderer) eventListener() {
	go func() {
		events := r.app.Backend.PollEvents()
		for e := range events {
			switch e.ID {
			case "q", "Q", "<C-c>":
				close(r.doneCh)
				return
			case "<Left>":
				r.redraw(func(tab *widgets.TabPane) { tab.FocusLeft() })
			case "<Right>":
				r.redraw(func(tab *widgets.TabPane) { tab.FocusRight() })
			case "<Tab>":
				r.redraw(func(tab *widgets.TabPane) {
					if tab.ActiveTabIndex == len(tab.TabNames)-1 {
						tab.ActiveTabIndex = 0
					} else {
						tab.FocusRight()
					}
				})
			case "<Resize>":
				r.redraw(nil)
			}
		}
	}()
}

func (r *tuiRenderer) Close() {
	r.app.Backend.Close()
}

func (r *tuiRenderer) Done() <-chan struct{} {
	return r.doneCh
}

func (r *tuiRenderer) Render(snapshot metrics.Snapshot) {
	r.draw(view.NewRenderData(snapshot))
}

func (r *tuiRenderer) redraw(change func(tab *widgets.TabPane)) {
	r.mu.Lock()
	if change != nil && r.tab != nil {
		change(r.tab)
	}
	data := r.lastData
	r.mu.Unlock()

	r.draw(data)
}

func (r *tuiRenderer) draw(data view.RenderData) {
	w, h := tui.TerminalDimensions()
	hintsHeight := 3
	tabHeight := 2
	contentStartX := 2             // take an offset from the main frame
	contentStartY := tabHeight + 1 // take an offset from tabs
	contentEndX := w - contentStartX
	contentEndY := h - hintsHeight - 1 // take an offset from the main frame and hints

	r.mu.Lock()
	defer r.mu.Unlock()

	if data.Timestamp > r.lastData.Timestamp {
		r.lastData = data
	}

	layout := components.Layout(r.lastData.Time, 0, 0, w, h-hintsHeight)
	hints := components.Hints(0, h-hintsHeight, w, h)
	tab := components.Tabs(0, w, tabHeight)
	if r.tab != nil {
		tab.ActiveTabIndex = r.tab.ActiveTabIndex
	}
	r.tab = tab

	r.once.Do(r.eventListener)
	r.app.Backend.Render(layout, tab, hints, r.content(r.lastData, contentStartX, contentStartY, contentEndX, contentEndY))
}

func (r *tuiRenderer) content(data view.RenderData, x1, y1, x2, y2 int) tui.Drawable {
	var content tui.Drawable
	switch r.tab.ActiveTabIndex {
	case 0:
		content = pages.SystemPage(data.System, x1, y1, x2, y2)
	case 1:
		content = pages.CPUPage(data.CPU, x1, y1, x2, y2)
	case 2:
		content = pages.MemoryPage(data.Memory, x1, y1, x2, y2)
	case 3:
		content = pages.DiskPage(data.Disk, x1, y1, x2, y2)
	}
	return content
}
