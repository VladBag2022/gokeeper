package store

import (
	"context"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/VladBag2022/gokeeper/internal/cmd"
	"github.com/VladBag2022/gokeeper/internal/crypt"
	pb "github.com/VladBag2022/gokeeper/internal/proto"
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

func storeSecret(secret *pb.Secret) {
	ctx := context.Background()

	rpcClient, err := cmd.NewGRPCClient()
	if err != nil {
		return
	}

	var password string
	prompt := &survey.Password{Message: "Encryption password"}
	if err = survey.AskOne(prompt, &password); err != nil {
		log.Errorf("failed to prompt encryption password: %s", err)
		return
	}

	coder, err := crypt.NewCoder([]byte(password))
	if err != nil {
		log.Errorf("failed to create coder: %s", err)
		return
	}
	secret.Data = coder.Encrypt(secret.GetData())

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
