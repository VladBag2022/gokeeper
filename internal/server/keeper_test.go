package server

import (
	"context"
	"github.com/stretchr/testify/mock"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"

	"github.com/VladBag2022/gokeeper/internal/mocks"
)

func newTestCtxWithUserID(userID string) context.Context {
	return metadata.NewIncomingContext(context.Background(), metadata.Pairs(UserIDKey, userID))
}

func TestKeeperServer_permitSecretID(t *testing.T) {
	tests := []struct {
		name    string
		ctxUserID  string
		secretID int64
		isUserSecret bool
		wantErr bool
	}{
		{
			name: "positive test",
			ctxUserID: "100",
			secretID: 100,
			isUserSecret: true,
			wantErr: false,
		},
		{
			name: "negative test - not user secret",
			ctxUserID: "100",
			secretID: 100,
			isUserSecret: false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := strconv.ParseInt(tt.ctxUserID, 10, 64)
			require.NoError(t, err)

			s := mocks.NewStore(t)
			s.On("IsUserSecret", mock.Anything, userID, tt.secretID).Return(tt.isUserSecret, nil)

			ks := &KeeperServer{store:s}
			ctx := newTestCtxWithUserID(tt.ctxUserID)
			err = ks.permitSecretID(ctx, tt.secretID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestKeeperServer_permitMetaID(t *testing.T) {
	tests := []struct {
		name    string
		ctxUserID  string
		metaID int64
		isUserMeta bool
		wantErr bool
	}{
		{
			name: "positive test",
			ctxUserID: "100",
			metaID: 100,
			isUserMeta: true,
			wantErr: false,
		},
		{
			name: "negative test - not user meta",
			ctxUserID: "100",
			metaID: 100,
			isUserMeta: false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := strconv.ParseInt(tt.ctxUserID, 10, 64)
			require.NoError(t, err)

			s := mocks.NewStore(t)
			s.On("IsUserMeta", mock.Anything, userID, tt.metaID).Return(tt.isUserMeta, nil)

			ks := &KeeperServer{store:s}
			ctx := newTestCtxWithUserID(tt.ctxUserID)
			err = ks.permitMetaID(ctx, tt.metaID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_userIDFromContext(t *testing.T) {
	tests := []struct {
		name       string
		setUserID bool
		ctxUserID  string
		wantUserID int64
		wantErr    bool
	}{
		{
			name:       "positive test",
			setUserID:  true,
			ctxUserID:  "100",
			wantUserID: 100,
			wantErr:    false,
		},
		{
			name:       "negative test - malformed userID",
			setUserID:  true,
			ctxUserID:  "1f00",
			wantUserID: 0,
			wantErr:    true,
		},
		{
			name:       "negative test - empty userID",
			setUserID:  true,
			ctxUserID:  "",
			wantUserID: 0,
			wantErr:    true,
		},
		{
			name:       "negative test - userID not set",
			setUserID:  false,
			ctxUserID:  "",
			wantUserID: 0,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.setUserID {
				ctx = newTestCtxWithUserID(tt.ctxUserID)
			}
			gotUserID, err := userIDFromContext(ctx)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantUserID, gotUserID)
			}
		})
	}
}
