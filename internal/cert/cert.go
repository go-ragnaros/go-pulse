// Package cert manages local encrypted token storage and validation.
package cert

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"

	"github.com/go-ragnaros/go-pulse/internal/store"
)

var ErrTokenExpired = errors.New("cert: token expired")
var ErrNodeMismatch = errors.New("cert: node identifier mismatch")
var ErrInvalidToken = errors.New("cert: token invalid or not trusted")

type tokenState struct {
	CertPEM []byte `json:"cert_pem"`
	KeyPEM  []byte `json:"key_pem"`
	NodeID  string `json:"node_id"`
}

// LocalToken holds the local connectivity token.
type LocalToken struct {
	cert   *x509.Certificate
	key    *ecdsa.PrivateKey
	nodeID string
	s      *store.StateStore
	path   string
}

// New loads or initialises the local token at path.
func New(s *store.StateStore, path string, nodeID string) (*LocalToken, error) {
	cc := &LocalToken{s: s, path: path, nodeID: nodeID}
	data, err := s.Read(path)
	if os.IsNotExist(err) {
		return cc.generatePlaceholder()
	}
	if err != nil {
		return nil, err
	}
	var st tokenState
	if err := json.Unmarshal(data, &st); err != nil {
		return nil, err
	}
	return cc, cc.load(st)
}

func (c *LocalToken) IsExpired() bool { return c.IsExpiredAt(time.Now().UTC()) }

func (c *LocalToken) IsExpiredAt(t time.Time) bool { return t.After(c.cert.NotAfter) }

func (c *LocalToken) VerifyAt(anchor *x509.Certificate, t time.Time) error {
	if err := c.cert.CheckSignatureFrom(anchor); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}
	if t.Before(c.cert.NotBefore) || t.After(c.cert.NotAfter) {
		return fmt.Errorf("%w: not valid at %v", ErrInvalidToken, t)
	}
	return nil
}

func (c *LocalToken) RequestBytes() ([]byte, error) {
	csr, err := x509.CreateCertificateRequest(rand.Reader, &x509.CertificateRequest{
		Subject: pkix.Name{CommonName: c.nodeID},
	}, c.key)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csr}), nil
}

func (c *LocalToken) Update(certPEM []byte) error {
	block, _ := pem.Decode(certPEM)
	if block == nil {
		return fmt.Errorf("%w: server returned invalid PEM", ErrInvalidToken)
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return err
	}
	c.cert = cert
	return c.persist()
}

func (c *LocalToken) ExpiresAt() time.Time { return c.cert.NotAfter }

func (c *LocalToken) generatePlaceholder() (*LocalToken, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: c.nodeID},
		NotBefore:    time.Now().UTC().Add(-time.Hour),
		NotAfter:     time.Now().UTC().Add(-time.Second),
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
	}
	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		return nil, err
	}
	c.key = key
	c.cert, err = x509.ParseCertificate(certDER)
	if err != nil {
		return nil, err
	}
	return c, c.persist()
}

func (c *LocalToken) load(st tokenState) error {
	certBlock, _ := pem.Decode(st.CertPEM)
	if certBlock == nil {
		return fmt.Errorf("%w: stored cert PEM invalid", ErrInvalidToken)
	}
	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return err
	}
	keyBlock, _ := pem.Decode(st.KeyPEM)
	if keyBlock == nil {
		return fmt.Errorf("%w: stored key PEM invalid", ErrInvalidToken)
	}
	key, err := x509.ParseECPrivateKey(keyBlock.Bytes)
	if err != nil {
		return err
	}
	c.cert, c.key, c.nodeID = cert, key, st.NodeID
	return nil
}

func (c *LocalToken) persist() error {
	keyDER, err := x509.MarshalECPrivateKey(c.key)
	if err != nil {
		return err
	}
	st := tokenState{
		CertPEM: pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: c.cert.Raw}),
		KeyPEM:  pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER}),
		NodeID:  c.nodeID,
	}
	data, err := json.Marshal(st)
	if err != nil {
		return err
	}
	return c.s.Write(c.path, data)
}
