// Package store provides AES-GCM encrypted local state persistence.
package store

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"os"
)

// StateStore reads and writes AES-GCM encrypted files.
// The key must be exactly 32 bytes (AES-256).
type StateStore struct {
	key [32]byte
}

// New creates a StateStore with the given 32-byte AES key.
func New(key [32]byte) *StateStore {
	return &StateStore{key: key}
}

// Write encrypts plaintext and writes it to path atomically (tmp → rename).
func (s *StateStore) Write(path string, plaintext []byte) error {
	block, err := aes.NewCipher(s.key[:])
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, ciphertext, 0600); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// Read decrypts and returns plaintext from path.
func (s *StateStore) Read(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(s.key[:])
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(data) < gcm.NonceSize() {
		return nil, errors.New("store: ciphertext too short")
	}
	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
