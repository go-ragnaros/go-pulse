// Package daemon provides background connectivity monitoring for long-running services.
package daemon

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-ragnaros/go-pulse/internal/checker"
)

const (
	DefaultInterval = time.Hour
	drainTimeout    = 5 * time.Second
)

// Watcher runs periodic connectivity checks in the background.
type Watcher struct {
	monitor  *checker.Monitor
	interval time.Duration
}

// New creates a Watcher.
func New(m *checker.Monitor, interval time.Duration) *Watcher {
	if interval <= 0 {
		interval = DefaultInterval
	}
	return &Watcher{monitor: m, interval: interval}
}

// Start launches the background check loop.
// On failure, cancel() is called to initiate graceful shutdown,
// followed by os.Exit(1) after drainTimeout.
func (w *Watcher) Start(ctx context.Context, cancel context.CancelFunc) {
	w.StartWithShutdown(ctx, cancel, func() { os.Exit(1) })
}

// StartWithShutdown is the testable variant: exitFn replaces os.Exit(1).
func (w *Watcher) StartWithShutdown(ctx context.Context, cancel context.CancelFunc, exitFn func()) {
	shutdownFn := func(reason string) {
		log.Printf("[pulse] initiating shutdown: %s", reason)
		cancel()
		select {
		case <-ctx.Done():
		case <-time.After(drainTimeout):
			log.Println("[pulse] drain timeout, forcing exit")
		}
		exitFn()
	}

	w.monitor.Check(ctx, shutdownFn)

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Println("[pulse] watcher stopped")
			return
		case <-ticker.C:
			w.monitor.Check(ctx, shutdownFn)
		}
	}
}
