package server

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/VladBag2022/gokeeper/internal/jwt"
	"github.com/VladBag2022/gokeeper/internal/mocks"
	pb "github.com/VladBag2022/gokeeper/internal/proto"
)

func newTestJWTManager() *jwt.Manager {
	return jwt.NewManager([]byte("secret"), time.Hour)
}

func TestAuthServer_SignIn(t *testing.T) {
	tests := []struct {
		name       string
		req        *pb.Credentials
		authorized bool
	}{
		{
			name: "positive test",
			req: &pb.Credentials{
				Username: "user1",
				Password: "pass1",
			},
			authorized: false,
		},
		{
			name: "negative test",
			req: &pb.Credentials{
				Username: "user1",
				Password: "pass1",
			},
			authorized: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := mocks.NewStore(t)
			if tt.authorized {
				s.On("SignIn", mock.Anything, tt.req).Return(int64(0), errors.New("wrong credentials"))
			} else {
				s.On("SignIn", mock.Anything, tt.req).Return(int64(0), nil)
			}

			as := &AuthServer{
				store:      s,
				jwtManager: newTestJWTManager(),
			}
			got, err := as.SignIn(context.Background(), tt.req)
			if (err != nil) != tt.authorized {
				t.Errorf("SignIn() error = %v, authorized %v", err, tt.authorized)
				return
			}
			if err == nil {
				assert.NotEmpty(t, got.GetToken())
			}
		})
	}
}

func TestAuthServer_SignUp(t *testing.T) {
	tests := []struct {
		name      string
		req       *pb.Credentials
		available bool
	}{
		{
			name: "positive test",
			req: &pb.Credentials{
				Username: "user1",
				Password: "pass1",
			},
			available: true,
		},
		{
			name: "negative test",
			req: &pb.Credentials{
				Username: "user1",
				Password: "pass1",
			},
			available: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := mocks.NewStore(t)
			s.On("IsUsernameAvailable", mock.Anything, tt.req.GetUsername()).Return(tt.available, nil)
			if tt.available {
				s.On("SignUp", mock.Anything, tt.req).Return(int64(0), nil)
			}

			as := &AuthServer{
				store:      s,
				jwtManager: newTestJWTManager(),
			}
			got, err := as.SignUp(context.Background(), tt.req)
			if (err != nil) == tt.available {
				t.Errorf("SignIn() error = %v, available %v", err, tt.available)
				return
			}
			if err == nil {
				assert.NotEmpty(t, got.GetToken())
			}
		})
	}
}
