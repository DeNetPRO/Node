package cmd

import (
	"dfile-secondary-node/account"
	"fmt"

	"github.com/spf13/cobra"
)

// accountListCmd represents the list command
var accountListCmd = &cobra.Command{
	Use:   "list",
	Short: "display addresses of all blockchain accounts",
	Long:  "display addresses of all blockchain accounts",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Accounts:")
		accounts := account.List()
		for i, a := range accounts {
			fmt.Println(i+1, a)
		}
	},
}

func init() {
	accountCmd.AddCommand(accountListCmd)
}
