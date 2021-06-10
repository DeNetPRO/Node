package cmd

import (
	"dfile-secondary-node/account"
	blockchainprovider "dfile-secondary-node/blockchain_provider"
	"dfile-secondary-node/config"
	"dfile-secondary-node/shared"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

const confUpdateFatalMessage = "Fatal error while configuration update"

// accountListCmd represents the list command
var configUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "updates your account configuration",
	Long:  "updates your account configuration",
	Run: func(cmd *cobra.Command, args []string) {

		accounts := account.List()

		if len(accounts) > 1 {
			fmt.Println("Which account configuration would you like to change?")
			for i, a := range accounts {
				fmt.Println(i+1, a)
			}
		}

		var address string
		var password string

		for {

			if len(accounts) == 1 {
				address = accounts[0]
			} else {
				byteAddress, err := shared.ReadFromConsole()
				if err != nil {
					log.Fatal(confUpdateFatalMessage)
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
				log.Fatal(confUpdateFatalMessage)
			}
			password = string(bytePassword)
			if strings.Trim(password, " ") == "" {
				fmt.Println("Empty string can't be used as a password. Please enter passwords again")
				continue
			}

			err = account.Login(address, password)
			if err != nil {
				log.Fatal("Wrong password")
			}

			break
		}

		pathToConfigDir := filepath.Join(shared.AccsDirPath, address, shared.ConfDirName)
		pathToConfigFile := filepath.Join(pathToConfigDir, "config.json")

		var dFileConf config.SecondaryNodeConfig

		confFile, err := os.OpenFile(pathToConfigFile, os.O_RDWR, 0700)
		if err != nil {
			log.Fatal(confUpdateFatalMessage)
		}
		defer confFile.Close()

		fileBytes, err := io.ReadAll(confFile)
		if err != nil {
			log.Fatal(confUpdateFatalMessage)
		}

		err = json.Unmarshal(fileBytes, &dFileConf)
		if err != nil {
			log.Fatal(confUpdateFatalMessage)
		}

		stateBefore := dFileConf

		fmt.Println("Please enter disk space for usage in GB (should be positive number), or just press enter button to skip")

		err = config.SetStorageLimit(pathToConfigDir, config.State.Update, &dFileConf)
		if err != nil {
			log.Fatal(confUpdateFatalMessage)
		}

		fmt.Println("Please enter new ip address, or just press enter button to skip")

		splitIPAddr, err := config.SetIpAddr(&dFileConf, config.State.Update)
		if err != nil {
			log.Fatal(confUpdateFatalMessage)
		}

		fmt.Println("Please enter new http port number, or just press enter button to skip")

		err = config.SetPort(&dFileConf, config.State.Update)
		if err != nil {
			log.Fatal(confUpdateFatalMessage)
		}

		if stateBefore.IpAddress == dFileConf.IpAddress &&
			stateBefore.HTTPPort == dFileConf.HTTPPort &&
			stateBefore.StorageLimit == dFileConf.StorageLimit {
			fmt.Println("Nothing was changed")
			return
		}

		if stateBefore.IpAddress != dFileConf.IpAddress || stateBefore.HTTPPort != dFileConf.HTTPPort {
			blockchainprovider.UpdateNodeInfo(common.HexToAddress(address), password, dFileConf.HTTPPort, splitIPAddr)
		}

		confJSON, err := json.Marshal(dFileConf)
		if err != nil {
			log.Fatal(confUpdateFatalMessage)
		}

		err = confFile.Truncate(0)
		if err != nil {
			log.Fatal(confUpdateFatalMessage)
		}

		_, err = confFile.Seek(0, 0)
		if err != nil {
			log.Fatal(confUpdateFatalMessage)
		}

		_, err = confFile.Write(confJSON)
		if err != nil {
			log.Fatal(confUpdateFatalMessage)
		}

		confFile.Sync()

		fmt.Println("Config file is updated successfully")

	},
}

func init() {
	configCmd.AddCommand(configUpdateCmd)
}
