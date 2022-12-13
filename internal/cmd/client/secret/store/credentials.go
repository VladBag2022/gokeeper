package store

import (
	"strings"

	"github.com/spf13/cobra"

	pb "github.com/VladBag2022/gokeeper/internal/proto"
)

var credentialsCmd = &cobra.Command{
	Use: "credentials <username> <password>",
	Example: "credentials user password",
	Args: func(cmd *cobra.Command, args []string) error {
		return cobra.ExactArgs(2)(cmd, args)
	},
	Run: credentialsRun,
}

func init() {
	cmd.AddCommand(credentialsCmd)
}

func credentialsRun(_ *cobra.Command, args []string) {
	username := args[0]
	password := args[1]

	text := strings.Join([]string{username, password}, "")

	storeSecret(&pb.Secret{
		Data: []byte(text),
		Kind: pb.SecretKind_SECRET_TEXT,
	})
}
