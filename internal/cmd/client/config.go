package client

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	common "github.com/VladBag2022/gokeeper/internal/cmd"
	"github.com/VladBag2022/gokeeper/internal/cmd/client/meta"
	"github.com/VladBag2022/gokeeper/internal/cmd/client/secret"
)

// CLI is the entry point for GoKeeper client CLI.
type CLI struct {
	cmd *cobra.Command

	configFile string
	saveConfig bool
}

// NewCLI returns CLI.
func NewCLI() *CLI {
	var cli CLI

	cli.cmd = &cobra.Command{
		Use:     "gokeeper -s <host>:<port>",
		Example: "gokeeper -s 127.0.0.1:8080",
	}

	cobra.OnInitialize(common.InitConfig(&cli.configFile))

	cli.cmd.PersistentFlags().StringVarP(&cli.configFile, "config", "c", "", "config file")
	cli.cmd.PersistentFlags().StringP("server", "r", "", "server address: host:port")
	cli.cmd.PersistentFlags().StringP("jwt", "j", "", "JWT")
	cli.cmd.PersistentFlags().StringP("session-key", "k", "", "session key")
	cli.cmd.PersistentFlags().BoolVarP(&cli.saveConfig, "save", "s", false, "save configuration including acquired JWT")
	cli.cmd.PersistentFlags().BoolP("verbose", "v", false, "be verbose")

	for viperKey, flagName := range map[string]string{
		"ServerAddress": "server",
		"JWT":           "jwt",
		"SessionKey":    "session-key",
		"Verbose":       "verbose",
	} {
		if err := viper.BindPFlag(viperKey, cli.cmd.PersistentFlags().Lookup(flagName)); err != nil {
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

	if err := cli.cmd.MarkPersistentFlagRequired("server"); err != nil {
		log.Errorf("failed to mark server flag as required: %s", err)
	}

	cli.cmd.AddCommand(newSignCLI(), secret.NewCLI(), meta.NewCLI())

	return &cli
}
