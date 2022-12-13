package client

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	rpc "github.com/VladBag2022/gokeeper/internal/client"
	common "github.com/VladBag2022/gokeeper/internal/cmd"
)

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Errorf("failed to execute root command: %s", err)
	}

	common.SaveConfigOnDemand(saveConfig, configFile, "./gokeeper.yaml")
}

func NewRPCClient() (*rpc.Client, error) {
	authInterceptor := rpc.NewAuthInterceptor(viper.GetString("JWT"))

	cc, err := grpc.Dial(viper.GetString("ServerAddress"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(authInterceptor.Unary()))
	if err != nil {
		log.Errorf("failed to dial server: %s", err)
		return nil, err
	}

	return rpc.NewClient(cc), nil
}
