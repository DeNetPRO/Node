package cmd

import (
	"dfile-secondary-node/account"
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

// accountCreateCmd represents the accountCreate command
var accountCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create a new blockchain account",
	Long: `create a new blockchain account`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called")
		accountStr, err := account.CreateAccount("password")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(accountStr)
	},
}

func init() {
	accountCmd.AddCommand(accountCreateCmd)
}
