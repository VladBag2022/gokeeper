package store

import (
	"context"
)

type Credentials struct {
	Username string
	Password string
}

type SecretKind int32

type Secret struct {
	Data []byte
	Kind SecretKind
}

type Meta string

type ClientMeta struct {
	Meta Meta
	ID   int64
}

type ClientSecret struct {
	Secret Secret
	Meta   []ClientMeta
	ID     int64
}

type Store interface {
	IsUsernameAvailable(ctx context.Context, username string) (available bool, err error)
	SignIn(ctx context.Context, credentials Credentials) (id int64, err error)
	SignUp(ctx context.Context, credentials Credentials) (id int64, err error)

	StoreSecret(ctx context.Context, userID int64, secret Secret) (id int64, err error)
	UpdateSecret(ctx context.Context, id int64, secret Secret) error
	DeleteSecret(ctx context.Context, id int64) error

	StoreMeta(ctx context.Context, secretID int64, meta Meta) (id int64, err error)
	UpdateMeta(ctx context.Context, id int64, meta Meta) error
	DeleteMeta(ctx context.Context, id int64) error

	GetSecrets(ctx context.Context, userID int64) (secrets []ClientSecret, err error)
	GetEncryptedKey(ctx context.Context, userID int64) (secret ClientSecret, err error)

	IsUserSecret(ctx context.Context, userID, secretID int64) (userSecret bool, err error)
	IsUserMeta(ctx context.Context, userID, metaID int64) (userMeta bool, err error)
}
