package cmd

import (
	"dfile-secondary-node/account"
	bcProvider "dfile-secondary-node/blockchain_provider"
	"dfile-secondary-node/server"

	"dfile-secondary-node/config"
	"dfile-secondary-node/shared"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

const accLoginFatalError = "Fatal error while account log in"

// accountListCmd represents the list command
var accountLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "log in a blockchain accounts",
	Long:  "log in a blockchain accounts",
	Run: func(cmd *cobra.Command, args []string) {

		var accounts []string

		if len(args) == 1 {
			accounts = append(accounts, args[0])
		} else {
			accounts = account.List()
		}

		var address string
		var password string

		for {

			if len(args) == 1 {
				address = args[0]
			} else {
				byteAddress, err := shared.ReadFromConsole()
				if err != nil {
					log.Fatal(accLoginFatalError)
				}

				address = string(byteAddress)
			}

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
				log.Fatal(accLoginFatalError)
			}
			password = string(bytePassword)
			if strings.Trim(password, " ") == "" {
				fmt.Println("Empty string can't be used as a password. Please enter passwords again.")
				continue
			}

			break
		}

		err := account.Login(address, password)
		if err != nil {
			fmt.Println("Wrong password")
			return
		}

		pathToConfigDir := filepath.Join(shared.AccsDirPath, address, shared.ConfDirName)

		var dFileConf config.SecondaryNodeConfig

		pathToConfigFile := filepath.Join(pathToConfigDir, "config.json")

		stat, err := os.Stat(pathToConfigFile)

		if err != nil {
			errPart := strings.Split(err.Error(), ":")

			if strings.Trim(errPart[1], " ") != "no such file or directory" {
				log.Fatal(accLoginFatalError)
			}
		}

		if stat == nil {
			dFileConf, err = config.Create(address, password)
			if err != nil {
				log.Fatal(accLoginFatalError)
			}
		} else {
			confFile, err := os.Open(pathToConfigFile)
			if err != nil {
				log.Fatal(accLoginFatalError)
			}
			defer confFile.Close()

			fileBytes, err := io.ReadAll(confFile)
			if err != nil {
				log.Fatal(accLoginFatalError)
			}

			err = json.Unmarshal(fileBytes, &dFileConf)
			if err != nil {
				log.Fatal(accLoginFatalError)
			}

			if dFileConf.StorageLimit <= 0 {
				log.Fatal(accLoginFatalError)
			}
		}

		fmt.Println("Logged in")

		go bcProvider.StartMining(password)

		server.Start(address, dFileConf.HTTPPort)

	},
}

func init() {
	accountCmd.AddCommand(accountLoginCmd)
}
