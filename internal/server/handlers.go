package server

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	pb "github.com/VladBag2022/gokeeper/internal/proto"
)

// StoreSecret stores provided secret and returns new ID.
func (s *KeeperServer) StoreSecret(ctx context.Context, secret *pb.Secret) (*pb.ClientSecret, error) {
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	secretID, err := s.store.StoreSecret(ctx, userID, secret)
	return &pb.ClientSecret{
		Id:     secretID,
		Secret: secret,
	}, err
}

// UpdateSecret checks user permissions and updates secret by ID.
func (s *KeeperServer) UpdateSecret(ctx context.Context, secret *pb.ClientSecret) (*empty.Empty, error) {
	if err := s.permitSecretID(ctx, secret.GetId()); err != nil {
		return nil, err
	}
	return &empty.Empty{}, s.store.UpdateSecret(ctx, secret.GetId(), secret.GetSecret())
}

// DeleteSecret checks user permissions and deletes secret by ID.
func (s *KeeperServer) DeleteSecret(ctx context.Context, secret *pb.ClientSecret) (*empty.Empty, error) {
	if err := s.permitSecretID(ctx, secret.GetId()); err != nil {
		return nil, err
	}
	return &empty.Empty{}, s.store.DeleteSecret(ctx, secret.GetId())
}

// StoreMeta stores provided meta and returns new ID.
func (s *KeeperServer) StoreMeta(ctx context.Context, req *pb.StoreMetaRequest) (*pb.ClientMeta, error) {
	if err := s.permitSecretID(ctx, req.GetSecretId()); err != nil {
		return nil, err
	}
	metaID, err := s.store.StoreMeta(ctx, req.GetSecretId(), req.GetMeta())
	if err != nil {
		return nil, err
	}
	return &pb.ClientMeta{
		Id:   metaID,
		Meta: req.GetMeta(),
	}, nil
}

// UpdateMeta checks user permissions and updates meta by ID.
func (s *KeeperServer) UpdateMeta(ctx context.Context, meta *pb.ClientMeta) (*empty.Empty, error) {
	if err := s.permitMetaID(ctx, meta.GetId()); err != nil {
		return nil, err
	}
	return &empty.Empty{}, s.store.UpdateMeta(ctx, meta.GetId(), meta.GetMeta())
}

// DeleteMeta checks user permissions and deletes meta by ID.
func (s *KeeperServer) DeleteMeta(ctx context.Context, meta *pb.ClientMeta) (*empty.Empty, error) {
	if err := s.permitMetaID(ctx, meta.GetId()); err != nil {
		return nil, err
	}
	return &empty.Empty{}, s.store.DeleteMeta(ctx, meta.GetId())
}

// GetSecrets returns user' secrets if any.
func (s *KeeperServer) GetSecrets(ctx context.Context, _ *empty.Empty) (*pb.ClientSecrets, error) {
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return s.store.GetSecrets(ctx, userID)
}

// GetEncryptedKey returns user's encrypted key if any.
func (s *KeeperServer) GetEncryptedKey(ctx context.Context, _ *empty.Empty) (*pb.ClientSecret, error) {
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	return s.store.GetEncryptedKey(ctx, userID)
}
