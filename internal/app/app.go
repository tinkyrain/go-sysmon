package app

import (
	"context"
	"time"

	"github.com/tinkyrain/go-sysmon/internal/config"
	"github.com/tinkyrain/go-sysmon/internal/metrics"
)

func Run(
	ctx context.Context,
	cfg config.Config,
	collect func() (metrics.Snapshot, error),
	preview func(snapshot metrics.Snapshot),
) error {
	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done(): // The code will be placed here, if a "notification" is received indicating that the application is stopping
			return nil
		case <-ticker.C: // Collect metrics
			snapshot, err := collect()
			if err != nil {
				return err
			}
			preview(snapshot)
		}
	}
}
