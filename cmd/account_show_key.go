package cmd

import (
	"dfile-secondary-node/account"
	"dfile-secondary-node/paths"
	"dfile-secondary-node/shared"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/spf13/cobra"
)

const showKeyFatalMessage = "Fatal error while extracting private key"

var showKeyCmd = &cobra.Command{
	Use:   "key",
	Short: "reveals you the private key",
	Long:  "reveals you the private key",
	Run: func(cmd *cobra.Command, args []string) {
		const logInfo = "showKeyCmd->"
		fmt.Println("Never disclose this key. Anyone with your private keys can steal any assets held in your account")

		etherAccount, password, err := account.ValidateUser()
		if err != nil {
			shared.LogError(logInfo, err)
			log.Fatal(showKeyFatalMessage)
		}

		ks := keystore.NewKeyStore(paths.AccsDirPath, keystore.StandardScryptN, keystore.StandardScryptP)

		keyJson, err := ks.Export(*etherAccount, password, password)
		if err != nil {
			shared.LogError(logInfo, shared.GetDetailedError(err))
			log.Fatal(showKeyFatalMessage)
		}

		key, err := keystore.DecryptKey(keyJson, password)
		if err != nil {
			shared.LogError(logInfo, shared.GetDetailedError(err))
			log.Fatal(showKeyFatalMessage)
		}

		fmt.Println("Private Key:", hex.EncodeToString(key.PrivateKey.D.Bytes()))
	},
}

func init() {
	accountCmd.AddCommand(showKeyCmd)
}
