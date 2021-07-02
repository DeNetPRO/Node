package cmd

import (
	"context"
	"dfile-secondary-node/account"
	blockchainprovider "dfile-secondary-node/blockchain_provider"
	"dfile-secondary-node/cleaner"
	"dfile-secondary-node/config"
	"dfile-secondary-node/paths"
	"dfile-secondary-node/server"
	"dfile-secondary-node/shared"
	"dfile-secondary-node/upnp"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
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

		pathToConfigDir := filepath.Join(paths.AccsDirPath, etherAccount.Address.String(), paths.ConfDirName)

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
			confFile, err := os.OpenFile(pathToConfigFile, os.O_RDWR, 0700)
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

			if dFileConf.StorageLimit <= 0 {
				log.Fatal(accLoginFatalError)
			}

			if upnp.InternetDevice != nil {

				ip, err := upnp.InternetDevice.ExternalIP()
				if err != nil {
					shared.LogError(logInfo, shared.GetDetailedError(err))
				}

				if dFileConf.IpAddress != ip {

					fmt.Println("Updating public ip info...")

					splitIPAddr := strings.Split(ip, ".")

					ctx, _ := context.WithTimeout(context.Background(), time.Minute)

					err = blockchainprovider.UpdateNodeInfo(ctx, etherAccount.Address, password, dFileConf.HTTPPort, splitIPAddr)
					if err != nil {
						shared.LogError(logInfo, shared.GetDetailedError(err))
						log.Fatal(ipUpdateFatalError)
					}

					dFileConf.IpAddress = ip

					err = config.SaveAndClose(confFile, dFileConf) // we dont't use mutex because race condition while login is impossible
					if err != nil {
						shared.LogError(logInfo, shared.GetDetailedError(err))
						log.Fatal(ipUpdateFatalError)
					}

				}
			}

			shared.SendLogs = dFileConf.AgreeSendLogs
		}

		fmt.Println("Logged in")

		go blockchainprovider.StartMining(password)

		go cleaner.Start()

		server.Start(etherAccount.Address.String(), dFileConf.HTTPPort)
	},
}

func init() {
	accountCmd.AddCommand(accountLoginCmd)
}
