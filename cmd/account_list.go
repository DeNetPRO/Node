package cmd

import (
	"fmt"

	"git.denetwork.xyz/dfile/dfile-secondary-node/account"
	"github.com/spf13/cobra"
)

// AccountListCmd is executed when "list" flag is passed after "account" flag and is used for listing created/imported accounts.
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
