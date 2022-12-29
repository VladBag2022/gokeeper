// Package jwt is responsible for JWT generation and verifying.
package jwt

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Manager encapsulates JWT generation and verifying logic.
type Manager struct {
	secretKey     []byte
	tokenDuration time.Duration
}

// UserClaims are the extended JST claims containing user ID.
type UserClaims struct {
	jwt.StandardClaims

	UserID int64 `json:"user_id"`
}

// NewManager creates new Manager from secret key and token duration.
func NewManager(secretKey []byte, tokenDuration time.Duration) *Manager {
	return &Manager{secretKey, tokenDuration}
}

// Generate returns new valid JWT from provided user ID.
func (m *Manager) Generate(userID int64) (string, error) {
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(m.tokenDuration).Unix(),
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(m.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to create signed string from token: %s", err)
	}

	return tokenString, nil
}

// Verify checks JWT validity and returns encapsulated UserClaims.
func (m *Manager) Verify(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("unexpected token signing method")
			}

			return m.secretKey, nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
