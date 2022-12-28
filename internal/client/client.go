package client

import (
	"google.golang.org/grpc"

	pb "github.com/VladBag2022/gokeeper/internal/proto"
)

type Client struct {
	Auth   pb.AuthClient
	Keeper pb.KeeperClient
}

func NewClient(cc *grpc.ClientConn) *Client {
	return &Client{pb.NewAuthClient(cc), pb.NewKeeperClient(cc)}
}
