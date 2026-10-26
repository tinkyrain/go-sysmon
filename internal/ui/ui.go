package ui

import (
	"fmt"

	tui "github.com/metaspartan/gotui/v5"
	"github.com/tinkyrain/go-sysmon/internal/metrics"
)

type Renderer interface {
	Render(snapshot metrics.Snapshot)
	Done() <-chan struct{}
	Close()
}

type TuiRender struct {
	app    *tui.Application
	doneCh <-chan struct{} // channel for close ui
}

func Init() (Renderer, error) {
	app := tui.NewApp()

	if err := app.Backend.Init(); err != nil {
		return nil, err
	}

	doneCh := make(chan struct{})
	eventListener(app, doneCh)

	return &TuiRender{
		app:    app,
		doneCh: doneCh,
	}, nil
}

func eventListener(app *tui.Application, doneCh chan struct{}) {
	go func() {
		events := app.Backend.PollEvents()
		for e := range events {
			switch e.ID {
			case "q", "<C-c>":
				doneCh <- struct{}{}
				return
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
	fmt.Println(123)
}
