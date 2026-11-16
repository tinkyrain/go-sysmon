package app

import (
	"context"
	"time"

	"github.com/tinkyrain/go-sysmon/internal/config"
	"github.com/tinkyrain/go-sysmon/internal/metrics"
	"github.com/tinkyrain/go-sysmon/internal/ui"
	"github.com/tinkyrain/go-sysmon/internal/ui/view"
)

func Run(
	ctx context.Context,
	cfg config.Config,
	renderer ui.Renderer,
	collect func() (metrics.Snapshot, error),
) error {
	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	if err := render(collect, renderer); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-renderer.Done():
			return nil
		case <-ticker.C:
			if err := render(collect, renderer); err != nil {
				return err
			}
		}
	}
}

func render(collect func() (metrics.Snapshot, error), renderer ui.Renderer) error {
	snapshot, err := collect()
	if err != nil {
		return err
	}
	renderer.Render(view.NewRenderData(snapshot))
	return nil
}
