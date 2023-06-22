package cmd

import (
	"context"

	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/DeNetPRO/src/account"
	blckChain "github.com/DeNetPRO/src/blockchain_provider"
	"github.com/DeNetPRO/src/config"
	"github.com/DeNetPRO/src/logger"
	nodeFile "github.com/DeNetPRO/src/node_file"
	nodeTypes "github.com/DeNetPRO/src/node_types"

	"github.com/DeNetPRO/src/paths"
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

		nodeAccount, password, err := account.Unlock()
		if err != nil {
			logger.Log(logger.MarkLocation(location, err))
			log.Fatal(confUpdateFatalMessage)
		}

		fmt.Println("Started configuration update")

		pathToConfigDir := filepath.Join(paths.List().AccsDir, nodeAccount.Address.String(), "config", "config.json") //TODO fix
		pathToConfigFile := filepath.Join(pathToConfigDir, paths.List().ConfigFile)

		var nodeConfig nodeTypes.Config

		confFile, fileBytes, err := nodeFile.Read(pathToConfigFile)
		if err != nil {
			logger.Log(logger.MarkLocation(location, err))
			log.Fatal(confUpdateFatalMessage)
		}
		defer confFile.Close()

		err = json.Unmarshal(fileBytes, &nodeConfig)
		if err != nil {
			logger.Log(logger.MarkLocation(location, err))
			log.Fatal(confUpdateFatalMessage)
		}

		config.RPC = nodeConfig.RPC[nodeConfig.Network]

		stateBefore := nodeConfig

		err = config.SetNetwork(&nodeConfig)
		if err != nil {
			logger.Log(logger.MarkLocation(location, err))
			log.Fatal(confUpdateFatalMessage)
		}

		fmt.Println("\nHow much GB are you going to share? (should be positive number), or just press enter button to skip")

		err = config.SetStorageLimit(&nodeConfig, config.Stats().Update)
		if err != nil {
			logger.Log(logger.MarkLocation(location, err))
			log.Fatal(confUpdateFatalMessage)
		}

		fmt.Println("\nPlease enter new ip address, or just press enter button to skip")

		err = config.SetIpAddr(&nodeConfig, config.Stats().Update)
		if err != nil {
			logger.Log(logger.MarkLocation(location, err))
			log.Fatal(confUpdateFatalMessage)
		}

		fmt.Println("\nPlease enter new http port number, or just press enter button to skip")

		err = config.SetPort(&nodeConfig, config.Stats().Update)
		if err != nil {
			logger.Log(logger.MarkLocation(location, err))
			log.Fatal(confUpdateFatalMessage)
		}

		fmt.Println("Do you want to send bug reports to developers? [y/n] (or just press enter button to skip)")

		err = config.SwitchReports(&nodeConfig, config.Stats().Update)
		if err != nil {
			logger.Log(logger.MarkLocation(location, err))
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
				logger.Log(logger.MarkLocation(location, err))
				log.Fatal(confUpdateFatalMessage)
			}
		}

		err = config.Save(confFile, nodeConfig) // we dont't use mutex because race condition while config update is impossible
		if err != nil {
			logger.Log(logger.MarkLocation(location, err))
			log.Fatal(confUpdateFatalMessage)
		}

		fmt.Println("Config file is updated successfully")
	},
}

func init() {
	configCmd.AddCommand(configUpdateCmd)
}
