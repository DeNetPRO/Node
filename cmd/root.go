package cmd

import (
	"git.denetwork.xyz/dfile/dfile-secondary-node/account"
	"github.com/spf13/cobra"
)

// RootCmd is an entry point for executing CLI commands.
var rootCmd = &cobra.Command{
	Use:   "dfile-node",
	Short: "dfile-node provides unused space for users that need it",
	Long: `dfile-node is a CLI application that allows a user (miner) 
	to connect to the DeNet decentralized network and mine DFile tokens by granting access to their avaliable space for other users.`,
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
