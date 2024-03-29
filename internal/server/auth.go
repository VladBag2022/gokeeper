package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/VladBag2022/gokeeper/internal/jwt"
	pb "github.com/VladBag2022/gokeeper/internal/proto"
	"github.com/VladBag2022/gokeeper/internal/store"
)

// AuthServer implements gRPC Auth service.
type AuthServer struct {
	pb.UnimplementedAuthServer

	store      store.GRPCStore
	jwtManager *jwt.Manager
}

// NewAuthServer returns new AuthServer.
func NewAuthServer(store store.GRPCStore, jwtManager *jwt.Manager) *AuthServer {
	return &AuthServer{store: store, jwtManager: jwtManager}
}

// SignIn checks user credentials and returns JWT.
func (s *AuthServer) SignIn(ctx context.Context, req *pb.Credentials) (*pb.JWT, error) {
	id, err := s.store.SignIn(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "wrong credentials: %s", err)
	}

	token, err := s.jwtManager.Generate(id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate access token")
	}

	return &pb.JWT{Token: token}, nil
}

// SignUp registers new user and returns JWT.
func (s *AuthServer) SignUp(ctx context.Context, req *pb.Credentials) (*pb.JWT, error) {
	available, err := s.store.IsUsernameAvailable(ctx, req.GetUsername())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check whether the username is available")
	}

	if !available {
		return nil, status.Errorf(codes.AlreadyExists, "the username is not available")
	}

	id, err := s.store.SignUp(ctx, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create new user")
	}

	token, err := s.jwtManager.Generate(id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate access token")
	}

	return &pb.JWT{Token: token}, nil
}
