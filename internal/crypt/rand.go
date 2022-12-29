// Package crypt contains security functions.
package crypt

import (
	"crypto/rand"
	"fmt"
)

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	buf := make([]byte, n)
	
	// Note that err == nil only if we read len(buf) bytes.
	if _, err := rand.Read(buf); err != nil {
		return nil, fmt.Errorf("failed to read %d random byte(s): %s", n, err)
	}

	return buf, nil
}
