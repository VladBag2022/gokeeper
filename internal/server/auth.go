package server

import (
	"github.com/VladBag2022/gokeeper/internal/jwt"
	pb "github.com/VladBag2022/gokeeper/internal/proto"
	"github.com/VladBag2022/gokeeper/internal/storage"
)

type AuthServer struct {
	pb.UnimplementedAuthServer

	store      storage.Repository
	jwtManager *jwt.Manager
}

func NewAuthServer(store storage.Repository, jwtManager *jwt.Manager) *AuthServer {
	return &AuthServer{store: store, jwtManager: jwtManager}
}
