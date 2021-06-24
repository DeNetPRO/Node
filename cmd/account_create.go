package cmd

import (
	"dfile-secondary-node/account"
	"dfile-secondary-node/server"
	"dfile-secondary-node/shared"
	"fmt"
	"log"
	"os"
	"strconv"
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
		const logInfo = "accountCreateCmd->"
		var password1, password2 string

		fmt.Println("Password is required for account creation. It can't be restored, please save it in a safe place.")
		fmt.Println("Please enter your new password: ")

		for {
			bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				shared.LogError(logInfo, err)
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
				shared.LogError(logInfo, err)
				log.Println(accCreateFatalMessage)
			}

			password2 = string(bytePassword)

			if password1 == password2 {
				break
			}

			fmt.Println("Passwords do not match. Please enter passwords again.")
		}

		password := shared.GetHashPassword(password1)
		password1 = ""

		accountStr, nodeConfig, err := account.Create(password)
		if err != nil {
			shared.LogError(logInfo, err)
			log.Fatal(accCreateFatalMessage)
		}

		intPort, err := strconv.Atoi(nodeConfig.HTTPPort)
		if err != nil {
			shared.LogError(logInfo, err)
			log.Fatal(accCreateFatalMessage)
		}

		fmt.Println("forward port")
		if err := shared.InternetDevice.Forward(uint16(intPort), "node"); err != nil {
			shared.LogError(logInfo, err)
			log.Fatal(accCreateFatalMessage)
		}

		defer shared.InternetDevice.Clear(uint16(intPort))

		server.Start(accountStr, nodeConfig.HTTPPort)
	},
}

func init() {
	accountCmd.AddCommand(accountCreateCmd)
}
