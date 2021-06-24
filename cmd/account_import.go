package cmd

import (
	"dfile-secondary-node/account"
	"dfile-secondary-node/server"
	"dfile-secondary-node/shared"
	"log"
	"strconv"

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
			shared.LogError(logInfo, err)
			log.Fatal("Fatal error, couldn't import an account")
		}

		intPort, err := strconv.Atoi(nodeConfig.HTTPPort)
		if err != nil {
			log.Fatal(accCreateFatalMessage)
		}

		device, err := shared.ForwardPort(intPort)
		if err != nil {
			shared.LogError(logInfo, err)
			log.Fatal(accCreateFatalMessage)
		}

		defer device.Clear(uint16(intPort))

		server.Start(accountStr, nodeConfig.HTTPPort)
	},
}

func init() {
	accountCmd.AddCommand(accountImportCmd)
}
