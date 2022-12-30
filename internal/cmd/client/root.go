// Package client contains root command for GoKeeper client CLI.
package client

import (
	log "github.com/sirupsen/logrus"

	common "github.com/VladBag2022/gokeeper/internal/cmd"
)

// Execute starts GoKeeper client CLI.
func (c *CLI) Execute() {
	if err := c.cmd.Execute(); err != nil {
		log.Errorf("failed to execute root command: %s", err)
	}

	common.SaveConfigOnDemand(c.saveConfig, c.configFile, "./gokeeper.yaml")
}
