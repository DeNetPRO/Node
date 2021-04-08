package cmd

import (
	"dfile-secondary-node/account"
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

// accountListCmd represents the list command
var accountListCmd = &cobra.Command{
	Use:   "list",
	Short: "display addresses of all blockchain accounts",
	Long: `display addresses of all blockchain accounts`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list called")
		accounts, err := account.GetAllAccounts()
		if err != nil {
			log.Fatal(err)
		}
		for i, a := range accounts {
			fmt.Println(i+1, a)
		}
	},
}

func init() {
	accountCmd.AddCommand(accountListCmd)
}
