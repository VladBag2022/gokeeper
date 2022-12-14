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
	RootCmd = &cobra.Command{
		Use:     "gokeeper -s <host>:<port>",
		Example: "gokeeper -s 127.0.0.1:8080",
	}

	configFile string
	saveConfig bool
)

func init() {
	cobra.OnInitialize(common.InitConfig(&configFile))

	RootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file")
	RootCmd.PersistentFlags().StringP("server", "r", "", "server address: host:port")
	RootCmd.PersistentFlags().StringP("jwt", "j", "", "JWT")
	RootCmd.PersistentFlags().BoolVarP(&saveConfig, "save", "s", false, "save configuration including acquired JWT")

	for key, f := range map[string]string{
		"ServerAddress": "server",
		"JWT":           "jwt",
	} {
		if err := viper.BindPFlag(key, RootCmd.PersistentFlags().Lookup(f)); err != nil {
			log.Errorf("failed to bind flag %s. %s", f, err)
		}
	}

	for key, env := range map[string]string{
		"ServerAddress": "SERVER_ADDRESS",
		"JWT":           "JWT",
	} {
		if err := viper.BindEnv(key, env); err != nil {
			log.Errorf("failed to bind environment variable %s. %s", env, err)
		}
	}

	if err := RootCmd.MarkPersistentFlagRequired("server"); err != nil {
		log.Errorf("failed to mark server flag as required: %s", err)
	}

	RootCmd.AddCommand(signCmd, secret.Cmd, meta.Cmd)
}
