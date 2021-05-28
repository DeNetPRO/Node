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

		allMatch := false

		var address string
		var password string

		for !allMatch {

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
				fmt.Println("Empty string can't be used as a password. Please, enter passwords again.")
				continue
			}

			allMatch = true
		}

		err := account.Login(address, password)
		if err != nil {
			fmt.Println("Wrong password")
			return
		}

		confFiles := []string{}

		pathToConfig := filepath.Join(shared.AccDir, address, shared.ConfDir)

		err = filepath.WalkDir(pathToConfig,
			func(path string, info fs.DirEntry, err error) error {
				if err != nil {
					log.Fatal(accLoginFatalError)
				}

				if info.Name() != shared.ConfDir {
					confFiles = append(confFiles, info.Name())
				}

				return nil
			})
		if err != nil {
			log.Fatal(accLoginFatalError)
		}

		var dFileConf config.SecondaryNodeConfig

		if len(confFiles) == 0 {
			conf, err := config.Create(address)
			if err != nil {
				log.Fatal(accLoginFatalError)
			}

			dFileConf = conf
		} else {
			confFile, err := os.Open(filepath.Join(pathToConfig, confFiles[0]))
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

		fmt.Println("Success")

		// account.StartMining()
		// account.SendProof()

		server.Start(address, dFileConf.HTTPPort)
		// fmt.Println(dFileConf.HTTPPort)

	},
}

func init() {
	accountCmd.AddCommand(accountLoginCmd)
}
