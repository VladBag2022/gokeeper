package meta

import (
	"github.com/spf13/cobra"

	"github.com/VladBag2022/gokeeper/internal/cmd/client"
)

var cmd = &cobra.Command{
	Use: "meta",
}

func init() {
	client.RootCmd.AddCommand(cmd)
}
