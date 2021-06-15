package cmd

import (
	"dfile-secondary-node/account"
	bcProvider "dfile-secondary-node/blockchain_provider"
	"dfile-secondary-node/server"

	"dfile-secondary-node/config"
	"dfile-secondary-node/shared"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const accLoginFatalError = "Fatal error while account log in"

// accountListCmd represents the list command
var accountLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "log in a blockchain accounts",
	Long:  "log in a blockchain accounts",
	Run: func(cmd *cobra.Command, args []string) {
		const logInfo = "accountLoginCmd->"
		etherAccount, password, err := account.ValidateUser()
		if err != nil {
			shared.LogError(logInfo, err)
			log.Fatal(accLoginFatalError)
		}

		pathToConfigDir := filepath.Join(shared.AccsDirPath, etherAccount.Address.String(), shared.ConfDirName)

		var dFileConf config.SecondaryNodeConfig

		pathToConfigFile := filepath.Join(pathToConfigDir, "config.json")

		stat, err := os.Stat(pathToConfigFile)
		if err != nil {
			errPart := strings.Split(err.Error(), ":")

			if strings.Trim(errPart[1], " ") != "no such file or directory" {
				shared.LogError(logInfo, err)
				log.Fatal(accLoginFatalError)
			}
		}

		if stat == nil {
			dFileConf, err = config.Create(etherAccount.Address.String(), password)
			if err != nil {
				shared.LogError(logInfo, err)
				log.Fatal(accLoginFatalError)
			}
		} else {
			confFile, err := os.Open(pathToConfigFile)
			if err != nil {
				shared.LogError(logInfo, err)
				log.Fatal(accLoginFatalError)
			}
			defer confFile.Close()

			fileBytes, err := io.ReadAll(confFile)
			if err != nil {
				shared.LogError(logInfo, err)
				log.Fatal(accLoginFatalError)
			}

			err = json.Unmarshal(fileBytes, &dFileConf)
			if err != nil {
				shared.LogError(logInfo, err)
				log.Fatal(accLoginFatalError)
			}

			if dFileConf.StorageLimit <= 0 {
				shared.LogError(logInfo, err)
				log.Fatal(accLoginFatalError)
			}
		}

		fmt.Println("Logged in")

		go bcProvider.StartMining(password)

		server.Start(etherAccount.Address.String(), dFileConf.HTTPPort)
	},
}

func init() {
	accountCmd.AddCommand(accountLoginCmd)
}
