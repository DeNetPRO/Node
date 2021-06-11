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

	"github.com/spf13/cobra"
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

		etherAccount, password, err := account.ValidateUser()
		if err != nil {
			shared.LogError(err.Error())
			log.Fatal(confUpdateFatalMessage)
		}

		pathToConfigDir := filepath.Join(shared.AccsDirPath, etherAccount.Address.String(), shared.ConfDirName)
		pathToConfigFile := filepath.Join(pathToConfigDir, "config.json")

		var dFileConf config.SecondaryNodeConfig

		confFile, err := os.OpenFile(pathToConfigFile, os.O_RDWR, 0700)
		if err != nil {
			shared.LogError(err.Error())
			log.Fatal(confUpdateFatalMessage)
		}
		defer confFile.Close()

		fileBytes, err := io.ReadAll(confFile)
		if err != nil {
			shared.LogError(err.Error())
			log.Fatal(confUpdateFatalMessage)
		}

		err = json.Unmarshal(fileBytes, &dFileConf)
		if err != nil {
			shared.LogError(err.Error())
			log.Fatal(confUpdateFatalMessage)
		}

		stateBefore := dFileConf

		fmt.Println("Please enter disk space for usage in GB (should be positive number), or just press enter button to skip")

		err = config.SetStorageLimit(pathToConfigDir, config.State.Update, &dFileConf)
		if err != nil {
			shared.LogError(err.Error())
			log.Fatal(confUpdateFatalMessage)
		}

		fmt.Println("Please enter new ip address, or just press enter button to skip")

		splitIPAddr, err := config.SetIpAddr(&dFileConf, config.State.Update)
		if err != nil {
			shared.LogError(err.Error())
			log.Fatal(confUpdateFatalMessage)
		}

		fmt.Println("Please enter new http port number, or just press enter button to skip")

		err = config.SetPort(&dFileConf, config.State.Update)
		if err != nil {
			shared.LogError(err.Error())
			log.Fatal(confUpdateFatalMessage)
		}

		if stateBefore.IpAddress == dFileConf.IpAddress &&
			stateBefore.HTTPPort == dFileConf.HTTPPort &&
			stateBefore.StorageLimit == dFileConf.StorageLimit {
			fmt.Println("Nothing was changed")
			return
		}

		if stateBefore.IpAddress != dFileConf.IpAddress || stateBefore.HTTPPort != dFileConf.HTTPPort {
			blockchainprovider.UpdateNodeInfo(etherAccount.Address, password, dFileConf.HTTPPort, splitIPAddr)
		}

		confJSON, err := json.Marshal(dFileConf)
		if err != nil {
			shared.LogError(err.Error())
			log.Fatal(confUpdateFatalMessage)
		}

		err = confFile.Truncate(0)
		if err != nil {
			shared.LogError(err.Error())
			log.Fatal(confUpdateFatalMessage)
		}

		_, err = confFile.Seek(0, 0)
		if err != nil {
			shared.LogError(err.Error())
			log.Fatal(confUpdateFatalMessage)
		}

		_, err = confFile.Write(confJSON)
		if err != nil {
			shared.LogError(err.Error())
			log.Fatal(confUpdateFatalMessage)
		}

		confFile.Sync()

		fmt.Println("Config file is updated successfully")

	},
}

func init() {
	configCmd.AddCommand(configUpdateCmd)
}
