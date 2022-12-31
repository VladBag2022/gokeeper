package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/VladBag2022/gokeeper/internal/crypt"
)

func TestManager_Generate(t *testing.T) {
	tests := []struct {
		name   string
		userID int64
	}{
		{
			name:   "positive test",
			userID: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, err := crypt.GenerateRandomBytes(10)
			require.NoError(t, err)

			m := &Manager{
				secretKey:     k,
				tokenDuration: time.Minute,
			}

			got, err := m.Generate(tt.userID)
			require.NoError(t, err)
			assert.NotEmpty(t, got)
		})
	}
}

func TestManager_Verify(t *testing.T) {
	type fields struct {
		secretKey     []byte
		tokenDuration time.Duration
	}

	tests := []struct {
		name        string
		fields      fields
		accessToken string
		wantUserID  int64
		wantErr     bool
	}{
		{
			name: "positive test - 100 years duration",
			fields: fields{
				secretKey:     []byte("secret"),
				tokenDuration: time.Hour * 24 * 365 * 100,
			},
			accessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQ4MjQ2MDI1NjcsInVzZXJfaWQiOjF9." +
				"S6BxLPp_V5pRPg3-7BLf86RVtllqJHFb2S3CTLYRJt8",
			wantUserID: 0,
			wantErr:    false,
		},
		{
			name: "negative test - invalid token",
			fields: fields{
				secretKey:     []byte("secret"),
				tokenDuration: time.Hour * 24 * 365 * 100,
			},
			accessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjQ4MjQ2MDI1NjcsInVzZXJfaWQiOjF9." +
				"S6BxLPp_V5pRPg3-7BLf86RVtllqJHFb2S3CTLYRJt8xxxxxxxxxxxxx",
			wantUserID: 0,
			wantErr:    true,
		},
		{
			name: "negative test - userID 0, token expired",
			fields: fields{
				secretKey:     []byte("secret"),
				tokenDuration: time.Minute,
			},
			accessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzEwMDE4NzAsInVzZXJfaWQiOjB9." +
				"iWAB9nrD4FWcLP4gk5cVUh6AFJqf3SoFnYoVENw2J0M",
			wantUserID: 0,
			wantErr:    true,
		},
		{
			name: "negative test - userID 1, token expired",
			fields: fields{
				secretKey:     []byte("secret"),
				tokenDuration: time.Minute,
			},
			accessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzEwMDE5MzEsInVzZXJfaWQiOjF9." +
				"3hS-Uzg4nfGgKMV8SQU8CTt2rv2bC-clxc40yNREaXE",
			wantUserID: 1,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				secretKey:     tt.fields.secretKey,
				tokenDuration: tt.fields.tokenDuration,
			}
			got, err := m.Verify(tt.accessToken)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantUserID, got.UserID)
			}
		})
	}
}
