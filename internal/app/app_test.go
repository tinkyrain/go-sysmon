package app_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tinkyrain/go-sysmon/internal/app"
	"github.com/tinkyrain/go-sysmon/internal/config"
	"github.com/tinkyrain/go-sysmon/internal/metrics"
)

const (
	tickInterval  = 10 * time.Millisecond
	neverTick     = time.Hour
	safetyTimeout = 2 * time.Second
)

var errCollect = errors.New("collect failed")

type mockRenderer struct {
	mock.Mock
}

func (m *mockRenderer) Render(snapshot metrics.Snapshot) {
	m.Called(snapshot)
}

func (m *mockRenderer) Done() <-chan struct{} {
	args := m.Called()
	return args.Get(0).(<-chan struct{})
}

func (m *mockRenderer) Close() {
	m.Called()
}

func newMockRenderer(done <-chan struct{}) *mockRenderer {
	renderer := &mockRenderer{}
	renderer.On("Done").Return(done)

	return renderer
}

func blockedDone() <-chan struct{} {
	return make(<-chan struct{})
}

func unusedCollect(t *testing.T) func() (metrics.Snapshot, error) {
	t.Helper()

	return func() (metrics.Snapshot, error) {
		t.Error("collect must not be called")

		return metrics.Snapshot{}, nil
	}
}

func testSnapshot() metrics.Snapshot {
	return metrics.Snapshot{
		Time: time.Date(2026, time.October, 26, 12, 0, 0, 0, time.UTC),
	}
}

func TestRunContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	renderer := newMockRenderer(blockedDone())

	err := app.Run(ctx, config.Config{Interval: neverTick}, renderer, unusedCollect(t))

	require.NoError(t, err)
	renderer.AssertNumberOfCalls(t, "Render", 0)
	renderer.AssertNumberOfCalls(t, "Close", 0)
}

func TestRunRendererDone(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), safetyTimeout)
	defer cancel()

	done := make(chan struct{})
	close(done)
	renderer := newMockRenderer(done)

	err := app.Run(ctx, config.Config{Interval: neverTick}, renderer, unusedCollect(t))

	require.NoError(t, err)
	renderer.AssertNumberOfCalls(t, "Render", 0)
	renderer.AssertNumberOfCalls(t, "Close", 0)
}

func TestRunRenderSnapshot(t *testing.T) {
	expected := testSnapshot()

	ctx, cancel := context.WithTimeout(context.Background(), safetyTimeout)
	defer cancel()

	renderer := newMockRenderer(blockedDone())
	renderer.On("Render", mock.Anything).Run(func(mock.Arguments) { cancel() })

	collect := func() (metrics.Snapshot, error) {
		return expected, nil
	}

	err := app.Run(ctx, config.Config{Interval: tickInterval}, renderer, collect)

	require.NoError(t, err)
	renderer.AssertCalled(t, "Render", expected)
	renderer.AssertNumberOfCalls(t, "Close", 0)
}

func TestRunCollectError(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), safetyTimeout)
	defer cancel()

	renderer := newMockRenderer(blockedDone())

	collect := func() (metrics.Snapshot, error) {
		return metrics.Snapshot{}, errCollect
	}

	err := app.Run(ctx, config.Config{Interval: tickInterval}, renderer, collect)

	require.ErrorIs(t, err, errCollect)
	renderer.AssertNumberOfCalls(t, "Render", 0)
	renderer.AssertNumberOfCalls(t, "Close", 0)
}

func TestRunCollectErrorInLoop(t *testing.T) {
	expected := testSnapshot()

	ctx, cancel := context.WithTimeout(context.Background(), safetyTimeout)
	defer cancel()

	renderer := newMockRenderer(blockedDone())
	renderer.On("Render", mock.Anything).Return()

	calls := 0
	collect := func() (metrics.Snapshot, error) {
		calls++
		if calls == 1 {
			return expected, nil
		}

		return metrics.Snapshot{}, errCollect
	}

	err := app.Run(ctx, config.Config{Interval: tickInterval}, renderer, collect)

	require.ErrorIs(t, err, errCollect)
	assert.Equal(t, 2, calls)
	renderer.AssertNumberOfCalls(t, "Render", 1)
	renderer.AssertCalled(t, "Render", expected)
}
