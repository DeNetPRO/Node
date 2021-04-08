package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dfile-secondary-node",
	Short: "dfile-secondary-node is a CLI program that allows a user (miner) to connect to the DeNet decentralized network and mine DFile tokens.",
	Long: `dfile-secondary-node is a CLI program that allows a user (miner) 
	to connect to the DeNet decentralized network and mine DFile tokens.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("dfile-secondary-node is called!!!")
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
}


