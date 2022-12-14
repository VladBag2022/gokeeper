package meta

import (
	"context"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/VladBag2022/gokeeper/internal/cmd"
	pb "github.com/VladBag2022/gokeeper/internal/proto"
)

func init() {
	storeCmd.PersistentFlags().IntP("secret", "s", 0, "secret ID")
	storeCmd.PersistentFlags().IntP("meta", "m", 0, "meta ID to update")

	if err := viper.BindPFlags(storeCmd.PersistentFlags()); err != nil {
		log.Errorf("failed to bind flags: %s", err)
	}

	storeCmd.MarkFlagsMutuallyExclusive("secret", "meta")
	Cmd.AddCommand(storeCmd)
}

var storeCmd = &cobra.Command{
	Use:     "store [-s <secret_id>] [-m <meta_id>] <string>",
	Example: "store -s 1 top secret",
	Args: func(cmd *cobra.Command, args []string) error {
		return cobra.MinimumNArgs(1)(cmd, args)
	},
	Run: storeRun,
}

func storeRun(_ *cobra.Command, args []string) {
	ctx := context.Background()

	text := strings.Join(args[1:], " ")
	meta := &pb.Meta{Text: text}

	rpcClient, err := cmd.NewGRPCClient()
	if err != nil {
		return
	}

	metaID := viper.GetInt64("meta")
	secretID := viper.GetInt64("secret")
	if metaID > 0 {
		_, err = rpcClient.Keeper.UpdateMeta(ctx, &pb.ClientMeta{
			Meta: meta,
			Id:   metaID,
		})
		if err != nil {
			log.Errorf("failed to update meta: %s", err)
		}
		return
	}

	clientMeta, err := rpcClient.Keeper.StoreMeta(ctx, &pb.StoreMetaRequest{
		Meta:     meta,
		SecretId: secretID,
	})
	if err != nil {
		log.Errorf("failed to store meta: %s", err)
		return
	}
	fmt.Printf("Meta ID: %d\n", clientMeta.GetId())
}
