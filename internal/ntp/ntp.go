// Package ntp provides authoritative time via public NTP servers.
package ntp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"time"
)

const (
	maxRetries = 3
	MaxDrift   = 5 * time.Minute
	ntpEpoch   = 2208988800
)

var ErrAllServersFailed = errors.New("ntp: all servers failed after retries")
var ErrClockDrift = errors.New("ntp: local clock drift exceeds allowed threshold")

var defaultServers = []string{
	"time1.aliyun.com:123",
	"cn.pool.ntp.org:123",
	"pool.ntp.org:123",
}

func Now() (time.Time, error) {
	return QueryChain(defaultServers)
}

func QueryChain(servers []string) (time.Time, error) {
	var lastErr error
	for _, server := range servers {
		for i := range maxRetries {
			t, err := queryNTP(server)
			if err == nil {
				return t, nil
			}
			lastErr = fmt.Errorf("ntp: server %s attempt %d: %w", server, i+1, err)
		}
	}
	return time.Time{}, fmt.Errorf("%w: last error: %v", ErrAllServersFailed, lastErr)
}

func ValidateDrift() (time.Time, error) {
	return ValidateDriftFrom(defaultServers)
}

func ValidateDriftFrom(servers []string) (time.Time, error) {
	ntpTime, err := QueryChain(servers)
	if err != nil {
		return time.Time{}, err
	}
	drift := time.Since(ntpTime)
	if drift < 0 {
		drift = -drift
	}
	if drift > MaxDrift {
		return time.Time{}, fmt.Errorf("%w: drift=%v", ErrClockDrift, drift)
	}
	return ntpTime, nil
}

func queryNTP(server string) (time.Time, error) {
	conn, err := net.DialTimeout("udp", server, 3*time.Second)
	if err != nil {
		return time.Time{}, err
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(3 * time.Second))
	req := make([]byte, 48)
	req[0] = 0x1B
	if _, err := conn.Write(req); err != nil {
		return time.Time{}, err
	}
	resp := make([]byte, 48)
	if _, err := conn.Read(resp); err != nil {
		return time.Time{}, err
	}
	secs := binary.BigEndian.Uint32(resp[40:44])
	return time.Unix(int64(secs)-ntpEpoch, 0).UTC(), nil
}
