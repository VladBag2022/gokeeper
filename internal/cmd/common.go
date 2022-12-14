package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/VladBag2022/gokeeper/internal/client"
)

func InitConfig(configFile *string) func() {
	return func() {
		if *configFile == "" {
			*configFile = os.Getenv("CONFIG")
		}

		if *configFile != "" {
			viper.SetConfigFile(*configFile)
			err := viper.ReadInConfig()
			if err != nil {
				log.Errorf("failed to read config file: %s", err)
			}
		}
	}
}

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

func NewGRPCClient() (*client.Client, error) {
	authInterceptor := client.NewAuthInterceptor(viper.GetString("JWT"))

	cc, err := grpc.Dial(viper.GetString("ServerAddress"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(authInterceptor.Unary()))
	if err != nil {
		log.Errorf("failed to dial server: %s", err)
		return nil, err
	}

	return client.NewClient(cc), nil
}
