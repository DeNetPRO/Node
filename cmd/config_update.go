package cmd

import (
	"context"

	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/account"
	blckChain "git.denetwork.xyz/DeNet/dfile-secondary-node/blockchain_provider"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/config"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/logger"
	nodeFile "git.denetwork.xyz/DeNet/dfile-secondary-node/node_file"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/paths"
	"github.com/spf13/cobra"
)

const confUpdateFatalMessage = "Fatal error while configuration update"

// ConfigUpdateCmd is execited when "update" flag is passed after "config" flag and id used for updating
// user's configuration file.
var configUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "updates your account configuration",
	Long:  "updates your account configuration",
	Run: func(cmd *cobra.Command, args []string) {
		const location = "configUpdateCmd->"
		accounts := account.List()

		if len(accounts) > 1 {
			fmt.Println("Which account configuration would you like to change?")
			for i, a := range accounts {
				fmt.Println(i+1, a)
			}
		}

		nodeAccount, password, err := account.ValidateUser()
		if err != nil {
			logger.Log(logger.CreateDetails(location, err))
			log.Fatal(confUpdateFatalMessage)
		}

		pathToConfigDir := filepath.Join(paths.AccsDirPath, nodeAccount.Address.String(), paths.ConfDirName)
		pathToConfigFile := filepath.Join(pathToConfigDir, paths.ConfFileName)

		var nodeConfig config.NodeConfig

		confFile, fileBytes, err := nodeFile.Read(pathToConfigFile)
		if err != nil {
			logger.Log(logger.CreateDetails(location, err))
			log.Fatal(confUpdateFatalMessage)
		}
		defer confFile.Close()

		err = json.Unmarshal(fileBytes, &nodeConfig)
		if err != nil {
			logger.Log(logger.CreateDetails(location, err))
			log.Fatal(confUpdateFatalMessage)
		}

		stateBefore := nodeConfig

		network, err := config.SelectNetwork()
		if err != nil {
			logger.Log(logger.CreateDetails(location, err))
			log.Fatal(confUpdateFatalMessage)
		}

		nodeConfig.Network = network

		fmt.Println("Please enter disk space for usage in GB (should be positive number), or just press enter button to skip")

		err = config.SetStorageLimit(&nodeConfig, config.UpdateStatus)
		if err != nil {
			logger.Log(logger.CreateDetails(location, err))
			log.Fatal(confUpdateFatalMessage)
		}

		fmt.Println("Please enter new ip address, or just press enter button to skip")

		err = config.SetIpAddr(&nodeConfig, config.UpdateStatus)
		if err != nil {
			logger.Log(logger.CreateDetails(location, err))
			log.Fatal(confUpdateFatalMessage)
		}

		fmt.Println("Please enter new http port number, or just press enter button to skip")

		err = config.SetPort(&nodeConfig, config.UpdateStatus)
		if err != nil {
			logger.Log(logger.CreateDetails(location, err))
			log.Fatal(confUpdateFatalMessage)
		}

		fmt.Println("Do you want to send bug reports to developers? [y/n] (or just press enter button to skip)")

		err = config.ChangeAgreeSendLogs(&nodeConfig, config.UpdateStatus)
		if err != nil {
			logger.Log(logger.CreateDetails(location, err))
			log.Fatal(confUpdateFatalMessage)
		}

		if stateBefore.IpAddress == nodeConfig.IpAddress &&
			stateBefore.Network == nodeConfig.Network &&
			stateBefore.HTTPPort == nodeConfig.HTTPPort &&
			stateBefore.StorageLimit == nodeConfig.StorageLimit &&
			stateBefore.SendBugReports == nodeConfig.SendBugReports {
			fmt.Println("Nothing was changed")
			return
		}

		if stateBefore.IpAddress != nodeConfig.IpAddress || stateBefore.HTTPPort != nodeConfig.HTTPPort {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)

			defer cancel()

			err := blckChain.UpdateNodeInfo(ctx, nodeAccount.Address, password, nodeConfig.IpAddress, nodeConfig.HTTPPort)
			if err != nil {
				logger.Log(logger.CreateDetails(location, err))
				log.Fatal(confUpdateFatalMessage)
			}
		}

		err = config.Save(confFile, nodeConfig) // we dont't use mutex because race condition while config update is impossible
		if err != nil {
			logger.Log(logger.CreateDetails(location, err))
			log.Fatal(confUpdateFatalMessage)
		}

		fmt.Println("Config file is updated successfully")
	},
}

func init() {
	configCmd.AddCommand(configUpdateCmd)
}
