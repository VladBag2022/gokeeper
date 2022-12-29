package client

import (
	"crypto/aes"
	"crypto/sha512"
	"encoding/hex"
	"fmt"

	"github.com/VladBag2022/gokeeper/internal/crypt"
)

// SessionManager encapsulates primary key storing logic.
type SessionManager struct {
	Coder *crypt.Coder

	encryptedKey []byte
	sessionKey   []byte
}

// NewSessionManagerFromPassword creates new SessionManager from provided password.
func NewSessionManagerFromPassword(password string) (*SessionManager, error) {
	sessionKey, err := crypt.GenerateRandomBytes(aes.BlockSize * 2)
	if err != nil {
		return nil, fmt.Errorf("failed to generate session key: %s", err)
	}

	tempCoder, err := crypt.NewCoder(sessionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp coder: %s", err)
	}

	key := sha512.Sum512([]byte(password))

	encryptedKey := tempCoder.Encrypt(key[:])

	primaryCoder, err := crypt.NewCoder(key[:])
	if err != nil {
		return nil, fmt.Errorf("failed to create coder: %s", err)
	}

	return &SessionManager{primaryCoder, encryptedKey, sessionKey}, nil
}

// NewSessionManagerFromEncryptedKey creates new SessionManager from encrypted key and session key (both hex-formatted).
func NewSessionManagerFromEncryptedKey(encryptedKeyHex, sessionKeyHex string) (*SessionManager, error) {
	encryptedKey, err := hex.DecodeString(encryptedKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex to encrypted key: %s", err)
	}

	sessionKey, err := hex.DecodeString(sessionKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex to session key: %s", err)
	}

	tempCoder, err := crypt.NewCoder(sessionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp coder: %s", err)
	}

	key, err := tempCoder.Decrypt(encryptedKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode encrypted key: %s", err)
	}

	primaryCoder, err := crypt.NewCoder(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create coder: %s", err)
	}

	return &SessionManager{primaryCoder, encryptedKey, sessionKey}, nil
}

// GetEncryptedKey returns encrypted key in hex format.
func (s *SessionManager) GetEncryptedKey() string {
	return hex.EncodeToString(s.encryptedKey)
}

// GetSessionKey returns session key in hex format.
func (s *SessionManager) GetSessionKey() string {
	return hex.EncodeToString(s.sessionKey)
}
