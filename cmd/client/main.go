// Package main runs GoKeeper client CLI.
package main

import (
	"fmt"

	"github.com/VladBag2022/gokeeper/internal/cmd/client"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

// NA means "Not Assigned".
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
	fmt.Printf("Build commit: %s\n\n", buildCommit)

	client.Execute()
}
