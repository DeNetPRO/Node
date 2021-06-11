package cmd

import (
	"dfile-secondary-node/account"
	"dfile-secondary-node/config"
	"dfile-secondary-node/server"
	"dfile-secondary-node/shared"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

const accCreateFatalMessage = "Fatal error while creating an account"

// accountCreateCmd represents the accountCreate command
var accountCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create a new blockchain account",
	Long:  `create a new blockchain account`,
	Run: func(cmd *cobra.Command, args []string) {

		var password1, password2 string

		fmt.Println("Password is required for account creation. It can't be restored, please save it in a safe place.")
		fmt.Println("Please enter your new password: ")

		for {
			bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				shared.LogError(err.Error())
				log.Fatal(accCreateFatalMessage)
			}
			password1 = string(bytePassword)

			if strings.Trim(password1, " ") == "" {
				fmt.Println("Empty string can't be used as a password. Please enter passwords again.")
				continue
			}

			fmt.Println("Enter password again: ")
			bytePassword, err = term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				shared.LogError(err.Error())
				log.Fatal(accCreateFatalMessage)
			}

			password2 = string(bytePassword)

			if password1 == password2 {
				break
			}

			fmt.Println("Passwords do not match. Please enter passwords again.")
		}
		accountStr, err := account.Create(password1)
		if err != nil {
			shared.LogError(err.Error())
			log.Fatal(accCreateFatalMessage)
		}
		fmt.Println("Your new account's address is:", accountStr)

		secondaryNodeConfig, err := config.Create(accountStr, password1)
		if err != nil {
			shared.LogError(err.Error())
			log.Fatal(accCreateFatalMessage)
		}

		server.Start(accountStr, secondaryNodeConfig.HTTPPort)

	},
}

func init() {
	accountCmd.AddCommand(accountCreateCmd)
}
