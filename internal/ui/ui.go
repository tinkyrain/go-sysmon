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

type TuiRender struct {
	once sync.Once
	mu   sync.Mutex

	app    *tui.Application
	doneCh chan struct{} // channel for close ui

	tab      *widgets.TabPane
	lastData view.RenderData
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
				data := r.lastData
				r.mu.Unlock()

				r.draw(data)
			case "<Right>":
				r.mu.Lock()
				if r.tab != nil {
					r.tab.FocusRight()
				}
				data := r.lastData
				r.mu.Unlock()

				r.draw(data)
			case "<Tab>":
				r.mu.Lock()
				if r.tab != nil {
					if r.tab.ActiveTabIndex == len(r.tab.TabNames)-1 {
						r.tab.ActiveTabIndex = 0
					} else {
						r.tab.FocusRight()
					}
				}
				data := r.lastData
				r.mu.Unlock()

				r.draw(data)
			case "<Resize>":
				r.mu.Lock()
				data := r.lastData
				r.mu.Unlock()

				r.draw(data)
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
	r.draw(view.NewRenderData(snapshot))
}

func (r *TuiRender) draw(data view.RenderData) {
	w, h := tui.TerminalDimensions()
	hintsHeight := 3
	tabHeight := 2
	contentStartX := 2             // take an offset from the main frame
	contentStartY := tabHeight + 1 // take an offset from tabs
	contentEndX := w - contentStartX
	contentEndY := h - hintsHeight - 1 // take an offset from the main frame and hints

	r.mu.Lock()
	defer r.mu.Unlock()

	r.lastData = data

	layout := components.Layout(data.Time, 0, 0, w, h-hintsHeight)
	hints := components.Hint(0, h-hintsHeight, w, h)
	tab := components.Tab(0, w, tabHeight)
	if r.tab != nil {
		tab.ActiveTabIndex = r.tab.ActiveTabIndex
	}
	r.tab = tab

	r.once.Do(r.eventListener)
	r.app.Backend.Render(layout, tab, hints, r.content(data, contentStartX, contentStartY, contentEndX, contentEndY))
}

func (r *TuiRender) content(data view.RenderData, x1, y1, x2, y2 int) tui.Drawable {
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
	default:
		content = pages.SystemPage(data.System, x1, y1, x2, y2)
	}
	return content
}
