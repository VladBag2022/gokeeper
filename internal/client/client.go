package client

import (
	"google.golang.org/grpc"

	pb "github.com/VladBag2022/gokeeper/internal/proto"
)

type Client struct {
	Auth   pb.AuthClient
	Keeper pb.KeeperClient

	SessionKey []byte
}

func NewClient(cc *grpc.ClientConn, sessionKey []byte) *Client {
	return &Client{pb.NewAuthClient(cc), pb.NewKeeperClient(cc), sessionKey}
}
