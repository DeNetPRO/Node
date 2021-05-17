package cmd

import (
	"bufio"
	"dfile-secondary-node/account"
	"dfile-secondary-node/server"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// accountListCmd represents the list command
var accountCheckCmd = &cobra.Command{
	Use:   "list",
	Short: "display addresses of all blockchain accounts",
	Long:  `display addresses of all blockchain accounts`,
	Run: func(cmd *cobra.Command, args []string) {

		allMatch := false

		fmt.Println("Please enter account address you want to log in:")

		var address string
		var password string

		for !allMatch {
			reader := bufio.NewReader(os.Stdin)
			byteAddress, _, err := reader.ReadLine()
			if err != nil {
				log.Fatal("Fatal error while account log in.")
			}

			address = string(byteAddress)

			accounts := account.GetAllAccounts()

			addressMatches := false

			for _, a := range accounts {
				if a == address {
					addressMatches = true
				}
			}

			if !addressMatches {
				fmt.Println("There is no such account address:")
				for i, a := range accounts {
					fmt.Println(i+1, a)
				}
				continue
			}

			fmt.Println("Please enter your password:")

			bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				log.Fatal("Fatal error while account log in.")
			}
			password = string(bytePassword)
			if strings.Trim(password, " ") == "" {
				fmt.Println("Empty string can't be used as a password. Please, enter passwords again.")
				continue
			}

			allMatch = true
		}

		err := account.CheckPassword(password, address)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Wrong password")
			return
		}

		fmt.Println("Success")

		server.Start(address, "48658")

	},
}

func init() {
	accountCmd.AddCommand(accountCheckCmd)
}
