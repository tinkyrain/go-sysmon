package app

import (
	"context"
	"fmt"
	"go-sysmon/internal/config"
	"go-sysmon/internal/metrics"
	"time"
)

func Run(ctx context.Context, cfg config.Config, preview func(snapshot metrics.Snapshot)) error {
	collector := metrics.New(cfg.ProcRoot)
	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done(): // The code will be placed here, if a "notification" is received indicating that the application is stopping
			fmt.Println("Shutdown application")
			return nil
		case <-ticker.C: // Collect metrics
			snapshot, err := collector.Collect()
			if err != nil {
				return err
			}
			preview(snapshot)
		}
	}
}
