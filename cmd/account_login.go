package cmd

import (
	"dfile-secondary-node/account"
	"dfile-secondary-node/config"
	"dfile-secondary-node/server"
	"dfile-secondary-node/shared"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
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
			byteAddress, err := shared.ReadFromConsole()
			if err != nil {
				log.Fatal("Fatal error while account log in.")
			}

			address = string(byteAddress)

			accounts := account.GetAllAccounts()

			addressMatches := shared.ContainsAccount(accounts, address)

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
			fmt.Println("Wrong password")
			return
		}

		confFiles := []string{}

		confDir := "config"

		confFilePath := filepath.Join(shared.AccDir, address, confDir)

		err = filepath.WalkDir(confFilePath,
			func(path string, info fs.DirEntry, err error) error {
				if err != nil {
					log.Fatal("Fatal error while account log in.")
				}

				if info.Name() != confDir {
					confFiles = append(confFiles, info.Name())
				}

				return nil
			})
		if err != nil {
			log.Fatal("Fatal error while account log in.")
		}

		var portAddr string

		if len(confFiles) == 0 {
			config, err := config.Create(address)
			if err != nil {
				log.Fatal("Fatal error while account log in.")
			}

			portAddr = config.HTTPPort
		} else {
			confFile, err := os.Open(filepath.Join(confFilePath, confFiles[0]))
			if err != nil {
				log.Fatal("Fatal error while account log in.")
			}

			fileBytes, err := io.ReadAll(confFile)
			if err != nil {
				log.Fatal("Fatal error while account log in.")
			}

			var conf config.SecondaryNodeConfig

			err = json.Unmarshal(fileBytes, &conf)
			if err != nil {
				log.Fatal("Fatal error while account log in.")
			}

			portAddr = conf.HTTPPort
		}

		fmt.Println("Success")

		server.Start(address, portAddr)

	},
}

func init() {
	accountCmd.AddCommand(accountCheckCmd)
}
