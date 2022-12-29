package client

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	common "github.com/VladBag2022/gokeeper/internal/cmd"
	"github.com/VladBag2022/gokeeper/internal/cmd/client/meta"
	"github.com/VladBag2022/gokeeper/internal/cmd/client/secret"
)

var (
	rootCmd = &cobra.Command{
		Use:     "gokeeper -s <host>:<port>",
		Example: "gokeeper -s 127.0.0.1:8080",
	}

	configFile string
	saveConfig bool
)

func init() {
	cobra.OnInitialize(common.InitConfig(&configFile))

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file")
	rootCmd.PersistentFlags().StringP("server", "r", "", "server address: host:port")
	rootCmd.PersistentFlags().StringP("jwt", "j", "", "JWT")
	rootCmd.PersistentFlags().StringP("session-key", "k", "", "session key")
	rootCmd.PersistentFlags().BoolVarP(&saveConfig, "save", "s", false, "save configuration including acquired JWT")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "be verbose")

	for viperKey, flagName := range map[string]string{
		"ServerAddress": "server",
		"JWT":           "jwt",
		"SessionKey":    "session-key",
		"Verbose":       "verbose",
	} {
		if err := viper.BindPFlag(viperKey, rootCmd.PersistentFlags().Lookup(flagName)); err != nil {
			log.Errorf("failed to bind flag %s. %s", flagName, err)
		}
	}

	for viperKey, envVar := range map[string]string{
		"ServerAddress": "SERVER_ADDRESS",
		"JWT":           "JWT",
		"SessionKey":    "SESSION_KEY",
		"Verbose":       "VERBOSE",
	} {
		if err := viper.BindEnv(viperKey, envVar); err != nil {
			log.Errorf("failed to bind environment variable %s. %s", envVar, err)
		}
	}

	if err := rootCmd.MarkPersistentFlagRequired("server"); err != nil {
		log.Errorf("failed to mark server flag as required: %s", err)
	}

	rootCmd.AddCommand(signCmd, secret.Cmd, meta.Cmd)
}
