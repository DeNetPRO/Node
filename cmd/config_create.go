package cmd

import (
	"dfile-secondary-node/config"
	"dfile-secondary-node/shared"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/spf13/cobra"
)

// configCreateCmd represents the configCreate command
var configCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create new storage-node startup configuration file",
	Long:  `create new storage-node startup configuration file`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("configCreate called")

		secondaryNodeConfig, err := config.Create("")
		if err != nil {
			log.Fatal(err)
		}

		file, _ := json.MarshalIndent(secondaryNodeConfig, "", " ")

		configsDir, err := shared.GetConfigsDirectory()
		if err != nil {
			log.Fatal(err)
		}

		_ = ioutil.WriteFile(filepath.Join(configsDir, secondaryNodeConfig.Name+".json"), file, 0644)
	},
}

func init() {
	configCmd.AddCommand(configCreateCmd)
}
