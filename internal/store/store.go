// Package store contains store interfaces and implementations.
package store

import (
	"context"
)

// Credentials consists of username and password.
type Credentials struct {
	Username string
	Password string
}

// SecretKind represents the secret type.
type SecretKind int32

// Secret consists of its binary data and kind (type).
type Secret struct {
	Data []byte
	Kind SecretKind
}

// Meta represents ant text meta information.
type Meta string

// StoredMeta consists of Meta and its store ID.
type StoredMeta struct {
	Meta Meta
	ID   int64
}

// StoredSecret consists of Secret, bounded Meta and its store ID.
type StoredSecret struct {
	Secret Secret
	Meta   []StoredMeta
	ID     int64
}

// Store is the main store interface.
type Store interface {
	// IsUsernameAvailable checks whether provided username is available.
	IsUsernameAvailable(ctx context.Context, username string) (available bool, err error)

	// SignIn returns user ID from provided credentials.
	SignIn(ctx context.Context, credentials Credentials) (id int64, err error)

	// SignUp registers new user and returns his/her new ID.
	SignUp(ctx context.Context, credentials Credentials) (id int64, err error)

	// StoreSecret stores user secret and returns its new ID.
	StoreSecret(ctx context.Context, userID int64, secret Secret) (id int64, err error)

	// UpdateSecret updates secret by its ID.
	UpdateSecret(ctx context.Context, id int64, secret Secret) error

	// DeleteSecret deletes secret by its ID.
	DeleteSecret(ctx context.Context, id int64) error

	// StoreMeta stores secret meta and returns ins new ID.
	StoreMeta(ctx context.Context, secretID int64, meta Meta) (id int64, err error)

	// UpdateMeta updates meta by its ID.
	UpdateMeta(ctx context.Context, id int64, meta Meta) error

	// DeleteMeta deletes meta by its ID.
	DeleteMeta(ctx context.Context, id int64) error

	// GetSecrets returns user' secrets.
	GetSecrets(ctx context.Context, userID int64) (secrets []StoredSecret, err error)

	// GetEncryptedKey returns user's encrypted key.
	GetEncryptedKey(ctx context.Context, userID int64) (secret StoredSecret, err error)

	// IsUserSecret checks whether secret belongs to user.
	IsUserSecret(ctx context.Context, userID, secretID int64) (userSecret bool, err error)

	// IsUserMeta checks whether meta belongs to meta.
	IsUserMeta(ctx context.Context, userID, metaID int64) (userMeta bool, err error)
}
