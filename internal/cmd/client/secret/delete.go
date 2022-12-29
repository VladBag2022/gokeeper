package secret

import (
	"context"
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/VladBag2022/gokeeper/internal/cmd"
	pb "github.com/VladBag2022/gokeeper/internal/proto"
)

var deleteCmd = &cobra.Command{
	Use:     "delete <store_id>",
	Example: "delete 1",
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(1)(cmd, args); err != nil {
			return err
		}
		_, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse secret ID: %s", err)
		}

		return nil
	},
	Run: deleteRun,
}

func init() {
	Cmd.AddCommand(deleteCmd)
}

func deleteRun(_ *cobra.Command, args []string) {
	ctx := context.Background()

	secretID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return
	}

	rpcClient, err := cmd.NewGRPCClient()
	if err != nil {
		return
	}

	if _, err = rpcClient.Keeper.DeleteSecret(ctx, &pb.StoredSecret{Id: secretID}); err != nil {
		log.Errorf("failed to delete secret: %s", err)
	}
}
