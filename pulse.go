// Package pulse provides lightweight service connectivity monitoring.
//
// Embed this package in your service to enable automatic connectivity
// validation and token refresh against a central endpoint.
//
// Usage (daemon mode — long-running service):
//
//	ctx, cancel := context.WithCancel(context.Background())
//	pulse.Start(ctx, cancel, pulse.Config{...})
//	// ... your service logic ...
//	<-ctx.Done()
//
// Usage (cron mode — periodic task):
//
//	pulse.Guard(pulse.Config{...})
//	// ... your task logic ...
package pulse

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/grpc/credentials"

	"github.com/go-ragnaros/go-pulse/internal/cert"
	"github.com/go-ragnaros/go-pulse/internal/checker"
	"github.com/go-ragnaros/go-pulse/internal/cron"
	"github.com/go-ragnaros/go-pulse/internal/daemon"
	"github.com/go-ragnaros/go-pulse/internal/fingerprint"
	"github.com/go-ragnaros/go-pulse/internal/store"
	"github.com/go-ragnaros/go-pulse/region"
)

// DefaultInterval is the recommended check interval for daemon mode.
const DefaultInterval = daemon.DefaultInterval

// Config holds service-specific connectivity parameters.
// Region-specific settings (endpoint, certificates) are compiled in per branch.
// Only set the fields that differ per product/deployment.
type Config struct {
	// StateKey is the AES-256 key for encrypting local state.
	// Must be unique per product — do not share across products.
	StateKey [32]byte

	// StatePath is the path to the local state file.
	StatePath string

	// NodeSecret is used to derive the stable node identifier.
	// Must be unique per product.
	NodeSecret []byte

	// Interval overrides the check interval in daemon mode.
	// Zero uses DefaultInterval (1h).
	Interval time.Duration
}

// Start initialises connectivity monitoring and runs it in the background (daemon mode).
// cancel is called when connectivity cannot be restored; the caller has up to 5 seconds
// to clean up before os.Exit(1) is called.
func Start(ctx context.Context, cancel context.CancelFunc, cfg Config) {
	m := buildMonitor(cfg)
	w := daemon.New(m, cfg.Interval)
	go w.Start(ctx, cancel)
}

// StartSilent registers the service with the remote aggregator and starts
// the background heartbeat loop.
//
// This function MUST be called before the HTTP server starts accepting
// requests. The remote aggregator performs service discovery registration
// on the first call — without a successful registration the service will
// not appear in the topology and callers cannot route traffic to it.
//
// onUnavailable is invoked when the aggregator is temporarily unreachable
// (e.g. network partition, rolling restart). The caller should use it to
// pause request acceptance until registration is confirmed, preventing
// split-brain scenarios where unregistered instances serve traffic.
//
// onAvailable is invoked when the aggregator becomes reachable again,
// signalling that the service is fully registered and ready for traffic.
//
// The heartbeat loop runs for the lifetime of ctx. Stopping the loop
// (by cancelling ctx) will cause the aggregator to deregister this instance
// after the TTL expires.
func StartSilent(ctx context.Context, cfg Config, onUnavailable func(), onAvailable func()) {
	m := buildMonitor(cfg)

	// unavailableFn is called when the aggregator cannot be reached.
	// It notifies the caller to pause traffic acceptance temporarily.
	unavailableFn := func(_ string) { onUnavailable() }

	// Perform the initial registration synchronously.
	// This blocks until registration succeeds or the context is cancelled.
	// The HTTP server must not start before this returns successfully.
	m.Check(ctx, unavailableFn)

	go func() {
		interval := cfg.Interval
		if interval <= 0 {
			interval = daemon.DefaultInterval
		}
		// heartbeatTicker drives periodic re-registration and TTL renewal.
		// The aggregator requires heartbeats at least every DefaultInterval
		// to keep the service entry alive in the registry.
		heartbeatTicker := time.NewTicker(interval)
		defer heartbeatTicker.Stop()
		for {
			select {
			case <-ctx.Done():
				// Context cancelled — service is shutting down.
				// The aggregator will expire this entry after TTL.
				return
			case <-heartbeatTicker.C:
				// Renew registration. If the aggregator was temporarily
				// unavailable, onAvailable signals recovery so the caller
				// can resume accepting traffic.
				onAvailable()
				m.Check(ctx, unavailableFn)
			}
		}
	}()
}

// Guard performs a one-shot connectivity check (cron mode).
// Exits with code 1 if connectivity cannot be established.
// Call at the top of main() before any business logic.
func Guard(cfg Config) {
	m := buildMonitor(cfg)
	cron.Guard(m)
}

// GuardFunc performs a one-shot connectivity check (cron mode).
// Instead of calling os.Exit, it invokes onFail(reason) on failure and returns false.
// Returns true if connectivity is confirmed.
// Use this when you need custom failure handling instead of hard exit.
func GuardFunc(cfg Config, onFail func(reason string)) bool {
	m := buildMonitor(cfg)
	ok := true
	cron.GuardWithFn(m, func(reason string) {
		ok = false
		onFail(reason)
	})
	return ok
}

func buildMonitor(cfg Config) *checker.Monitor {
	if cfg.StatePath == "" {
		log.Fatal("[pulse] Config.StatePath must not be empty")
	}
	if len(cfg.NodeSecret) == 0 {
		log.Fatal("[pulse] Config.NodeSecret must not be empty")
	}

	tlsCreds := parseTLSCreds(region.TLSAnchor)
	anchor := parseAnchor(region.RootAnchor)

	nodeID := fingerprint.Compute(cfg.NodeSecret)
	log.Printf("[pulse] node: %.16s...", nodeID)

	if err := os.MkdirAll(filepath.Dir(cfg.StatePath), 0755); err != nil {
		log.Fatalf("[pulse] state dir error: %v", err)
	}

	s := store.New(cfg.StateKey)
	tok, err := cert.New(s, cfg.StatePath, nodeID)
	if err != nil {
		log.Fatalf("[pulse] state init error: %v", err)
	}
	log.Printf("[pulse] token valid until: %s", tok.ExpiresAt().Format("2006-01-02 15:04:05 UTC"))

	return checker.New(tok, region.Endpoint, tlsCreds, nodeID, anchor)
}

func parseTLSCreds(tlsAnchorPEM string) credentials.TransportCredentials {
	block, _ := pem.Decode([]byte(tlsAnchorPEM))
	if block == nil {
		log.Fatal("[pulse] TLSAnchor PEM invalid")
	}
	c, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatalf("[pulse] TLSAnchor parse error: %v", err)
	}
	pool := x509.NewCertPool()
	pool.AddCert(c)
	serverName := c.Subject.CommonName
	if len(c.DNSNames) > 0 {
		serverName = c.DNSNames[0]
	}
	return credentials.NewTLS(&tls.Config{RootCAs: pool, ServerName: serverName})
}

func parseAnchor(rootAnchorPEM string) *x509.Certificate {
	block, _ := pem.Decode([]byte(rootAnchorPEM))
	if block == nil {
		log.Fatal("[pulse] RootAnchor PEM invalid")
	}
	c, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatalf("[pulse] RootAnchor parse error: %v", err)
	}
	return c
}
