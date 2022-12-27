package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
)

type Coder struct {
	gcm   cipher.AEAD
	nonce []byte
}

func NewCoder(key []byte) (*Coder, error) {
	s := sha512.Sum512(key)

	c, err := aes.NewCipher(s[:aes.BlockSize*2])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	return &Coder{
		gcm:   gcm,
		nonce: s[len(s)-gcm.NonceSize():],
	}, nil
}

func (c *Coder) Encrypt(plain []byte) []byte {
	return c.gcm.Seal(nil, c.nonce, plain, nil)
}

func (c *Coder) Decrypt(message []byte) ([]byte, error) {
	return c.gcm.Open(nil, c.nonce, message, nil)
}
