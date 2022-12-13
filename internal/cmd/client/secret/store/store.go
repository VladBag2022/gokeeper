package store

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/VladBag2022/gokeeper/internal/cmd/client"
	"github.com/VladBag2022/gokeeper/internal/cmd/client/secret"
	pb "github.com/VladBag2022/gokeeper/internal/proto"
)

var cmd = &cobra.Command{
	Use: "store",
}

func init() {
	cmd.PersistentFlags().IntP("id", "i", 0, "secret ID to update")

	if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
		log.Errorf("failed to bind flags: %s", err)
	}

	secret.Cmd.AddCommand(cmd)
}

func storeSecret(secret *pb.Secret) {
	ctx := context.Background()

	rpcClient, err := client.NewRPCClient()
	if err != nil {
		return
	}

	secretID := viper.GetInt64("id")
	if secretID > 0 {
		_, err = rpcClient.Keeper.UpdateSecret(ctx, &pb.ClientSecret{
			Secret: secret,
			Id:     secretID,
		})
		if err != nil {
			log.Errorf("failed to update secret: %s", err)
		}
		return
	}

	clientSecret, err := rpcClient.Keeper.StoreSecret(ctx, secret)
	if err != nil {
		log.Errorf("failed to store secret: %s", err)
		return
	}
	fmt.Printf("Secret ID: %d\n", clientSecret.GetId())
}
