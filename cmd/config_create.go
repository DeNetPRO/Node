
package cmd

import (
	"dfile-secondary-node/common"
	"dfile-secondary-node/config"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"path/filepath"
)

// configCreateCmd represents the configCreate command
var configCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create new storage-node startup configuration file",
	Long: `create new storage-node startup configuration file`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("configCreate called")

		secondaryNodeConfig := config.SecondaryNodeConfig{}

		err := secondaryNodeConfig.Create()
		if err != nil{
			log.Fatal(err)
		}

		file, _ := json.MarshalIndent(secondaryNodeConfig, "", " ")

		configsDir, err := common.GetConfigsDirectory()
		if err != nil {
			log.Fatal(err)
		}


		_ = ioutil.WriteFile(filepath.Join(configsDir, secondaryNodeConfig.Name + ".json"), file, 0644)
	},
}

func init() {
	configCmd.AddCommand(configCreateCmd)
}
