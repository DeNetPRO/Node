package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ConfigCmd is executed when "config" flag is passed and is used to call a command for provided extra flag.
//If there's no extra flag, lists flags that can be passed along with "account" flag.
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "config is a command for managing configuration",
	Long:  "config is a command for managing configuration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`config:
		config update: updates your account configuration`)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
