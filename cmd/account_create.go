package cmd

import (
	"fmt"
	"log"
	"strings"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/account"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/cleaner"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/hash"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/logger"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/server"
	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
)

const accCreateFatalMessage = "Fatal error while creating an account"

// AccountCreateCmd is executed when "create" flag is passed after "account" flag and is used for crypto wallet creation.
var accountCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create a new blockchain account",
	Long:  `create a new blockchain account`,
	Run: func(cmd *cobra.Command, args []string) {
		const location = "accountCreateCmd->"
		var password1, password2 string

		fmt.Println("Password is required for account creation. It can't be restored, please save it in a safe place.")
		fmt.Println("Please enter your new password: ")

		for {
			bytePassword, err := gopass.GetPasswdMasked()
			if err != nil {
				logger.Log(logger.CreateDetails(location, err))
				log.Fatal(accCreateFatalMessage)
			}
			password1 = string(bytePassword)

			if strings.Trim(password1, " ") == "" {
				fmt.Println("Empty string can't be used as a password. Please enter passwords again.")
				continue
			}

			fmt.Println("Enter password again: ")
			bytePassword, err = gopass.GetPasswdMasked()
			if err != nil {
				logger.Log(logger.CreateDetails(location, err))
				log.Println(accCreateFatalMessage)
			}

			password2 = string(bytePassword)

			if password1 == password2 {
				break
			}

			fmt.Println("Passwords do not match. Please enter passwords again.")
		}

		password := hash.Password(password1)
		password1 = ""

		_, nodeConfig, err := account.Create(password)
		if err != nil {
			logger.Log(logger.CreateDetails(location, err))
			log.Fatal(accCreateFatalMessage)
		}

		go cleaner.Start()

		server.Start(nodeConfig.HTTPPort)
	},
}

func init() {
	accountCmd.AddCommand(accountCreateCmd)
}
