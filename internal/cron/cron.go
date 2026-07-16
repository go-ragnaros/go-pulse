// Package cron provides one-shot connectivity validation for periodic tasks.
package cron

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-ragnaros/go-pulse/internal/checker"
)

const checkTimeout = 30 * time.Second

// Guard performs a one-shot connectivity check at process startup.
// Exits with code 1 if the service is unreachable or token refresh fails.
func Guard(m *checker.Monitor) {
	GuardWithFn(m, func(reason string) {
		log.Printf("[pulse] connectivity check failed, exiting: %s", reason)
		os.Exit(1)
	})
}

// GuardWithFn is the testable variant: shutdownFn replaces os.Exit(1).
func GuardWithFn(m *checker.Monitor, shutdownFn func(reason string)) {
	ctx, cancel := context.WithTimeout(context.Background(), checkTimeout)
	defer cancel()
	m.Check(ctx, shutdownFn)
}
