// Package client contains root command for GoKeeper client CLI.
package client

import (
	log "github.com/sirupsen/logrus"

	common "github.com/VladBag2022/gokeeper/internal/cmd"
)

// Execute start GoKeeper client CLI.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Errorf("failed to execute root command: %s", err)
	}

	common.SaveConfigOnDemand(saveConfig, configFile, "./gokeeper.yaml")
}
