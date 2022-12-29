// Package secret contains "store", "delete" and "get" commands for secrets.
package secret

import (
	"github.com/spf13/cobra"

	"github.com/VladBag2022/gokeeper/internal/cmd/client/secret/store"
)

// Cmd is the primary command - "secret".
var Cmd = &cobra.Command{
	Use: "secret",
}

func init() {
	Cmd.AddCommand(store.Cmd)
}
