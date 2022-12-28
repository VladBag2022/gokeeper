package store

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/VladBag2022/gokeeper/internal/client"
	"github.com/VladBag2022/gokeeper/internal/cmd"
	"github.com/VladBag2022/gokeeper/internal/proto"
)

var Cmd = &cobra.Command{
	Use: "store",
}

func init() {
	Cmd.PersistentFlags().IntP("id", "i", 0, "secret ID to update")

	if err := viper.BindPFlags(Cmd.PersistentFlags()); err != nil {
		log.Errorf("failed to bind flags: %s", err)
	}
}

func Secret(secret *proto.Secret) {
	ctx := context.Background()

	rpcClient, err := cmd.NewGRPCClient()
	if err != nil {
		return
	}

	if secret.GetKind() != proto.SecretKind_SECRET_ENCRYPTED_KEY {
		key, gErr := rpcClient.Keeper.GetEncryptedKey(ctx, &empty.Empty{})
		if gErr != nil {
			log.Errorf("failed to get encrypted key: %s", gErr)
			return
		}

		sessionManager, sErr := client.NewSessionManagerFromEncryptedKey(
			string(key.GetSecret().GetData()),
			viper.GetString("SessionKey"))
		if sErr != nil {
			log.Errorf("failed to create session manager: %s", sErr)
			return
		}

		secret.Data = sessionManager.Coder.Encrypt(secret.GetData())
	}

	secretID := viper.GetInt64("id")
	if secretID > 0 {
		_, err = rpcClient.Keeper.UpdateSecret(ctx, &proto.ClientSecret{
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
