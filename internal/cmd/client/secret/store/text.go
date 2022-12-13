package store

import (
	"strings"

	"github.com/spf13/cobra"

	pb "github.com/VladBag2022/gokeeper/internal/proto"
)

var textCmd = &cobra.Command{
	Use: "text secret string",
	Run: textRun,
}

func textRun(_ *cobra.Command, args []string) {
	text := strings.Join(args, " ")
	storeSecret(&pb.Secret{
		Data: []byte(text),
		Kind: pb.SecretKind_SECRET_TEXT,
	})
}
