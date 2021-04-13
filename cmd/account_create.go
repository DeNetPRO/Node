package cmd

import (
	"dfile-secondary-node/account"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
)

// accountCreateCmd represents the accountCreate command
var accountCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create a new blockchain account",
	Long: `create a new blockchain account`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called")

		var password1, password2 string
		passwordMatch := false

		for !passwordMatch {
			fmt.Println("Enter password: ")
			bytePassword, err := terminal.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				log.Fatal(err)
			}
			password1 = string(bytePassword)

			fmt.Println("Enter password again: ")
			bytePassword, err = terminal.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				log.Fatal(err)
			}

			password2 = string(bytePassword)

			if password1 == password2 {
				passwordMatch = true
			} else{
				fmt.Println("Passwords do not match! Please, enter passwords again")
			}

		}
		accountStr, err := account.CreateAccount("password")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Account created with address: ", accountStr)

	},
}

func init() {
	accountCmd.AddCommand(accountCreateCmd)
}
