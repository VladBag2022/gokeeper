package store

import (
	"strings"

	"github.com/spf13/cobra"

	pb "github.com/VladBag2022/gokeeper/internal/proto"
)

var textCmd = &cobra.Command{
	Use:     "text <string>",
	Example: "text secret text",
	Run:     textRun,
}

func init() {
	Cmd.AddCommand(textCmd)
}

func textRun(_ *cobra.Command, args []string) {
	text := strings.Join(args, " ")
	Secret(&pb.Secret{
		Data: []byte(text),
		Kind: pb.SecretKind_SECRET_TEXT,
	})
}
