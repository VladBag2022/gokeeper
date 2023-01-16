// Package client contains the gRPC client for GoKeeper gRPC service.
package client

import (
	"google.golang.org/grpc"

	pb "github.com/VladBag2022/gokeeper/internal/proto"
)

// Client is the GoKeeper gRPC client.
type Client struct {
	Auth   pb.AuthClient
	Keeper pb.KeeperClient
}

// NewClient creates new Client from gRPC connection.
func NewClient(cc *grpc.ClientConn) *Client {
	return &Client{pb.NewAuthClient(cc), pb.NewKeeperClient(cc)}
}
