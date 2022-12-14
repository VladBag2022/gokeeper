package server

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthInterceptor_authorize(t *testing.T) {
	tests := []struct {
		name    string
		method  string
		wantErr bool
	}{
		{
			name:    "positive test - SignIn",
			method:  "SignIn",
			wantErr: false,
		},
		{
			name:    "positive test - SignUp",
			method:  "SignUp",
			wantErr: false,
		},
		{
			name:    "negative test - GetSecrets",
			method:  "GetSecrets",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &AuthInterceptor{
				jwtManager: newTestJWTManager(),
			}
			_, err := i.authorize(context.Background(), tt.method)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
