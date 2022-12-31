// Package meta contains "store" and "delete" commands for meta.
package meta

import "github.com/spf13/cobra"

// NewCLI returns meta CLI.
func NewCLI() *cobra.Command {
	cli := &cobra.Command{
		Use: "meta",
	}

	cli.AddCommand(newStoreCLI(), newDeleteCLI())

	return cli
}
