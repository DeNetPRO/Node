package cmd

import (
	"dfile-secondary-node/account"
	"dfile-secondary-node/config"
	"dfile-secondary-node/shared"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

const confUpdateFatalMessage = "Fatal error while configuration update"

// accountListCmd represents the list command
var configUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "updates your account configuration",
	Long:  "updates your account configuration",
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("Which account configuration would you like to change?")
		accounts := account.List()
		for i, a := range accounts {
			fmt.Println(i+1, a)
		}

		allMatch := false

		var address string
		var password string

		for !allMatch {
			byteAddress, err := shared.ReadFromConsole()
			if err != nil {
				log.Fatal(confUpdateFatalMessage)
			}

			address = string(byteAddress)

			addressMatches := shared.ContainsAccount(accounts, address)

			if !addressMatches {
				fmt.Println("There is no such account address:")
				for i, a := range accounts {
					fmt.Println(i+1, a)
				}
				continue
			}

			fmt.Println("Please enter your password:")

			bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				log.Fatal(confUpdateFatalMessage)
			}
			password = string(bytePassword)
			if strings.Trim(password, " ") == "" {
				fmt.Println("Empty string can't be used as a password. Please enter passwords again")
				continue
			}

			allMatch = true
		}

		confFilePath := filepath.Join(shared.AccsDirPath, address, shared.ConfDirName)

		confFiles := []string{}

		err := filepath.WalkDir(confFilePath,
			func(path string, info fs.DirEntry, err error) error {
				if err != nil {
					log.Fatal(confUpdateFatalMessage)
				}

				if info.Name() != shared.ConfDirName {
					confFiles = append(confFiles, info.Name())
				}

				return nil
			})
		if err != nil {
			log.Fatal(confUpdateFatalMessage)
		}

		if len(confFiles) == 0 {
			log.Fatal("Config file is not found in your account directory")
		}

		var dFileConf config.SecondaryNodeConfig

		pathToConfig := filepath.Join(shared.AccsDirPath, address, shared.ConfDirName)
		confFile, err := os.OpenFile(filepath.Join(pathToConfig, confFiles[0]), os.O_RDWR, 0700)
		if err != nil {
			log.Fatal(confUpdateFatalMessage)
		}
		defer confFile.Close()

		fileBytes, err := io.ReadAll(confFile)
		if err != nil {
			log.Fatal(confUpdateFatalMessage)
		}

		err = json.Unmarshal(fileBytes, &dFileConf)
		if err != nil {
			log.Fatal(confUpdateFatalMessage)
		}

		fmt.Println("You can change your http port number or storage limit")

		fmt.Println("Please enter disk space for usage in GB (should be positive number), or just press enter button to skip")

		spaceValueIsCorrect := false

		regNum := regexp.MustCompile(("[0-9]+"))

		for !spaceValueIsCorrect {

			availableSpace := shared.GetAvailableSpace(pathToConfig)

			fmt.Println("Available space:", availableSpace, "GB")

			space, err := shared.ReadFromConsole()
			if err != nil {
				log.Fatal(confUpdateFatalMessage)
			}

			if space == "" {
				spaceValueIsCorrect = true
				continue
			}

			match := regNum.MatchString(space)

			if !match {
				fmt.Println("Value is incorrect, please try again")
				continue
			}

			intSpace, err := strconv.Atoi(space)
			if err != nil {
				fmt.Println("Value is incorrect, please try again")
				continue
			}

			if intSpace < 0 || intSpace >= availableSpace {
				fmt.Println("Value is incorrect, please try again")
				continue
			}

			spaceValueIsCorrect = true
			dFileConf.StorageLimit = intSpace

		}

		fmt.Println("Please enter new http port address, or just press enter button to skip")

		portHTTPValueIsCorrect := false
		regPort := regexp.MustCompile("[0-9]+|")

		for !portHTTPValueIsCorrect {

			httpPort, err := shared.ReadFromConsole()
			if err != nil {
				log.Fatal(confUpdateFatalMessage)
			}

			if httpPort == "" {
				portHTTPValueIsCorrect = true
				continue
			}

			match := regPort.MatchString(httpPort)
			if !match {
				fmt.Println("Value is incorrect, please try again")
				continue

			}

			intHttpPort, err := strconv.Atoi(httpPort)
			if err != nil {
				fmt.Println("Value is incorrect, please try again")
				continue
			}
			if intHttpPort < 49152 || intHttpPort > 65535 {
				fmt.Println("Value is incorrect, please try again")
				continue

			}

			portHTTPValueIsCorrect = true
			dFileConf.HTTPPort = fmt.Sprint(intHttpPort)
		}

		confJSON, err := json.Marshal(dFileConf)
		if err != nil {
			log.Fatal(confUpdateFatalMessage)
		}

		err = confFile.Truncate(0)
		if err != nil {
			log.Fatal(confUpdateFatalMessage)
		}

		_, err = confFile.Seek(0, 0)
		if err != nil {
			log.Fatal(confUpdateFatalMessage)
		}

		_, err = confFile.Write(confJSON)
		if err != nil {
			log.Fatal(confUpdateFatalMessage)
		}

		confFile.Sync()

		fmt.Println("Config file is updated successfully")

	},
}

func init() {
	configCmd.AddCommand(configUpdateCmd)
}
