package cmd

import (
	"fmt"
	"log"

	"git.denetwork.xyz/dfile/dfile-secondary-node/account"
	"git.denetwork.xyz/dfile/dfile-secondary-node/cleaner"
	"git.denetwork.xyz/dfile/dfile-secondary-node/logger"
	"git.denetwork.xyz/dfile/dfile-secondary-node/server"
	"github.com/spf13/cobra"
)

// accountListCmd represents the list command
var accountImportCmd = &cobra.Command{
	Use:   "import",
	Short: "imports your account by private key",
	Long:  "imports your account by private key",
	Run: func(cmd *cobra.Command, args []string) {
		const location = "accountImportCmd->"
		_, nodeConfig, err := account.Import()
		if err != nil {
			logger.Log(logger.CreateDetails(location, err))
			log.Fatal("Fatal error, couldn't import an account")
		}

		account.IpAddr = fmt.Sprint(nodeConfig.IpAddress, ":", nodeConfig.HTTPPort)

		go cleaner.Start()

		server.Start(nodeConfig.HTTPPort)
	},
}

func init() {
	accountCmd.AddCommand(accountImportCmd)
}
