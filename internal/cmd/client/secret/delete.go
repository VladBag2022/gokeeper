package secret

import (
	"context"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/VladBag2022/gokeeper/internal/cmd/client"
	pb "github.com/VladBag2022/gokeeper/internal/proto"
)

var deleteCmd = &cobra.Command{
	Use: "delete <store_id>",
	Example: "delete 1",
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}
		_, err := strconv.ParseInt(args[0], 10,  64)
		if err != nil {
			return err
		}
		return nil
	},
	Run: deleteRun,
}

func deleteRun(_ *cobra.Command, args []string) {
	ctx := context.Background()

	secretID, err := strconv.ParseInt(args[0], 10,  64)
	if err != nil {
		return 
	}

	rpcClient, err := client.NewRPCClient()
	if err != nil {
		return
	}

	if _, err = rpcClient.Keeper.DeleteSecret(ctx, &pb.ClientSecret{Id: secretID}); err != nil {
		log.Errorf("failed to delete secret: %s", err)
	}
}
