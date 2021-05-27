package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// accountCmd represents the account command
var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "account is a command that lets you to manage accounts in the DFile blockchain network",
	Long:  "account is a command that lets you to manage accounts in the DFile blockchain network",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`account:
		account list: lists all accounts that were created
		account create: creates a new account
		account login: asks for account address you want to log in and logs you in`)
	},
}

func init() {
	rootCmd.AddCommand(accountCmd)
}
