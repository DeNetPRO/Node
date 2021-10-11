package cmd

import (
	"git.denetwork.xyz/DeNet/dfile-secondary-node/account"
	"github.com/spf13/cobra"
)

// RootCmd is an entry point for executing CLI commands.
var rootCmd = &cobra.Command{
	Use:   "denet-node",
	Short: "denet-node provides unused space for users that need it",
	Long: `denet-node is a CLI application that allows a user (miner) 
	to connect to the DeNet decentralized network and mine tokens by granting access to their avaliable space for other users.`,
	Run: func(cmd *cobra.Command, args []string) {
		accs := account.List()
		if len(accs) == 0 {
			accountCreateCmd.Run(accountCreateCmd, []string{})
		} else {
			accountLoginCmd.Run(accountLoginCmd, []string{})
		}
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
