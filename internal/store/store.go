package store

import (
	"context"

	pb "github.com/VladBag2022/gokeeper/internal/proto"
)

type Store interface {
	IsUsernameAvailable(ctx context.Context, username string) (available bool, err error)
	SignIn(ctx context.Context, credentials *pb.Credentials) (id int64, err error)
	SignUp(ctx context.Context, credentials *pb.Credentials) (id int64, err error)

	StoreSecret(ctx context.Context, userID int64, secret *pb.Secret) (id int64, err error)
	UpdateSecret(ctx context.Context, id int64, secret *pb.Secret) error
	DeleteSecret(ctx context.Context, id int64) error

	StoreMeta(ctx context.Context, secretID int64, meta *pb.Meta) (id int64, err error)
	UpdateMeta(ctx context.Context, id int64, meta *pb.Meta) error
	DeleteMeta(ctx context.Context, id int64) error

	GetSecrets(ctx context.Context, userID int64) (secrets *pb.ClientSecrets, err error)

	IsUserSecret(ctx context.Context, userID, secretID int64) (userSecret bool, err error)
	IsUserMeta(ctx context.Context, userID, metaID int64) (userMeta bool, err error)
}
