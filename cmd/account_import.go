package cmd

import (
	"dfile-secondary-node/account"
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
		err := account.Import()
		if err != nil {
			shared.LogError(err.Error())
			log.Fatal("Fatal error, couldn't import an account")
		}

	},
}

func init() {
	accountCmd.AddCommand(accountImportCmd)
}
