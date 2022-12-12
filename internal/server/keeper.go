package server

import (
	"context"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/VladBag2022/gokeeper/internal/proto"
	"github.com/VladBag2022/gokeeper/internal/storage"
)

type KeeperServer struct {
	pb.UnimplementedKeeperServer

	store storage.Repository
}

func NewKeeperServer(store storage.Repository) *KeeperServer {
	return &KeeperServer{store: store}
}

func (s *KeeperServer) permitSecretID(ctx context.Context, secretID int64) error {
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return err
	}
	userSecret, err := s.store.IsUserSecret(ctx, userID, secretID)
	if err != nil {
		return err
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
		return err
	}
	if !userMeta {
		return status.Errorf(codes.PermissionDenied, "IDOR attack attempt detected")
	}
	return nil
}

func userIDFromContext(ctx context.Context) (userID int64, err error) {
	var s string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get(UserIDKey)
		if len(values) > 0 {
			s = values[0]
		}
	}
	if len(s) == 0 {
		return 0, status.Error(codes.Unauthenticated, "missing userID")
	}
	return strconv.ParseInt(s, 10, 64)
}
