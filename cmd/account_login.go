package cmd

import (
	"context"

	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"git.denetwork.xyz/dfile/dfile-secondary-node/account"
	blockchainprovider "git.denetwork.xyz/dfile/dfile-secondary-node/blockchain_provider"
	"git.denetwork.xyz/dfile/dfile-secondary-node/cleaner"
	"git.denetwork.xyz/dfile/dfile-secondary-node/config"
	"git.denetwork.xyz/dfile/dfile-secondary-node/logger"
	"git.denetwork.xyz/dfile/dfile-secondary-node/paths"
	"git.denetwork.xyz/dfile/dfile-secondary-node/server"
	"git.denetwork.xyz/dfile/dfile-secondary-node/shared"
	"git.denetwork.xyz/dfile/dfile-secondary-node/upnp"
	"github.com/spf13/cobra"
)

const accLoginFatalError = "Error while account log in"
const ipUpdateFatalError = "Couldn't update public ip info"

// accountListCmd represents the list command
var accountLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "log in a blockchain accounts",
	Long:  "log in a blockchain accounts",
	Run: func(cmd *cobra.Command, args []string) {
		const logInfo = "accountLoginCmd->"
		etherAccount, password, err := account.ValidateUser()
		if err != nil {
			logger.Log(logger.CreateDetails(logInfo, err))
			log.Fatal(accLoginFatalError)
		}

		pathToConfigDir := filepath.Join(paths.AccsDirPath, etherAccount.Address.String(), paths.ConfDirName)

		var nodeConfig config.SecondaryNodeConfig

		pathToConfigFile := filepath.Join(pathToConfigDir, "config.json")

		stat, err := os.Stat(pathToConfigFile)
		err = shared.CheckStatErr(err)
		if err != nil {
			log.Fatal(accLoginFatalError)
		}

		if stat == nil {
			nodeConfig, err = config.Create(etherAccount.Address.String(), password)
			if err != nil {
				logger.Log(logger.CreateDetails(logInfo, err))
				log.Fatal("couldn't create config file")
			}
		} else {
			confFile, err := os.OpenFile(pathToConfigFile, os.O_RDWR, 0700)
			if err != nil {
				logger.Log(logger.CreateDetails(logInfo, err))
				log.Fatal("couldn't open config file")
			}
			defer confFile.Close()

			fileBytes, err := io.ReadAll(confFile)
			if err != nil {
				logger.Log(logger.CreateDetails(logInfo, err))
				log.Fatal("couldn't read config file")
			}

			err = json.Unmarshal(fileBytes, &nodeConfig)
			if err != nil {
				logger.Log(logger.CreateDetails(logInfo, err))
				log.Fatal("couldn't read config file")
			}

			if nodeConfig.StorageLimit <= 0 {
				log.Fatal(accLoginFatalError)
			}

			if upnp.InternetDevice != nil {
				ip, err := upnp.InternetDevice.PublicIP()
				if err != nil {
					logger.Log(logger.CreateDetails(logInfo, err))
				}

				if nodeConfig.IpAddress != ip {

					fmt.Println("Updating public ip info...")

					splitIPAddr := strings.Split(ip, ".")

					ctx, _ := context.WithTimeout(context.Background(), time.Minute)

					err = blockchainprovider.UpdateNodeInfo(ctx, etherAccount.Address, password, nodeConfig.HTTPPort, splitIPAddr)
					if err != nil {
						logger.Log(logger.CreateDetails(logInfo, err))
						log.Fatal(ipUpdateFatalError)
					}

					nodeConfig.IpAddress = ip

					err = config.Save(confFile, nodeConfig) // we dont't use mutex because race condition while login is impossible
					if err != nil {
						logger.Log(logger.CreateDetails(logInfo, err))
						log.Fatal(ipUpdateFatalError)
					}
				}
			}

			logger.SendLogs = nodeConfig.AgreeSendLogs
		}

		account.NodeIpAddr = fmt.Sprint(nodeConfig.IpAddress, ":", nodeConfig.HTTPPort)

		fmt.Println("Logged in")

		go blockchainprovider.StartMining(password)

		go cleaner.Start()

		server.Start(nodeConfig.HTTPPort)
	},
}

func init() {
	accountCmd.AddCommand(accountLoginCmd)
}
