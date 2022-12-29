package server

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	common "github.com/VladBag2022/gokeeper/internal/cmd"
)

var (
	rootCmd = &cobra.Command{
		Use: "gokeeperd -a <host>:<port> -d <database_dsn> -c <cert_pem> -k <key_pem> -j <jwt_duration>",
		Example: "gokeeperd -a 127.0.0.1:8080 -d postgresql://user:password@127.0.0.1:5432/database " +
			"-c cert.pem -k key.pem -j 10m",
		Run: rootRun,
	}

	configFile string
	saveConfig bool
)

func init() {
	cobra.OnInitialize(common.InitConfig(&configFile))

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file")
	rootCmd.PersistentFlags().StringP("address", "a", "", "GRPC listen address: host:port")
	rootCmd.PersistentFlags().StringP("database", "d", "", "Postgres database DSN")
	rootCmd.PersistentFlags().StringP("cert", "t", "", "TLS cert file in PEM format")
	rootCmd.PersistentFlags().StringP("key", "k", "", "TLS key file in PEM format")
	rootCmd.PersistentFlags().DurationP("jwt", "j", 0, "JWT duration")
	rootCmd.PersistentFlags().BoolVarP(&saveConfig, "save", "s", false, "save configuration including acquired JWT")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "be verbose")

	for viperKey, flagName := range map[string]string{
		"ListenAddress": "address",
		"DatabaseDSN":   "database",
		"CertPEM":       "cert",
		"KeyPEM":        "key",
		"JWTDuration":   "jwt",
		"Verbose":       "verbose",
	} {
		if err := viper.BindPFlag(viperKey, rootCmd.PersistentFlags().Lookup(flagName)); err != nil {
			log.Errorf("failed to bind flag %s: %s", flagName, err)
		}
	}

	for viperKey, envVar := range map[string]string{
		"ListenAddress": "SERVER_ADDRESS",
		"DatabaseDSN":   "DATABASE_DSN",
		"CertPEMFile":   "CERT_PEM",
		"KeyPEMFile":    "KEY_PEM",
		"JWTDuration":   "JWT_DURATION",
		"Verbose":       "VERBOSE",
	} {
		if err := viper.BindEnv(viperKey, envVar); err != nil {
			log.Errorf("failed to bind environment variable %s: %s", envVar, err)
		}
	}

	viper.SetDefault("ListenAddress", "127.0.0.1:8080")
	viper.SetDefault("CertPEMFile", "cert.pem")
	viper.SetDefault("KeyPEMFile", "key.pem")
	viper.SetDefault("JWTDuration", time.Minute*10)

	if err := rootCmd.MarkPersistentFlagRequired("database"); err != nil {
		log.Errorf("failed to mark database flag as required: %s", err)
	}
}
