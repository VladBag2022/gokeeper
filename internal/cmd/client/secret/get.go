package secret

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"

	"github.com/VladBag2022/gokeeper/internal/cmd"
	"github.com/VladBag2022/gokeeper/internal/crypt"
	pb "github.com/VladBag2022/gokeeper/internal/proto"
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

	rpcClient, err := cmd.NewGRPCClient(true)
	if err != nil {
		return
	}

	coder, err := crypt.NewCoder(rpcClient.SessionKey)
	if err != nil {
		log.Errorf("failed to create coder: %s", err)
		return
	}

	secrets, err := rpcClient.Keeper.GetSecrets(ctx, &empty.Empty{})
	if err != nil {
		log.Errorf("failed to get secrets: %s", err)
		return
	}

fs:
	for _, secret := range secrets.GetSecrets() {
		data, dErr := coder.Decrypt(secret.GetSecret().GetData())
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
