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

	"github.com/VladBag2022/gokeeper/internal/jwt"
	pb "github.com/VladBag2022/gokeeper/internal/proto"
	"github.com/VladBag2022/gokeeper/internal/server"
	"github.com/VladBag2022/gokeeper/internal/storage"
	"github.com/VladBag2022/gokeeper/internal/utils"
)

func rootRun(_ *cobra.Command, _ []string) {
	ctx := context.Background()

	listener, err := net.Listen("tcp", viper.GetString("ListenAddress"))
	if err != nil {
		log.Errorf("failed to listen: %s", err)
		return
	}

	jwtKey, err := utils.GenerateRandomBytes(32)
	if err != nil {
		log.Errorf("failed to generate random JWT key: %s", err)
		return
	}

	jwtManager := jwt.NewManager(jwtKey, viper.GetDuration("JWTDuration"))

	databaseDSN := viper.GetString("DatabaseDSN")
	if len(databaseDSN) == 0 {
		log.Error("Postgres database DSN is required")
		return
	}

	store, err := storage.NewPostgresStore(ctx, databaseDSN)
	if err != nil {
		log.Errorf("failed to connect to Postgres store: %s", err)
		return
	}
	defer func() {
		if err = store.Close(); err != nil {
			log.Errorf("failed to close store: %s", err)
		}
	}()

	authInterceptor := server.NewAuthInterceptor(jwtManager)
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
	)

	authServer := server.NewAuthServer(store, jwtManager)
	keeperServer := server.NewKeeperServer(store)

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
}
