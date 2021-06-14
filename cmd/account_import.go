package cmd

import (
	"dfile-secondary-node/account"
	"dfile-secondary-node/server"
	"dfile-secondary-node/shared"
	"log"

	"github.com/spf13/cobra"
)

// accountListCmd represents the list command
var accountImportCmd = &cobra.Command{
	Use:   "import",
	Short: "imports your account by private key",
	Long:  "imports your account by private key",
	Run: func(cmd *cobra.Command, args []string) {
		const info = "accountImportCmd"
		accountStr, nodeConfig, err := account.Import()
		if err != nil {
			shared.LogError(info + ":" + err.Error())
			log.Fatal("Fatal error, couldn't import an account")
		}

		server.Start(accountStr, nodeConfig.HTTPPort)
	},
}

func init() {
	accountCmd.AddCommand(accountImportCmd)
}
