package cmd

import (
	"dfile-secondary-node/account"
	"dfile-secondary-node/cleaner"
	"dfile-secondary-node/logger"
	"dfile-secondary-node/server"
	"log"

	"github.com/spf13/cobra"
)

// accountListCmd represents the list command
var accountImportCmd = &cobra.Command{
	Use:   "import",
	Short: "imports your account by private key",
	Long:  "imports your account by private key",
	Run: func(cmd *cobra.Command, args []string) {
		const logInfo = "accountImportCmd->"
		accountStr, nodeConfig, err := account.Import()
		if err != nil {
			logger.LogError(logInfo, err)
			log.Fatal("Fatal error, couldn't import an account")
		}

		go cleaner.Start()

		server.Start(accountStr, nodeConfig.HTTPPort)
	},
}

func init() {
	accountCmd.AddCommand(accountImportCmd)
}
