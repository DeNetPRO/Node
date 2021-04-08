package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// accountCmd represents the account command
var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "account is a command that allows you to manage accounts in the DFile blockchain network",
	Long:  `account is a command that allows you to manage accounts in the DFile blockchain network`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("account called")
	},
}

func init() {
	rootCmd.AddCommand(accountCmd)
}
