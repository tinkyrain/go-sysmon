package app_test

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tinkyrain/go-sysmon/internal/app"
	"github.com/tinkyrain/go-sysmon/internal/config"
	"github.com/tinkyrain/go-sysmon/internal/metrics"
)

func TestRunSuccessPreview(t *testing.T) {
	config, err := config.Load("")
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), config.Interval)
	defer cancel()

	config.Interval = 1 * time.Millisecond

	ctx, stopSignals := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stopSignals()

	collect := func() (metrics.Snapshot, error) {
		return metrics.Snapshot{}, nil
	}

	previewCalled := false
	preview := func(metrics.Snapshot) {
		previewCalled = true
		stopSignals()
	}

	_ = app.Run(ctx, config, collect, preview)

	assert.Equal(t, true, previewCalled)
}

func TestRunCollectError(t *testing.T) {
	config, err := config.Load("")
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), config.Interval)
	defer cancel()

	config.Interval = 1 * time.Millisecond

	ctx, stopSignals := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stopSignals()

	collect := func() (metrics.Snapshot, error) {
		return metrics.Snapshot{}, fmt.Errorf("Error")
	}

	preview := func(metrics.Snapshot) {}

	error := app.Run(ctx, config, collect, preview)

	assert.Error(t, error)
}
