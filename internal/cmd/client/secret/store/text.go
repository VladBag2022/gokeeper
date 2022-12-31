package store

import (
	"strings"

	"github.com/spf13/cobra"

	pb "github.com/VladBag2022/gokeeper/internal/proto"
)

func newTextCLI() *cobra.Command {
	return &cobra.Command{
		Use:     "text <string>",
		Example: "text secret text",
		Run:     textRun,
	}
}

func textRun(_ *cobra.Command, args []string) {
	text := strings.Join(args, " ")

	Secret(&pb.Secret{
		Data: []byte(text),
		Kind: pb.SecretKind_SECRET_TEXT,
	})
}
