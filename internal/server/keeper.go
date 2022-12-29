// Package server contains the GoKeeper gRPC server.
package server

import (
	"context"
	"fmt"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/VladBag2022/gokeeper/internal/proto"
	"github.com/VladBag2022/gokeeper/internal/store"
)

// KeeperServer implements gRPC Keeper service.
type KeeperServer struct {
	pb.UnimplementedKeeperServer

	store store.GRPCStore
}

// NewKeeperServer returns new KeeperServer.
func NewKeeperServer(store store.GRPCStore) *KeeperServer {
	return &KeeperServer{store: store}
}

func (s *KeeperServer) permitSecretID(ctx context.Context, secretID int64) error {
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return err
	}

	userSecret, err := s.store.IsUserSecret(ctx, userID, secretID)
	if err != nil {
		return fmt.Errorf("failed to check wether secret belongs to user: %s", err)
	}

	if !userSecret {
		return status.Errorf(codes.PermissionDenied, "IDOR attack attempt detected")
	}

	return nil
}

func (s *KeeperServer) permitMetaID(ctx context.Context, metaID int64) error {
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return err
	}

	userMeta, err := s.store.IsUserMeta(ctx, userID, metaID)
	if err != nil {
		return fmt.Errorf("failed to check wether meta belongs to user: %s", err)
	}

	if !userMeta {
		return status.Errorf(codes.PermissionDenied, "IDOR attack attempt detected")
	}

	return nil
}

func userIDFromContext(ctx context.Context) (userID int64, err error) {
	var userIDStr string

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get(UserIDKey)
		if len(values) > 0 {
			userIDStr = values[0]
		}
	}

	if len(userIDStr) == 0 {
		return 0, fmt.Errorf("failed to get userID from context: %s",
			status.Error(codes.Unauthenticated, "missing userID"))
	}

	userID, err = strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse userID from context: %s", err)
	}

	return userID, nil
}
