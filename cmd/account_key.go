package cmd

import (
	"encoding/hex"
	"fmt"
	"log"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/account"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/encryption"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/logger"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/paths"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/spf13/cobra"
)

const showKeyFatalMessage = "Fatal error while extracting private key"

// ShowKeyCmd is executed when when "key" flag is passed after "account" flag and is used for revealing crypto wallet's private key.
var showKeyCmd = &cobra.Command{
	Use:   "key",
	Short: "discloses your private key",
	Long:  "discloses your private key",
	Run: func(cmd *cobra.Command, args []string) {
		const location = "showKeyCmd->"
		fmt.Println("Never disclose this key. Anyone with your private keys can steal any assets held in your account")

		nodeAccount, password, err := account.Unlock()
		if err != nil {
			logger.Log(logger.MarkLocation(location, err))
			log.Fatal(showKeyFatalMessage)
		}

		scryptN, scryptP := encryption.GetScryptParams()

		ks := keystore.NewKeyStore(paths.List().AccsDir, scryptN, scryptP)

		keyJson, err := ks.Export(*nodeAccount, password, password)
		if err != nil {
			logger.Log(logger.MarkLocation(location, err))
			log.Fatal(showKeyFatalMessage)
		}

		key, err := keystore.DecryptKey(keyJson, password)
		if err != nil {
			logger.Log(logger.MarkLocation(location, err))
			log.Fatal(showKeyFatalMessage)
		}

		fmt.Println("Private Key:", hex.EncodeToString(key.PrivateKey.D.Bytes()))
	},
}

func init() {
	accountCmd.AddCommand(showKeyCmd)
}
