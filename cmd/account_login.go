package cmd

import (
	"context"
	"errors"

	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/account"
	blckChain "git.denetwork.xyz/DeNet/dfile-secondary-node/blockchain_provider"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/cleaner"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/config"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/errs"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/logger"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/networks"
	nodeFile "git.denetwork.xyz/DeNet/dfile-secondary-node/node_file"
	nodeTypes "git.denetwork.xyz/DeNet/dfile-secondary-node/node_types"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/rpcserver"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/paths"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/upnp"
	"github.com/spf13/cobra"
)

type previuosConfig struct {
	Address              string          `json:"nodeAddress"`
	IpAddress            string          `json:"ipAddress"`
	HTTPPort             string          `json:"portHTTP"`
	Network              string          `json:"network"`
	RPC                  string          `json:"rpc"`
	StorageLimit         int             `json:"storageLimit"`
	StoragePaths         []string        `json:"storagePaths"`
	UsedStorageSpace     int64           `json:"usedStorageSpace"`
	SendBugReports       bool            `json:"sendBugReports"`
	RegisteredInNetworks map[string]bool `json:"registeredInNetworks"`
}

const accLoginFatalError = "Error while account log in "
const ipUpdateFatalError = "Couldn't update public ip info"

// AccountLoginCmd is executed when "login" flag is passed after "account" flag and is used for logging in to an account.
var accountLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "logs in to your account",
	Long:  "logs in to your account",
	Run: func(cmd *cobra.Command, args []string) {
		const location = "accountLoginCmd->"
		nodeAccount, password, err := account.Unlock()
		if err != nil {
			logger.Log(logger.MarkLocation(location, err))
			log.Fatal(accLoginFatalError)
		}

		paths.SetConfigPath(nodeAccount.Address.String())

		var nodeConfig nodeTypes.Config

		stat, err := os.Stat(paths.List().ConfigDir)
		if err != nil {
			logger.Log(logger.MarkLocation(location, err))

			if !errors.Is(err, os.ErrNotExist) {
				log.Fatal(accLoginFatalError)
			}
		}

		configWasUpdated := false

		if stat == nil {
			nodeConfig, err = config.Create(nodeAccount.Address.String())
			if err != nil {
				logger.Log(logger.MarkLocation(location, err))
				log.Fatal("couldn't create config file")
			}
		} else {
			confFile, fileBytes, err := nodeFile.Read(paths.List().ConfigFile)
			if err != nil {
				logger.Log(logger.MarkLocation(location, err))
				log.Fatal("couldn't open config file")
			}
			defer confFile.Close()

			err = json.Unmarshal(fileBytes, &nodeConfig)
			if err != nil {

				var conf previuosConfig
				err = json.Unmarshal(fileBytes, &conf)
				if err != nil {
					logger.Log(logger.MarkLocation(location, err))
					log.Fatal("couldn't read config file")
				}

				nodeConfig = nodeTypes.Config{
					Address:   conf.Address,
					IpAddress: conf.IpAddress,
					HTTPPort:  conf.HTTPPort,
					Network:   conf.Network,
					RPC: map[string]string{"kovan": conf.RPC,
						"polygon": "https://polygon-rpc.com"},
					StorageLimit:         conf.StorageLimit,
					StoragePaths:         conf.StoragePaths,
					UsedStorageSpace:     conf.UsedStorageSpace,
					SendBugReports:       conf.SendBugReports,
					RegisteredInNetworks: conf.RegisteredInNetworks,
				}

				configWasUpdated = true

			}

			config.RPC = nodeConfig.RPC[nodeConfig.Network]

			if nodeConfig.StorageLimit <= 0 {
				log.Fatal(accLoginFatalError + " storage limit is " + fmt.Sprint(nodeConfig.StorageLimit))
			}

			err = networks.Set(nodeConfig.Network)
			if err != nil {
				log.Fatal(accLoginFatalError + errs.List().Network.Error())
			}

			registeredInNetwork, registrationExists := nodeConfig.RegisteredInNetworks[nodeConfig.Network]

			if !registeredInNetwork || !registrationExists {

				fmt.Println("registering node in", nodeConfig.Network)

				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)

				defer cancel()

				err = blckChain.RegisterNode(ctx, nodeAccount.Address, password, nodeConfig)
				if err != nil {
					logger.Log(logger.MarkLocation(location, err))
					log.Fatal(accLoginFatalError + ": couldn't register node in " + nodeConfig.Network)
				}

				nodeConfig.RegisteredInNetworks[nodeConfig.Network] = true
				configWasUpdated = true

			}

			if upnp.InternetDevice != nil {
				ip, err := upnp.InternetDevice.PublicIP()
				if err != nil {
					logger.Log(logger.MarkLocation(location, err))
				}

				if registeredInNetwork && nodeConfig.IpAddress != ip {

					fmt.Println("Updating public ip info...")

					ctx, cancel := context.WithTimeout(context.Background(), time.Minute)

					defer cancel()

					err = blckChain.UpdateNodeInfo(ctx, nodeAccount.Address, password, nodeConfig.IpAddress, nodeConfig.HTTPPort)
					if err != nil {
						logger.Log(logger.MarkLocation(location, err))
						log.Fatal(ipUpdateFatalError)
					}

					nodeConfig.IpAddress = ip
					configWasUpdated = true

				}
			}

			if configWasUpdated {
				err = config.Save(confFile, nodeConfig)
				if err != nil {
					logger.Log(logger.MarkLocation(location, err))
					log.Fatal(ipUpdateFatalError)
				}
			}

		}

		if len(nodeConfig.StoragePaths) == 0 {
			err := errors.New("path to storage is not specified in config")
			logger.Log(logger.MarkLocation(location, err))
			log.Fatal(err)
		}

		paths.SetStoragePaths(nodeConfig.StoragePaths)
		logger.SendReports = nodeConfig.SendBugReports

		fmt.Println("Logged in")

		go blckChain.StartMakingProofs(nodeAccount.Address, password, nodeConfig)

		go cleaner.Start()

		err = rpcserver.Start(nodeConfig.HTTPPort)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	accountCmd.AddCommand(accountLoginCmd)
}
