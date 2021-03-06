package cmd

import (
	"log"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/account"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/cleaner"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/logger"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/server"
	"github.com/spf13/cobra"
)

// AccountListCmd is executed when "import" flag is passed after "account" flag and is used for importing crypto wallet.
var accountImportCmd = &cobra.Command{
	Use:   "import",
	Short: "imports your wallet by private key",
	Long:  "imports your wallet by private key",
	Run: func(cmd *cobra.Command, args []string) {
		const location = "accountImportCmd->"
		_, nodeConfig, err := account.Import()
		if err != nil {
			logger.Log(logger.CreateDetails(location, err))
			log.Fatal("Fatal error, couldn't import an account")
		}

		go cleaner.Start()

		server.Start(nodeConfig.HTTPPort)
	},
}

func init() {
	accountCmd.AddCommand(accountImportCmd)
}
