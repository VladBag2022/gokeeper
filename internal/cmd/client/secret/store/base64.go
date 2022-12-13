package store

import (
	"github.com/spf13/cobra"

	pb "github.com/VladBag2022/gokeeper/internal/proto"
)

var base64Cmd = &cobra.Command{
	Use: "base64 base64-string",
	Args: func(cmd *cobra.Command, args []string) error {
		return cobra.ExactArgs(1)(cmd, args)
	},
	Run: base64Run,
}

func base64Run(_ *cobra.Command, args []string) {
	storeSecret(&pb.Secret{
		Data: []byte(args[0]),
		Kind: pb.SecretKind_SECRET_TEXT,
	})
}
