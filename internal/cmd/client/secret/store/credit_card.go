package store

import (
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"

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
	month, err := strconv.ParseInt(args[1], 10, 32)
	if err != nil {
		log.Errorf("failed to parse month: %s", err)
		return
	}

	year, err := strconv.ParseInt(args[2], 10, 32)
	if err != nil {
		log.Errorf("failed to parse year: %s", err)
		return
	}

	cvv, err := strconv.ParseInt(args[3], 10, 32)
	if err != nil {
		log.Errorf("failed to parse CVV: %s", err)
		return
	}

	data, err := proto.Marshal(&pb.CreditCard{
		Number: args[0],
		Month:  int32(month),
		Year:   int32(year),
		Cvv:    int32(cvv),
		Owner:  strings.Join(args[4:], " "),
	})
	if err != nil {
		log.Errorf("failed to marshal credit card: %s", err)
		return
	}

	storeSecret(&pb.Secret{
		Data: data,
		Kind: pb.SecretKind_SECRET_CREDIT_CARD,
	})
}
