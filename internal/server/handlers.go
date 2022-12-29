package server

import (
	"context"
	"fmt"

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
	if err != nil {
		return nil, fmt.Errorf("failed to store secret: %s", err)
	}

	return &pb.ClientSecret{
		Id:     secretID,
		Secret: secret,
	}, nil
}

// UpdateSecret checks user permissions and updates secret by ID.
func (s *KeeperServer) UpdateSecret(ctx context.Context, secret *pb.ClientSecret) (*empty.Empty, error) {
	if err := s.permitSecretID(ctx, secret.GetId()); err != nil {
		return nil, err
	}

	if err := s.store.UpdateSecret(ctx, secret.GetId(), secret.GetSecret()); err != nil {
		return nil, fmt.Errorf("failed to update secret in store: %s", err)
	}

	return &empty.Empty{}, nil
}

// DeleteSecret checks user permissions and deletes secret by ID.
func (s *KeeperServer) DeleteSecret(ctx context.Context, secret *pb.ClientSecret) (*empty.Empty, error) {
	if err := s.permitSecretID(ctx, secret.GetId()); err != nil {
		return nil, err
	}

	if err := s.store.DeleteSecret(ctx, secret.GetId()); err != nil {
		return nil, fmt.Errorf("failed to delete secret from store: %s", err)
	}

	return &empty.Empty{}, nil
}

// StoreMeta stores provided meta and returns new ID.
func (s *KeeperServer) StoreMeta(ctx context.Context, req *pb.StoreMetaRequest) (*pb.ClientMeta, error) {
	if err := s.permitSecretID(ctx, req.GetSecretId()); err != nil {
		return nil, err
	}

	metaID, err := s.store.StoreMeta(ctx, req.GetSecretId(), req.GetMeta())
	if err != nil {
		return nil, fmt.Errorf("failed to store meta: %s", err)
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

	if err := s.store.UpdateMeta(ctx, meta.GetId(), meta.GetMeta()); err != nil {
		return nil, fmt.Errorf("failed to update meta in store: %s", err)
	}

	return &empty.Empty{}, nil
}

// DeleteMeta checks user permissions and deletes meta by ID.
func (s *KeeperServer) DeleteMeta(ctx context.Context, meta *pb.ClientMeta) (*empty.Empty, error) {
	if err := s.permitMetaID(ctx, meta.GetId()); err != nil {
		return nil, err
	}

	if err := s.store.DeleteMeta(ctx, meta.GetId()); err != nil {
		return nil, fmt.Errorf("failed to delete meta from store: %s", err)
	}

	return &empty.Empty{}, nil
}

// GetSecrets returns user' secrets if any.
func (s *KeeperServer) GetSecrets(ctx context.Context, _ *empty.Empty) (*pb.ClientSecrets, error) {
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	secrets, err := s.store.GetSecrets(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user' secrets from store: %s", err)
	}

	return secrets, nil
}

// GetEncryptedKey returns user's encrypted key if any.
func (s *KeeperServer) GetEncryptedKey(ctx context.Context, _ *empty.Empty) (*pb.ClientSecret, error) {
	userID, err := userIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	secret, err := s.store.GetEncryptedKey(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user's encrypted key from store: %s", err)
	}

	return secret, nil
}
