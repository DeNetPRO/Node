package cmd

import (
	"context"

	"encoding/json"
	"fmt"
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
		const logLoc = "accountLoginCmd->"
		etherAccount, password, err := account.ValidateUser()
		if err != nil {
			logger.Log(logger.CreateDetails(logLoc, err))
			log.Fatal(accLoginFatalError)
		}

		pathToConfigDir := filepath.Join(paths.AccsDirPath, etherAccount.Address.String(), paths.ConfDirName)

		var nodeConfig config.SecondaryNodeConfig

		pathToConfigFile := filepath.Join(pathToConfigDir, paths.ConfFileName)

		stat, err := os.Stat(pathToConfigFile)
		err = shared.CheckStatErr(err)
		if err != nil {
			logger.Log(logger.CreateDetails(logLoc, err))
			log.Fatal(accLoginFatalError)
		}

		if stat == nil {
			nodeConfig, err = config.Create(etherAccount.Address.String(), password)
			if err != nil {
				logger.Log(logger.CreateDetails(logLoc, err))
				log.Fatal("couldn't create config file")
			}
		} else {
			confFile, fileBytes, err := shared.ReadFile(pathToConfigFile)
			if err != nil {
				logger.Log(logger.CreateDetails(logLoc, err))
				log.Fatal("couldn't open config file")
			}
			defer confFile.Close()

			err = json.Unmarshal(fileBytes, &nodeConfig)
			if err != nil {
				logger.Log(logger.CreateDetails(logLoc, err))
				log.Fatal("couldn't read config file")
			}

			if nodeConfig.StorageLimit <= 0 {
				log.Fatal(accLoginFatalError)
			}

			if upnp.InternetDevice != nil {
				ip, err := upnp.InternetDevice.PublicIP()
				if err != nil {
					logger.Log(logger.CreateDetails(logLoc, err))
				}

				if nodeConfig.IpAddress != ip {

					fmt.Println("Updating public ip info...")

					splitIPAddr := strings.Split(ip, ".")

					ctx, _ := context.WithTimeout(context.Background(), time.Minute)

					err = blockchainprovider.UpdateNodeInfo(ctx, etherAccount.Address, password, nodeConfig.HTTPPort, splitIPAddr)
					if err != nil {
						logger.Log(logger.CreateDetails(logLoc, err))
						log.Fatal(ipUpdateFatalError)
					}

					nodeConfig.IpAddress = ip

					err = config.Save(confFile, nodeConfig) // we dont't use mutex because race condition while login is impossible
					if err != nil {
						logger.Log(logger.CreateDetails(logLoc, err))
						log.Fatal(ipUpdateFatalError)
					}
				}
			}

			logger.SendLogs = nodeConfig.AgreeSendLogs
		}

		account.IpAddr = fmt.Sprint(nodeConfig.IpAddress, ":", nodeConfig.HTTPPort)

		rating, err := os.Stat(paths.RatingFilePath)
		err = shared.CheckStatErr(err)
		if err != nil {
			logger.Log(logger.CreateDetails(logLoc, err))
			log.Fatal(accLoginFatalError)
		}

		if rating == nil {
			file, err := os.Create(paths.RatingFilePath)
			if err != nil {
				logger.Log(logger.CreateDetails(logLoc, err))
				log.Fatal(accLoginFatalError)
			}

			ratingInfo := shared.NewRatingInfo()

			if nodeConfig.IpAddress == "46.101.202.151" {
				ratingInfo.Rating = 100
				ratingInfo.ConnectedNodes["68.183.215.241:55050"] = 0
				ratingInfo.ConnectedNodes["157.230.98.89:55050"] = 0
				ratingInfo.NumberOfAuthorityConn = 99
			}

			err = shared.WriteFile(file, ratingInfo)
			if err != nil {
				logger.Log(logger.CreateDetails(logLoc, err))
				log.Fatal(accLoginFatalError)
			}

			file.Close()
		}

		fmt.Println("Logged in")

		go blockchainprovider.StartMining(password)

		go cleaner.Start()

		server.Start(nodeConfig.HTTPPort)
	},
}

func init() {
	accountCmd.AddCommand(accountLoginCmd)
}
