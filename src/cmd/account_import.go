package cmd

import (
	"fmt"
	"log"

	"github.com/DeNetPRO/src/account"
	"github.com/DeNetPRO/src/logger"
	"github.com/DeNetPRO/src/rpcserver"

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
			fmt.Println(err)
			logger.Log(logger.MarkLocation(location, err))
			log.Fatal("Fatal error, couldn't import an account")
		}

		err = rpcserver.Start(nodeConfig.HTTPPort)
		if err != nil {
			log.Fatal(err)
		}

	},
}

func init() {
	accountCmd.AddCommand(accountImportCmd)
}
