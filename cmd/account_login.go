package cmd

import (
	"context"
	"dfile-secondary-node/account"
	bcProvider "dfile-secondary-node/blockchain_provider"
	blockchainprovider "dfile-secondary-node/blockchain_provider"
	"dfile-secondary-node/config"
	"dfile-secondary-node/server"
	"dfile-secondary-node/shared"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

const accLoginFatalError = "Fatal error while account log in"
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
			log.Fatal(accLoginFatalError)
		}

		pathToConfigDir := filepath.Join(shared.AccsDirPath, etherAccount.Address.String(), shared.ConfDirName)

		var dFileConf config.SecondaryNodeConfig

		pathToConfigFile := filepath.Join(pathToConfigDir, "config.json")

		stat, err := os.Stat(pathToConfigFile)
		err = shared.CheckStatErr(err)
		if err != nil {
			log.Fatal(accLoginFatalError)
		}

		if stat == nil {
			dFileConf, err = config.Create(etherAccount.Address.String(), password)
			if err != nil {
				shared.LogError(logInfo, shared.GetDetailedError(err))
				log.Fatal("couldn't create config file")
			}
		} else {
			confFile, err := os.Open(pathToConfigFile)
			if err != nil {
				shared.LogError(logInfo, shared.GetDetailedError(err))
				log.Fatal("couldn't open config file")
			}
			defer confFile.Close()

			fileBytes, err := io.ReadAll(confFile)
			if err != nil {
				shared.LogError(logInfo, shared.GetDetailedError(err))
				log.Fatal("couldn't read config file")
			}

			err = json.Unmarshal(fileBytes, &dFileConf)
			if err != nil {
				shared.LogError(logInfo, shared.GetDetailedError(err))
				log.Fatal("couldn't read config file")
			}

			ip, err := shared.InternetDevice.ExternalIP()
			if err != nil {
				shared.LogError(logInfo, shared.GetDetailedError(err))
				log.Fatal("couldn't check public ip")
			}

			if dFileConf.StorageLimit <= 0 {
				log.Fatal(accLoginFatalError)
			}

			if dFileConf.IpAddress != ip {

				fmt.Println("Updating public ip info...")

				splitIPAddr, err := config.SetIpAddr(&dFileConf, config.State.Update)
				if err != nil {
					shared.LogError(logInfo, shared.GetDetailedError(err))
					log.Fatal(ipUpdateFatalError)
				}

				ctx, _ := context.WithTimeout(context.Background(), time.Minute)

				err = blockchainprovider.UpdateNodeInfo(ctx, etherAccount.Address, password, dFileConf.HTTPPort, splitIPAddr)
				if err != nil {
					shared.LogError(logInfo, shared.GetDetailedError(err))
					log.Fatal(ipUpdateFatalError)
				}

				confJSON, err := json.Marshal(dFileConf)
				if err != nil {
					shared.LogError(logInfo, shared.GetDetailedError(err))
					log.Fatal(ipUpdateFatalError)
				}

				err = confFile.Truncate(0)
				if err != nil {
					shared.LogError(logInfo, shared.GetDetailedError(err))
					log.Fatal(ipUpdateFatalError)
				}

				_, err = confFile.Seek(0, 0)
				if err != nil {
					shared.LogError(logInfo, shared.GetDetailedError(err))
					log.Fatal(ipUpdateFatalError)
				}

				_, err = confFile.Write(confJSON)
				if err != nil {
					shared.LogError(logInfo, shared.GetDetailedError(err))
					log.Fatal(ipUpdateFatalError)
				}

			}

			shared.SendLogs = dFileConf.AgreeSendLogs
		}

		fmt.Println("Logged in")

		go bcProvider.StartMining(password)

		server.Start(etherAccount.Address.String(), dFileConf.HTTPPort)
	},
}

func init() {
	accountCmd.AddCommand(accountLoginCmd)
}
