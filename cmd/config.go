package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "config is a command that lets you change your account configuration",
	Long:  "config is a command that lets you change your account configuration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`config:
		config update: updates your account configuration`)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
