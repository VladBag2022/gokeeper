package store

import (
	"context"
	"fmt"

	pb "github.com/VladBag2022/gokeeper/internal/proto"
)

// GRPCAdapter is the gRPC adapter for main Store.
type GRPCAdapter struct {
	store Store
}

// NewGRPCAdapter returns new GRPCAdapter.
func NewGRPCAdapter(store Store) *GRPCAdapter {
	return &GRPCAdapter{
		store: store,
	}
}

// IsUsernameAvailable checks whether provided username is available.
func (a *GRPCAdapter) IsUsernameAvailable(ctx context.Context, username string) (available bool, err error) {
	available, err = a.store.IsUsernameAvailable(ctx, username)
	if err != nil {
		return available, fmt.Errorf("failed to call IsUsernameAvailable from gRPC store adapter: %s", err)
	}

	return available, nil
}

// SignIn returns user ID from provided credentials.
func (a *GRPCAdapter) SignIn(ctx context.Context, credentials *pb.Credentials) (id int64, err error) {
	id, err = a.store.SignIn(ctx, fromCredentialsGRPC(credentials))
	if err != nil {
		return 0, fmt.Errorf("failed to call SignIn from gRPC store adapter: %s", err)
	}

	return id, nil
}

// SignUp registers new user and returns his/her new ID.
func (a *GRPCAdapter) SignUp(ctx context.Context, credentials *pb.Credentials) (id int64, err error) {
	id, err = a.store.SignUp(ctx, fromCredentialsGRPC(credentials))
	if err != nil {
		return 0, fmt.Errorf("failed to call SignUp from gRPC store adapter: %s", err)
	}

	return id, nil
}

// StoreSecret stores user secret and returns its new ID.
func (a *GRPCAdapter) StoreSecret(ctx context.Context, userID int64, secret *pb.Secret) (id int64, err error) {
	id, err = a.store.StoreSecret(ctx, userID, fromSecretGRPC(secret))
	if err != nil {
		return 0, fmt.Errorf("failed to call StoreSecret from gRPC store adapter: %s", err)
	}

	return id, nil
}

// UpdateSecret updates secret by its ID.
func (a *GRPCAdapter) UpdateSecret(ctx context.Context, id int64, secret *pb.Secret) error {
	if err := a.store.UpdateSecret(ctx, id, fromSecretGRPC(secret)); err != nil {
		return fmt.Errorf("failed to call UpdateSecret from gRPC store adapter: %s", err)
	}

	return nil
}

// DeleteSecret deletes secret by its ID.
func (a *GRPCAdapter) DeleteSecret(ctx context.Context, id int64) error {
	if err := a.store.DeleteSecret(ctx, id); err != nil {
		return fmt.Errorf("failed to call DeleteSecret from gRPC store adapter: %s", err)
	}

	return nil
}

// StoreMeta stores secret meta and returns ins new ID.
func (a *GRPCAdapter) StoreMeta(ctx context.Context, secretID int64, meta *pb.Meta) (id int64, err error) {
	id, err = a.store.StoreMeta(ctx, secretID, fromMetaGRPC(meta))
	if err != nil {
		return 0, fmt.Errorf("failed to call StoreMeta from gRPC store adapter: %s", err)
	}

	return id, nil
}

// UpdateMeta updates meta by its ID.
func (a *GRPCAdapter) UpdateMeta(ctx context.Context, id int64, meta *pb.Meta) error {
	if err := a.store.UpdateMeta(ctx, id, fromMetaGRPC(meta)); err != nil {
		return fmt.Errorf("failed to call UpdateMeta from gRPC store adapter: %s", err)
	}

	return nil
}

// DeleteMeta deletes meta by its ID.
func (a *GRPCAdapter) DeleteMeta(ctx context.Context, id int64) error {
	if err := a.store.DeleteMeta(ctx, id); err != nil {
		return fmt.Errorf("failed to call DeleteMeta from gRPC store adapter: %s", err)
	}

	return nil
}

// GetSecrets returns client' secrets.
func (a *GRPCAdapter) GetSecrets(ctx context.Context, userID int64) (secrets *pb.ClientSecrets, err error) {
	ss, err := a.store.GetSecrets(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetSecrets from gRPC store adapter: %s", err)
	}

	return toClientSecretsGRPC(ss), nil
}

// GetEncryptedKey returns client's encrypted key.
func (a *GRPCAdapter) GetEncryptedKey(ctx context.Context, userID int64) (secret *pb.ClientSecret, err error) {
	k, err := a.store.GetEncryptedKey(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetEncryptedKey from gRPC store adapter: %s", err)
	}

	return toClientSecretGRPC(k), nil
}

// IsUserSecret checks whether secret belongs to user.
func (a *GRPCAdapter) IsUserSecret(ctx context.Context, userID, secretID int64) (userSecret bool, err error) {
	userSecret, err = a.store.IsUserSecret(ctx, userID, secretID)
	if err != nil {
		err = fmt.Errorf("failed to call IsUserSecret from gRPC store adapter: %s", err)
	}

	return userSecret, err
}

// IsUserMeta checks whether meta belongs to meta.
func (a *GRPCAdapter) IsUserMeta(ctx context.Context, userID, metaID int64) (userMeta bool, err error) {
	userMeta, err = a.store.IsUserMeta(ctx, userID, metaID)
	if err != nil {
		err = fmt.Errorf("failed to call IsUserMeta from gRPC store adapter: %s", err)
	}

	return userMeta, err
}

func fromCredentialsGRPC(credentials *pb.Credentials) Credentials {
	return Credentials{
		Username: credentials.GetUsername(),
		Password: credentials.GetPassword(),
	}
}

func fromSecretGRPC(secret *pb.Secret) Secret {
	return Secret{
		Data: secret.GetData(),
		Kind: SecretKind(secret.GetKind()),
	}
}

func fromMetaGRPC(meta *pb.Meta) Meta {
	return Meta(meta.GetText())
}

func toClientSecretGRPC(secret ClientSecret) (pbSecret *pb.ClientSecret) {
	return &pb.ClientSecret{
		Secret: &pb.Secret{
			Data: secret.Secret.Data,
			Kind: pb.SecretKind(secret.Secret.Kind),
		},
		Id: secret.ID,
	}
}

func toClientMetaGRPC(meta ClientMeta) (pbMeta *pb.ClientMeta) {
	return &pb.ClientMeta{
		Meta: &pb.Meta{
			Text: string(meta.Meta),
		},
		Id: meta.ID,
	}
}

func toClientMetasGRPC(metas []ClientMeta) (pbMetas []*pb.ClientMeta) {
	pbMetas = make([]*pb.ClientMeta, len(metas))

	for i, m := range metas {
		pbMetas[i] = toClientMetaGRPC(m)
	}

	return pbMetas
}

func toClientSecretsGRPC(secrets []ClientSecret) (pbSecrets *pb.ClientSecrets) {
	pbSecrets = &pb.ClientSecrets{}

	pbSecrets.Secrets = make([]*pb.ClientSecret, len(secrets))
	for i, s := range secrets {
		pbSecrets.Secrets[i] = toClientSecretGRPC(s)
		pbSecrets.Secrets[i].Meta = toClientMetasGRPC(s.Meta)
	}

	return pbSecrets
}
