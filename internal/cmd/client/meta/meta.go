// Package meta contains "store" and "delete" commands for meta.
package meta

import (
	"github.com/spf13/cobra"
)

// Cmd is the primary command - "meta".
var Cmd = &cobra.Command{
	Use: "meta",
}
