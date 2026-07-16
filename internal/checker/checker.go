// Package checker implements the connectivity validation and token refresh flow.
package checker

import (
	"context"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/go-ragnaros/go-pulse/internal/cert"
	internalntp "github.com/go-ragnaros/go-pulse/internal/ntp"
	pb "github.com/go-ragnaros/go-pulse/proto/licensepb"
)

// Monitor performs time-validated, signature-verified connectivity checks.
type Monitor struct {
	token    *cert.LocalToken
	endpoint string
	creds    credentials.TransportCredentials
	nodeID   string
	anchor   *x509.Certificate
}

// New creates a Monitor.
func New(token *cert.LocalToken, endpoint string, creds credentials.TransportCredentials, nodeID string, anchor *x509.Certificate) *Monitor {
	return &Monitor{
		token:    token,
		endpoint: endpoint,
		creds:    creds,
		nodeID:   nodeID,
		anchor:   anchor,
	}
}

// Check runs the full connectivity validation flow.
// On failure, shutdownFn(reason) is called.
func (m *Monitor) Check(ctx context.Context, shutdownFn func(reason string)) {
	ntpTime, err := internalntp.ValidateDrift()
	if err != nil {
		var reason string
		switch {
		case errors.Is(err, internalntp.ErrClockDrift):
			reason = fmt.Sprintf("time sync drift detected: %v", err)
		case errors.Is(err, internalntp.ErrAllServersFailed):
			reason = fmt.Sprintf("time sync unreachable: %v", err)
		default:
			reason = fmt.Sprintf("time sync error: %v", err)
		}
		log.Printf("[pulse] %s", reason)
		shutdownFn(reason)
		return
	}

	expired := m.token.IsExpiredAt(ntpTime)
	sigErr := m.token.VerifyAt(m.anchor, ntpTime)

	if !expired && sigErr == nil {
		return
	}

	if expired {
		log.Printf("[pulse] token expired, refreshing")
	} else {
		log.Printf("[pulse] token validation failed, refreshing: %v", sigErr)
	}

	if err := m.refresh(ctx, ntpTime); err != nil {
		reason := fmt.Sprintf("service connectivity check failed: %v", err)
		log.Printf("[pulse] %s", reason)
		shutdownFn(reason)
	}
}

func (m *Monitor) refresh(ctx context.Context, ntpTime time.Time) error {
	conn, err := grpc.NewClient(m.endpoint, grpc.WithTransportCredentials(m.creds))
	if err != nil {
		return fmt.Errorf("endpoint unreachable: %w", err)
	}
	defer conn.Close()

	client := pb.NewHeartbeatServiceClient(conn)
	reqPEM, err := m.token.RequestBytes()
	if err != nil {
		return fmt.Errorf("prepare request: %w", err)
	}
	resp, err := client.Sync(ctx, &pb.SyncRequest{
		NodeId:   m.nodeID,
		TokenPem: reqPEM,
	})
	if err != nil {
		return fmt.Errorf("service refused: %w", err)
	}
	if err := m.token.Update(resp.TokenPem); err != nil {
		return fmt.Errorf("update token: %w", err)
	}
	if err := m.token.VerifyAt(m.anchor, ntpTime); err != nil {
		return fmt.Errorf("token anchor mismatch (endpoint changed?): %w", err)
	}
	return nil
}
