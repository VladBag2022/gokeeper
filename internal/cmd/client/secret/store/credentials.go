package store

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"

	pb "github.com/VladBag2022/gokeeper/internal/proto"
)

var credentialsCmd = &cobra.Command{
	Use:     "credentials <username> <password>",
	Example: "credentials user password",
	Args: func(cmd *cobra.Command, args []string) error {
		return cobra.ExactArgs(2)(cmd, args)
	},
	Run: credentialsRun,
}

func init() {
	Cmd.AddCommand(credentialsCmd)
}

func credentialsRun(_ *cobra.Command, args []string) {
	data, err := proto.Marshal(&pb.Credentials{
		Username: args[0],
		Password: args[1],
	})
	if err != nil {
		log.Errorf("failed to marshal credentials: %s", err)
		return
	}

	storeSecret(&pb.Secret{
		Data: data,
		Kind: pb.SecretKind_SECRET_TEXT,
	})
}
