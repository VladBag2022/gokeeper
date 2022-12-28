package client

import (
	"crypto/aes"
	"crypto/sha512"
	"encoding/hex"

	log "github.com/sirupsen/logrus"

	"github.com/VladBag2022/gokeeper/internal/crypt"
)

type SessionManager struct {
	Coder *crypt.Coder

	encryptedKey []byte
	sessionKey   []byte
}

func NewSessionManagerFromPassword(password string) (*SessionManager, error) {
	sessionKey, err := crypt.GenerateRandomBytes(aes.BlockSize * 2)
	if err != nil {
		log.Errorf("failed to generate session key: %s", err)
		return nil, err
	}

	tempCoder, err := crypt.NewCoder(sessionKey)
	if err != nil {
		log.Errorf("failed to create temp coder: %s", err)
		return nil, err
	}

	key := sha512.Sum512([]byte(password))

	encryptedKey := tempCoder.Encrypt(key[:])

	primaryCoder, err := crypt.NewCoder(key[:])
	if err != nil {
		log.Errorf("failed to create coder: %s", err)
		return nil, err
	}

	return &SessionManager{primaryCoder, encryptedKey, sessionKey}, nil
}

func NewSessionManagerFromEncryptedKey(encryptedKeyHex, sessionKeyHex string) (*SessionManager, error) {
	encryptedKey, err := hex.DecodeString(encryptedKeyHex)
	if err != nil {
		log.Errorf("failed to decode hex to encrypted key: %s", err)
		return nil, err
	}

	sessionKey, err := hex.DecodeString(sessionKeyHex)
	if err != nil {
		log.Errorf("failed to decode hex to session key: %s", err)
		return nil, err
	}

	tempCoder, err := crypt.NewCoder(sessionKey)
	if err != nil {
		log.Errorf("failed to create temp coder: %s", err)
		return nil, err
	}

	key, err := tempCoder.Decrypt(encryptedKey)
	if err != nil {
		log.Errorf("failed to decode encrypted key: %s", err)
		return nil, err
	}

	primaryCoder, err := crypt.NewCoder(key)
	if err != nil {
		log.Errorf("failed to create coder: %s", err)
		return nil, err
	}

	return &SessionManager{primaryCoder, encryptedKey, sessionKey}, nil
}

func (s *SessionManager) GetEncryptedKey() string {
	return hex.EncodeToString(s.encryptedKey)
}

func (s *SessionManager) GetSessionKey() string {
	return hex.EncodeToString(s.sessionKey)
}
