// Package secret contains "store", "delete" and "get" commands for secrets.
package secret

import (
	"github.com/spf13/cobra"

	"github.com/VladBag2022/gokeeper/internal/cmd/client/secret/store"
)

// NewCLI returns secret CLI.
func NewCLI() *cobra.Command {
	cli := &cobra.Command{
		Use: "secret",
	}

	cli.AddCommand(newGetCLI(), store.NewCLI(), newDeleteCLI())

	return cli
}
