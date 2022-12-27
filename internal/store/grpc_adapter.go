package store

import (
	"context"

	pb "github.com/VladBag2022/gokeeper/internal/proto"
)

type GRPCAdapter struct {
	store Store
}

func NewGRPCAdapter(store Store) *GRPCAdapter {
	return &GRPCAdapter{
		store: store,
	}
}

func (a *GRPCAdapter) IsUsernameAvailable(ctx context.Context, username string) (available bool, err error) {
	return a.store.IsUsernameAvailable(ctx, username)
}

func (a *GRPCAdapter) SignIn(ctx context.Context, credentials *pb.Credentials) (id int64, err error) {
	return a.store.SignIn(ctx, fromCredentialsGRPC(credentials))
}

func (a *GRPCAdapter) SignUp(ctx context.Context, credentials *pb.Credentials) (id int64, err error) {
	return a.store.SignUp(ctx, fromCredentialsGRPC(credentials))
}

func (a *GRPCAdapter) StoreSecret(ctx context.Context, userID int64, secret *pb.Secret) (id int64, err error) {
	return a.store.StoreSecret(ctx, userID, fromSecretGRPC(secret))
}

func (a *GRPCAdapter) UpdateSecret(ctx context.Context, id int64, secret *pb.Secret) error {
	return a.store.UpdateSecret(ctx, id, fromSecretGRPC(secret))
}

func (a *GRPCAdapter) DeleteSecret(ctx context.Context, id int64) error {
	return a.store.DeleteSecret(ctx, id)
}

func (a *GRPCAdapter) StoreMeta(ctx context.Context, secretID int64, meta *pb.Meta) (id int64, err error) {
	return a.store.StoreMeta(ctx, secretID, fromMetaGRPC(meta))
}

func (a *GRPCAdapter) UpdateMeta(ctx context.Context, id int64, meta *pb.Meta) error {
	return a.store.UpdateMeta(ctx, id, fromMetaGRPC(meta))
}

func (a *GRPCAdapter) DeleteMeta(ctx context.Context, id int64) error {
	return a.store.DeleteMeta(ctx, id)
}

func (a *GRPCAdapter) GetSecrets(ctx context.Context, userID int64) (secrets *pb.ClientSecrets, err error) {
	secrets = &pb.ClientSecrets{}

	ss, err := a.store.GetSecrets(ctx, userID)
	if err != nil {
		return secrets, err
	}

	secrets = toClientSecretsGRPC(ss)
	return secrets, err
}

func (a *GRPCAdapter) IsUserSecret(ctx context.Context, userID, secretID int64) (userSecret bool, err error) {
	return a.store.IsUserSecret(ctx, userID, secretID)
}

func (a *GRPCAdapter) IsUserMeta(ctx context.Context, userID, metaID int64) (userMeta bool, err error) {
	return a.store.IsUserMeta(ctx, userID, metaID)
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
