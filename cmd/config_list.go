package cmd

import (
	"dfile-secondary-node/config"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"strconv"
)

// configListCmd represents the configList command
var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "display all startup configs",
	Long:  `display all startup configs`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("configList called")
		configs, err := config.GetConfigsList()
		if err != nil {
			log.Fatal(err)
		}
		i := 1
		for k, _ := range configs {
			fmt.Println(strconv.Itoa(i) + " " + k)
			i++
		}
	},
}

func init() {
	configCmd.AddCommand(configListCmd)

}
