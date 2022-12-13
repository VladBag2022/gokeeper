package meta

import (
	"context"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/VladBag2022/gokeeper/internal/cmd/client"
	pb "github.com/VladBag2022/gokeeper/internal/proto"
)

var deleteCmd = &cobra.Command{
	Use: "delete <meta_id>",
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

func init() {
	cmd.AddCommand(deleteCmd)
}

func deleteRun(_ *cobra.Command, args []string) {
	ctx := context.Background()

	metaID, err := strconv.ParseInt(args[0], 10,  64)
	if err != nil {
		return
	}

	rpcClient, err := client.NewRPCClient()
	if err != nil {
		return
	}

	if _, err = rpcClient.Keeper.DeleteMeta(ctx, &pb.ClientMeta{Id: metaID}); err != nil {
		log.Errorf("failed to delete meta: %s", err)
	}
}
