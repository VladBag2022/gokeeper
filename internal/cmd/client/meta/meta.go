package meta

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use: "meta",
}

func init() {
	Cmd.AddCommand(storeCmd, deleteCmd)
}
