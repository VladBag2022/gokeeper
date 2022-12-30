// Package cmd contains GoKeeper CLI for both client and server.
package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/VladBag2022/gokeeper/internal/client"
)

// InitConfig returns function for config initialisation which loads config from file in case it was provided.
func InitConfig(configFile *string) func() {
	return func() {
		if viper.GetBool("Verbose") {
			log.SetLevel(log.DebugLevel)
		}

		if *configFile == "" {
			*configFile = os.Getenv("CONFIG")
		}

		if *configFile != "" {
			log.Infof("using config file %s", *configFile)
			viper.SetConfigFile(*configFile)
			err := viper.ReadInConfig()

			if err != nil {
				log.Errorf("failed to read config file: %s", err)
			}
		}
	}
}

// SaveConfigOnDemand is a helper function for saving configuration to file.
func SaveConfigOnDemand(saveConfig bool, configFile, defaultConfigFile string) {
	if saveConfig {
		if len(configFile) == 0 {
			configFile = defaultConfigFile
		}

		if err := viper.WriteConfigAs(configFile); err != nil {
			log.Errorf("failed to write config: %s", err)
		}
	}
}

// NewGRPCClient creates new Client using viper config.
func NewGRPCClient() (*client.Client, error) {
	authInterceptor := client.NewAuthInterceptor(viper.GetString("JWT"))

	clientConn, err := grpc.Dial(viper.GetString("ServerAddress"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(authInterceptor.Unary()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial server: %w", err)
	}

	return client.NewClient(clientConn), nil
}
