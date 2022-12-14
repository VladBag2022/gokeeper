package store

import (
	"strings"

	"github.com/spf13/cobra"

	pb "github.com/VladBag2022/gokeeper/internal/proto"
)

var creditCardCmd = &cobra.Command{
	Use:     "credit-card <number> <month> <year> <cvv> <username>",
	Example: "credit-card 1234123412341234 10 2020 123 User User",
	Args: func(cmd *cobra.Command, args []string) error {
		return cobra.MinimumNArgs(5)(cmd, args)
	},
	Run: creditCardRun,
}

func init() {
	Cmd.AddCommand(creditCardCmd)
}

func creditCardRun(_ *cobra.Command, args []string) {
	number := args[0]
	month := args[1]
	year := args[2]
	cvv := args[3]
	name := strings.Join(args[4:], " ")

	text := strings.Join([]string{number, month, year, cvv, name}, "")

	storeSecret(&pb.Secret{
		Data: []byte(text),
		Kind: pb.SecretKind_SECRET_TEXT,
	})
}
