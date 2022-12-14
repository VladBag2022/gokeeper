package secret

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/VladBag2022/gokeeper/internal/cmd"
)

var getCmd = &cobra.Command{
	Use: "get",
	Run: getRun,
}

func init() {
	Cmd.AddCommand(getCmd)
}

func getRun(_ *cobra.Command, _ []string) {
	ctx := context.Background()

	rpcClient, err := cmd.NewGRPCClient()
	if err != nil {
		return
	}

	secrets, err := rpcClient.Keeper.GetSecrets(ctx, &empty.Empty{})
	if err != nil {
		log.Errorf("failed to get secrets: %s", err)
		return
	}

	for _, secret := range secrets.GetSecrets() {
		fmt.Printf("[%d] %s\n", secret.GetId(), secret.GetSecret().GetData())
		for _, meta := range secret.GetMeta() {
			fmt.Printf("* [%d] %s\n", meta.GetId(), meta.GetMeta().GetText())
		}
	}
}
