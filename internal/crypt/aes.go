package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	"fmt"
)

// Coder is the AES GCM crypter/decryoter.
type Coder struct {
	gcm   cipher.AEAD
	nonce []byte
}

// NewCoder returns new Coder from provided secret key.
func NewCoder(key []byte) (*Coder, error) {
	hash := sha512.Sum512(key)

	c, err := aes.NewCipher(hash[:aes.BlockSize*2])
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	return &Coder{
		gcm:   gcm,
		nonce: hash[len(hash)-gcm.NonceSize():],
	}, nil
}

// Encrypt encrypts provided buffer.
func (c *Coder) Encrypt(plain []byte) []byte {
	return c.gcm.Seal(nil, c.nonce, plain, nil)
}

// Decrypt decrypts provided buffer.
func (c *Coder) Decrypt(message []byte) ([]byte, error) {
	plain, err := c.gcm.Open(nil, c.nonce, message, nil)
	if err != nil {
		return plain, fmt.Errorf("failed to decrypt ciphertext: %w", err)
	}

	return plain, nil
}
