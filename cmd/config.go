package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "config is a command that allows you to manage the storage-node startup configuration files for mining",
	Long: `config is a command that allows you to manage the storage-node startup configuration files for mining`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("config called")
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
