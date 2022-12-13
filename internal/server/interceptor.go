package server

import (
	"context"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/VladBag2022/gokeeper/internal/jwt"
)

const (
	AuthorizationKey = "authorization"
	UserIDKey        = "user_id"
)

type AuthInterceptor struct {
	jwtManager *jwt.Manager
}

func NewAuthInterceptor(jwtManager *jwt.Manager) *AuthInterceptor {
	return &AuthInterceptor{jwtManager}
}

func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		userID, err := i.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}
		metadata.NewIncomingContext(ctx, metadata.Pairs(UserIDKey, strconv.FormatInt(userID, 10)))
		return handler(ctx, req)
	}
}

func (i *AuthInterceptor) authorize(ctx context.Context, method string) (userID int64, err error) {
	if strings.Contains(method, "Sign") {
		return 0, nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md[AuthorizationKey]
	if len(values) == 0 {
		return 0, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]
	claims, err := i.jwtManager.Verify(accessToken)
	if err != nil {
		return 0, status.Errorf(codes.Unauthenticated, "access token is invalid: %s", err)
	}

	return claims.UserID, nil
}
