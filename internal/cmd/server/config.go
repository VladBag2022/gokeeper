package server

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	RootCmd = &cobra.Command{
		Use: "gokeeperd -a 127.0.0.1:8080 -d postgresql://user:password@127.0.0.1:5432/database -c cert.pem -k key.pem",
		Run: rootRun,
	}

	configFile string
)

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file")
	RootCmd.PersistentFlags().StringP("address", "a", "", "GRPC listen address: host:port")
	RootCmd.PersistentFlags().StringP("database", "d", "", "Postgres database DSN")
	RootCmd.PersistentFlags().StringP("cert", "r", "", "TLS cert file in PEM format")
	RootCmd.PersistentFlags().StringP("key", "k", "", "TLS key file in PEM format")
	RootCmd.PersistentFlags().DurationP("jwt", "j", 0, "JWT duration")

	for key, f := range map[string]string{
		"ListenAddress": "address",
		"DatabaseDSN":   "database",
		"CertPEM":       "cert",
		"KeyPEM":        "key",
		"JWTDuration":   "jwt",
	} {
		if err := viper.BindPFlag(key, RootCmd.PersistentFlags().Lookup(f)); err != nil {
			log.Errorf("failed to bind flag %s. %s", f, err)
		}
	}

	for key, env := range map[string]string{
		"ListenAddress": "SERVER_ADDRESS",
		"DatabaseDSN":   "DATABASE_DSN",
		"CertPEMFile":   "CERT_PEM",
		"KeyPEMFile":    "KEY_PEM",
		"JWTDuration":   "JWT_DURATION",
	} {
		if err := viper.BindEnv(key, env); err != nil {
			log.Errorf("failed to bind environment variable %s. %s", env, err)
		}
	}

	viper.SetDefault("ListenAddress", "127.0.0.1:8080")
	viper.SetDefault("CertPEMFile", "cert.pem")
	viper.SetDefault("KeyPEMFile", "key.pem")
	viper.SetDefault("JWTDuration", time.Minute*10)
}

func initConfig() {
	if configFile == "" {
		configFile = os.Getenv("CONFIG")
	}

	if configFile != "" {
		viper.SetConfigFile(configFile)
		err := viper.ReadInConfig()
		if err != nil {
			log.Errorf("failed to read config file: %s", err)
		}
	}
}
