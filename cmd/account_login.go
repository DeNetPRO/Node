package cmd

import (
	"context"
	"errors"

	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/account"
	blckChain "git.denetwork.xyz/DeNet/dfile-secondary-node/blockchain_provider"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/cleaner"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/config"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/errs"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/logger"
	nodeFile "git.denetwork.xyz/DeNet/dfile-secondary-node/node_file"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/paths"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/server"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/upnp"
	"github.com/spf13/cobra"
)

const accLoginFatalError = "Error while account log in "
const ipUpdateFatalError = "Couldn't update public ip info"

// AccountLoginCmd is executed when "login" flag is passed after "account" flag and is used for logging in to an account.
var accountLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "logs in to your account",
	Long:  "logs in to your account",
	Run: func(cmd *cobra.Command, args []string) {
		const location = "accountLoginCmd->"
		nodeAccount, password, err := account.ValidateUser()
		if err != nil {
			logger.Log(logger.CreateDetails(location, err))
			log.Fatal(accLoginFatalError)
		}

		pathToConfigDir := filepath.Join(paths.AccsDirPath, nodeAccount.Address.String(), paths.ConfDirName)

		var nodeConfig config.NodeConfig

		pathToConfigFile := filepath.Join(pathToConfigDir, paths.ConfFileName)

		stat, err := os.Stat(pathToConfigFile)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			logger.Log(logger.CreateDetails(location, err))
			log.Fatal(accLoginFatalError)
		}

		if stat == nil {
			nodeConfig, err = config.Create(nodeAccount.Address.String())
			if err != nil {
				logger.Log(logger.CreateDetails(location, err))
				log.Fatal("couldn't create config file")
			}
		} else {
			confFile, fileBytes, err := nodeFile.Read(pathToConfigFile)
			if err != nil {
				logger.Log(logger.CreateDetails(location, err))
				log.Fatal("couldn't open config file")
			}
			defer confFile.Close()

			err = json.Unmarshal(fileBytes, &nodeConfig)
			if err != nil {
				logger.Log(logger.CreateDetails(location, err))
				log.Fatal("couldn't read config file")
			}

			if nodeConfig.StorageLimit <= 0 {
				log.Fatal(accLoginFatalError + " storage limit is " + fmt.Sprint(nodeConfig.StorageLimit))
			}

			_, supportedNet := blckChain.Networks[nodeConfig.Network]

			if !supportedNet {
				log.Fatal(accLoginFatalError + errs.NetworkCheck.Error())
			}

			blckChain.CurrentNetwork = nodeConfig.Network

			registeredInNetwork, registrationExists := nodeConfig.RegisteredInNetworks[nodeConfig.Network]

			if !registeredInNetwork || !registrationExists {

				fmt.Println("registering node in", nodeConfig.Network)

				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)

				defer cancel()

				err = blckChain.RegisterNode(ctx, nodeAccount.Address.String(), password, nodeConfig.IpAddress, nodeConfig.HTTPPort)
				if err != nil {
					logger.Log(logger.CreateDetails(location, err))
					log.Fatal(accLoginFatalError + ": couldn't register node in " + nodeConfig.Network)
				}

				nodeConfig.RegisteredInNetworks[nodeConfig.Network] = true

				err = config.Save(confFile, nodeConfig)
				if err != nil {
					logger.Log(logger.CreateDetails(location, err))
					log.Fatal(ipUpdateFatalError)
				}

			}

			if upnp.InternetDevice != nil {
				ip, err := upnp.InternetDevice.PublicIP()
				if err != nil {
					logger.Log(logger.CreateDetails(location, err))
				}

				if registeredInNetwork && nodeConfig.IpAddress != ip {

					fmt.Println("Updating public ip info...")

					ctx, cancel := context.WithTimeout(context.Background(), time.Minute)

					defer cancel()

					err = blckChain.UpdateNodeInfo(ctx, nodeAccount.Address, password, nodeConfig.IpAddress, nodeConfig.HTTPPort)
					if err != nil {
						logger.Log(logger.CreateDetails(location, err))
						log.Fatal(ipUpdateFatalError)
					}

					nodeConfig.IpAddress = ip

					err = config.Save(confFile, nodeConfig)
					if err != nil {
						logger.Log(logger.CreateDetails(location, err))
						log.Fatal(ipUpdateFatalError)
					}
				}
			}

		}

		if len(nodeConfig.StoragePaths) == 0 {
			err := errors.New("path to storage is not specified in config")
			logger.Log(logger.CreateDetails(location, err))
			log.Fatal(err)
		}

		paths.StoragePaths = nodeConfig.StoragePaths
		logger.SendLogs = nodeConfig.AgreeSendLogs

		fmt.Println("Logged in")

		go blckChain.StartMakingProofs(password)

		go cleaner.Start()

		server.Start(nodeConfig.HTTPPort)
	},
}

func init() {
	accountCmd.AddCommand(accountLoginCmd)
}
