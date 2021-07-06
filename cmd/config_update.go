package cmd

import (
	"context"
	"dfile-secondary-node/account"
	blockchainprovider "dfile-secondary-node/blockchain_provider"
	"dfile-secondary-node/config"
	"dfile-secondary-node/logger"
	"dfile-secondary-node/paths"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

const confUpdateFatalMessage = "Fatal error while configuration update"

// accountListCmd represents the list command
var configUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "updates your account configuration",
	Long:  "updates your account configuration",
	Run: func(cmd *cobra.Command, args []string) {
		const logInfo = "configUpdateCmd->"
		accounts := account.List()

		if len(accounts) > 1 {
			fmt.Println("Which account configuration would you like to change?")
			for i, a := range accounts {
				fmt.Println(i+1, a)
			}
		}

		etherAccount, password, err := account.ValidateUser()
		if err != nil {
			logger.Log(logger.CreateDetails(logInfo, err))
			log.Fatal(confUpdateFatalMessage)
		}

		pathToConfigDir := filepath.Join(paths.AccsDirPath, etherAccount.Address.String(), paths.ConfDirName)
		pathToConfigFile := filepath.Join(pathToConfigDir, "config.json")

		var dFileConf config.SecondaryNodeConfig

		confFile, err := os.OpenFile(pathToConfigFile, os.O_RDWR, 0700)
		if err != nil {
			logger.Log(logger.CreateDetails(logInfo, err))
			log.Fatal(confUpdateFatalMessage)
		}
		defer confFile.Close()

		fileBytes, err := io.ReadAll(confFile)
		if err != nil {
			logger.Log(logger.CreateDetails(logInfo, err))
			log.Fatal(confUpdateFatalMessage)
		}

		err = json.Unmarshal(fileBytes, &dFileConf)
		if err != nil {
			logger.Log(logger.CreateDetails(logInfo, err))
			log.Fatal(confUpdateFatalMessage)
		}

		stateBefore := dFileConf

		fmt.Println("Please enter disk space for usage in GB (should be positive number), or just press enter button to skip")

		err = config.SetStorageLimit(pathToConfigDir, config.State.Update, &dFileConf)
		if err != nil {
			logger.Log(logger.CreateDetails(logInfo, err))
			log.Fatal(confUpdateFatalMessage)
		}

		fmt.Println("Please enter new ip address, or just press enter button to skip")

		splitIPAddr, err := config.SetIpAddr(&dFileConf, config.State.Update)
		if err != nil {
			logger.Log(logger.CreateDetails(logInfo, err))
			log.Fatal(confUpdateFatalMessage)
		}

		fmt.Println("Please enter new http port number, or just press enter button to skip")

		err = config.SetPort(&dFileConf, config.State.Update)
		if err != nil {
			logger.Log(logger.CreateDetails(logInfo, err))
			log.Fatal(confUpdateFatalMessage)
		}

		fmt.Println("Do you want to send bug reports to developers? [y/n] (or just press enter button to skip)")

		err = config.ChangeAgreeSendLogs(&dFileConf, config.State.Update)
		if err != nil {
			logger.Log(logger.CreateDetails(logInfo, err))
			log.Fatal(confUpdateFatalMessage)
		}

		if stateBefore.IpAddress == dFileConf.IpAddress &&
			stateBefore.HTTPPort == dFileConf.HTTPPort &&
			stateBefore.StorageLimit == dFileConf.StorageLimit &&
			stateBefore.AgreeSendLogs == dFileConf.AgreeSendLogs {
			fmt.Println("Nothing was changed")
			return
		}

		if stateBefore.IpAddress != dFileConf.IpAddress || stateBefore.HTTPPort != dFileConf.HTTPPort {
			ctx, _ := context.WithTimeout(context.Background(), time.Minute)

			err := blockchainprovider.UpdateNodeInfo(ctx, etherAccount.Address, password, dFileConf.HTTPPort, splitIPAddr)
			if err != nil {
				logger.Log(logger.CreateDetails(logInfo, err))
				log.Fatal(confUpdateFatalMessage)
			}
		}

		err = config.SaveAndClose(confFile, dFileConf) // we dont't use mutex because race condition while config update is impossible
		if err != nil {
			logger.Log(logger.CreateDetails(logInfo, err))
			log.Fatal(confUpdateFatalMessage)
		}

		fmt.Println("Config file is updated successfully")
	},
}

func init() {
	configCmd.AddCommand(configUpdateCmd)
}
