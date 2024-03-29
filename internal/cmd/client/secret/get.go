package secret

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/proto"

	"github.com/VladBag2022/gokeeper/internal/client"
	"github.com/VladBag2022/gokeeper/internal/cmd"
	pb "github.com/VladBag2022/gokeeper/internal/proto"
)

func newGetCLI() *cobra.Command {
	return &cobra.Command{
		Use: "get",
		Run: getRun,
	}
}

func getRun(_ *cobra.Command, _ []string) {
	ctx := context.Background()

	rpcClient, err := cmd.NewGRPCClient()
	if err != nil {
		return
	}

	key, err := rpcClient.Keeper.GetEncryptedKey(ctx, &empty.Empty{})
	if err != nil {
		log.Errorf("failed to get encrypted key: %s", err)

		return
	}

	sessionManager, err := client.NewSessionManagerFromEncryptedKey(
		string(key.GetSecret().GetData()),
		viper.GetString("SessionKey"))
	if err != nil {
		log.Errorf("failed to create session manager: %s", err)

		return
	}

	secrets, err := rpcClient.Keeper.GetSecrets(ctx, &empty.Empty{})
	if err != nil {
		log.Errorf("failed to get secrets: %s", err)

		return
	}

	printSecrets(secrets, sessionManager)
}

func printSecrets(secrets *pb.StoredSecrets, sessionManager *client.SessionManager) {
fs:
	for _, secret := range secrets.GetSecrets() {
		data, dErr := sessionManager.Coder.Decrypt(secret.GetSecret().GetData())
		if dErr != nil {
			log.Errorf("failed to decrypt secret data: %s", dErr)

			continue
		}

		var text string
		switch secret.GetSecret().GetKind() {
		case pb.SecretKind_SECRET_CREDENTIALS:
			credentials := &pb.Credentials{}
			if uErr := proto.Unmarshal(data, credentials); uErr != nil {
				log.Errorf("failed to unmarshal credentials: %s", uErr)

				continue fs
			}
			text = credentials.String()
		case pb.SecretKind_SECRET_TEXT, pb.SecretKind_SECRET_BLOB:
			text = string(data)
		case pb.SecretKind_SECRET_CREDIT_CARD:
			card := &pb.CreditCard{}
			if uErr := proto.Unmarshal(data, card); uErr != nil {
				log.Errorf("failed to unmarshal credit card: %s", uErr)

				continue fs
			}
			text = card.String()
		}

		fmt.Printf("[%d] %s\n", secret.GetId(), text)
		for _, meta := range secret.GetMeta() {
			fmt.Printf("* [%d] %s\n", meta.GetId(), meta.GetMeta().GetText())
		}
	}
}
