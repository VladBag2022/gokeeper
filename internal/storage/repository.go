package storage

import (
	"context"

	pb "github.com/VladBag2022/gokeeper/internal/proto"
)

type Repository interface {
	IsUsernameAvailable(ctx context.Context, username string) (available bool, err error)
	SignIn(ctx context.Context, credentials *pb.Credentials) (id int64, err error)
	SignUp(ctx context.Context, credentials *pb.Credentials) (id int64, err error)

	StoreSecret(ctx context.Context, userID int64, secret *pb.Secret) (secretID int64, err error)
	UpdateSecret(ctx context.Context, secretID int64, secret *pb.Secret) error
	DeleteSecret(ctx context.Context, secretID int64) error

	StoreMetaInfo(ctx context.Context, metaInfo *pb.MetaInfo) (id int64, err error)
	UpdateMetaInfo(ctx context.Context, metaInfoID int64, metaInfo *pb.MetaInfo) error
	DeleteMetaInfo(ctx context.Context, metaInfoID int64) error

	LinkMetaInfo(ctx context.Context, metaInfoID, secretID int64) error
	UnlinkMetaInfo(ctx context.Context, metaInfoID, secretID int64) error
	CountMetaInfoLinks(ctx context.Context, metaInfoID int64) (counter int64, err error)
}