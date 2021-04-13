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

		var password_1, password_2 string
		password_match := false

		for !password_match {
			fmt.Println("Enter password: ")
			bytePassword, err := terminal.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				log.Fatal(err)
			}
			password_1 = string(bytePassword)

			fmt.Println("Enter password again: ")
			bytePassword, err = terminal.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				log.Fatal(err)
			}

			password_2 = string(bytePassword)

			if password_1 == password_2 {
				password_match = true
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
