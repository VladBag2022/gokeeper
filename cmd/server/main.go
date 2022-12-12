package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/VladBag2022/gokeeper/internal/cmd/server"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

const NA string = "N/A"

func main() {
	if len(buildVersion) == 0 {
		buildVersion = NA
	}
	if len(buildDate) == 0 {
		buildDate = NA
	}
	if len(buildCommit) == 0 {
		buildCommit = NA
	}
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	if err := server.RootCmd.Execute(); err != nil {
		log.Errorf("failed to execute root command. %s", err)
	}
}
