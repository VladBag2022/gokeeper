package client

import (
	log "github.com/sirupsen/logrus"

	common "github.com/VladBag2022/gokeeper/internal/cmd"
)

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Errorf("failed to execute root command: %s", err)
	}

	common.SaveConfigOnDemand(saveConfig, configFile, "./gokeeper.yaml")
}
