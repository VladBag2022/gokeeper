package store

import (
	"context"

	pb "github.com/VladBag2022/gokeeper/internal/proto"
)

// GRPCStore is the store interface suitable for gRPC services.
type GRPCStore interface {
	// IsUsernameAvailable checks whether provided username is available.
	IsUsernameAvailable(ctx context.Context, username string) (available bool, err error)

	// SignIn returns user ID from provided credentials.
	SignIn(ctx context.Context, credentials *pb.Credentials) (id int64, err error)

	// SignUp registers new user and returns his/her new ID.
	SignUp(ctx context.Context, credentials *pb.Credentials) (id int64, err error)

	// StoreSecret stores user secret and returns its new ID.
	StoreSecret(ctx context.Context, userID int64, secret *pb.Secret) (id int64, err error)

	// UpdateSecret updates secret by its ID.
	UpdateSecret(ctx context.Context, id int64, secret *pb.Secret) error

	// DeleteSecret deletes secret by its ID.
	DeleteSecret(ctx context.Context, id int64) error

	// StoreMeta stores secret meta and returns ins new ID.
	StoreMeta(ctx context.Context, secretID int64, meta *pb.Meta) (id int64, err error)

	// UpdateMeta updates meta by its ID.
	UpdateMeta(ctx context.Context, id int64, meta *pb.Meta) error

	// DeleteMeta deletes meta by its ID.
	DeleteMeta(ctx context.Context, id int64) error

	// GetSecrets returns client' secrets.
	GetSecrets(ctx context.Context, userID int64) (secrets *pb.ClientSecrets, err error)

	// GetEncryptedKey returns client's encrypted key.
	GetEncryptedKey(ctx context.Context, userID int64) (secret *pb.ClientSecret, err error)

	// IsUserSecret checks whether secret belongs to user.
	IsUserSecret(ctx context.Context, userID, secretID int64) (userSecret bool, err error)

	// IsUserMeta checks whether meta belongs to meta.
	IsUserMeta(ctx context.Context, userID, metaID int64) (userMeta bool, err error)
}
