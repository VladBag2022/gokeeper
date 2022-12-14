package client

import (
	"context"
	"strings"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/VladBag2022/gokeeper/internal/server"
)

type AuthInterceptor struct {
	jwt string
}

func NewAuthInterceptor(jwt string) *AuthInterceptor {
	return &AuthInterceptor{jwt}
}

func (i *AuthInterceptor) Unary() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		log.Infof("call %s", method)
		if strings.Contains(method, "Sign") {
			log.Debug("not attaching JWT during call")
			return invoker(ctx, method, req, reply, cc, opts...)
		}
		log.Debug("attaching JWT during call")
		return invoker(i.attachToken(ctx), method, req, reply, cc, opts...)
	}
}

func (i *AuthInterceptor) attachToken(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, server.AuthorizationKey, i.jwt)
}
