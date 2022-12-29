package client

import (
	"context"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/VladBag2022/gokeeper/internal/client"
	"github.com/VladBag2022/gokeeper/internal/cmd"
	"github.com/VladBag2022/gokeeper/internal/cmd/client/secret/store"
	pb "github.com/VladBag2022/gokeeper/internal/proto"
)

var signCmd = &cobra.Command{
	Use:     "sign -u <username> -n",
	Example: "sign -u username -n",
	Run:     signRun,
}

func init() {
	signCmd.PersistentFlags().StringP("username", "u", "", "username to use")
	signCmd.PersistentFlags().BoolP("new", "n", false, "sign-in with new account")

	if err := viper.BindPFlags(signCmd.PersistentFlags()); err != nil {
		log.Errorf("failed to bind flags: %s", err)
	}

	if err := signCmd.MarkPersistentFlagRequired("username"); err != nil {
		log.Errorf("failed to mark username flag as required: %s", err)
	}
}

func signRun(_ *cobra.Command, _ []string) {
	ctx := context.Background()

	rpcClient, err := cmd.NewGRPCClient()
	if err != nil {
		return
	}

	credentials := &pb.Credentials{
		Username: viper.GetString("Username"),
	}

	prompt := &survey.Password{Message: "Password"}
	if err = survey.AskOne(prompt, &credentials.Password); err != nil {
		log.Errorf("failed to prompt password: %s", err)

		return
	}

	sessionManager, err := client.NewSessionManagerFromPassword(credentials.GetPassword())
	if err != nil {
		log.Errorf("failed to create session manager: %s", err)

		return
	}

	var jwt *pb.JWT
	if viper.GetBool("New") {
		jwt, err = rpcClient.Auth.SignUp(ctx, credentials)
		if err != nil {
			log.Errorf("failed to sign up: %s", err)

			return
		}
	} else {
		jwt, err = rpcClient.Auth.SignIn(ctx, credentials)
		if err != nil {
			log.Errorf("failed to sign in: %s", err)

			return
		}
	}

	viper.Set("JWT", jwt.GetToken())
	fmt.Printf("JWT acquired: %s\n", jwt.GetToken())

	viper.Set("SessionKey", sessionManager.GetSessionKey())
	fmt.Printf("Session key: %s\n", sessionManager.GetSessionKey())

	store.Secret(&pb.Secret{
		Data: []byte(sessionManager.GetEncryptedKey()),
		Kind: pb.SecretKind_SECRET_ENCRYPTED_KEY,
	})
}
