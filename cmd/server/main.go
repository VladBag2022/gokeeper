// Package main runs GoKeeper server CLI.
package main

import (
	"fmt"

	"github.com/VladBag2022/gokeeper/internal/cmd/server"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

// NAText means "Not Assigned".
const NAText string = "N/A"

func main() {
	if len(buildVersion) == 0 {
		buildVersion = NAText
	}

	if len(buildDate) == 0 {
		buildDate = NAText
	}

	if len(buildCommit) == 0 {
		buildCommit = NAText
	}
	
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n\n", buildCommit)

	server.Execute()
}
