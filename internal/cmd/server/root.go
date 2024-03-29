// Package server contains root command for GoKeeper server CLI.
package server

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	common "github.com/VladBag2022/gokeeper/internal/cmd"
	"github.com/VladBag2022/gokeeper/internal/crypt"
	"github.com/VladBag2022/gokeeper/internal/jwt"
	pb "github.com/VladBag2022/gokeeper/internal/proto"
	"github.com/VladBag2022/gokeeper/internal/server"
	"github.com/VladBag2022/gokeeper/internal/store"
)

// Execute starts GoKeeper server CLI.
func (c *CLI) Execute() {
	if err := c.cmd.Execute(); err != nil {
		log.Errorf("failed to execute root command: %s", err)
	}
}

func (c *CLI) run(_ *cobra.Command, _ []string) {
	ctx := context.Background()

	listener, err := net.Listen("tcp", viper.GetString("ListenAddress"))
	if err != nil {
		log.Errorf("failed to listen: %s", err)

		return
	}

	jwtKey, err := crypt.GenerateRandomBytes(32)
	if err != nil {
		log.Errorf("failed to generate random JWT key: %s", err)

		return
	}

	jwtManager := jwt.NewManager(jwtKey, viper.GetDuration("JWTDuration"))

	postgresStore, err := store.NewPostgresStore(ctx, viper.GetString("DatabaseDSN"))
	if err != nil {
		log.Errorf("failed to connect to Postgres store: %s", err)

		return
	}

	defer func() {
		if err = postgresStore.Close(); err != nil {
			log.Errorf("failed to close store: %s", err)
		}
	}()

	authInterceptor := server.NewAuthInterceptor(jwtManager)
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
	)

	grpcAdapter := store.NewGRPCAdapter(postgresStore)

	authServer := server.NewAuthServer(grpcAdapter, jwtManager)
	keeperServer := server.NewKeeperServer(grpcAdapter)

	pb.RegisterAuthServer(grpcServer, authServer)
	pb.RegisterKeeperServer(grpcServer, keeperServer)

	go func() {
		if err = grpcServer.Serve(listener); err != nil {
			log.Errorf("failed to serve: %s", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	<-sigChan

	grpcServer.GracefulStop()

	common.SaveConfigOnDemand(c.saveConfig, c.configFile, "./gokeeperd.yaml")
}
