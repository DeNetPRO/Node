package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

//AccountCmd is executed when "account" flag is passed to call a command for provided extra flag.
//If there's no extra flag, lists flags that can be passed along with "account" flag.
var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "account is a command for managing wallets",
	Long:  "account is a command for managing wallets",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`account:
		account list: lists all accounts that were created
		account create: creates a new account
		account import: imports your account by private key
		account login: asks for account address you want to log in and logs you in
		account key: discloses your private key`)
	},
}

func init() {
	rootCmd.AddCommand(accountCmd)
}
