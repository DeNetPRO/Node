package cmd

import (
	"dfile-secondary-node/account"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dfile-secondary-node",
	Short: "dfile-secondary-node is a CLI application that allows a user (miner) to connect to the DeNet decentralized network and mine DFile tokens.",
	Long: `dfile-secondary-node is a CLI application that allows a user (miner) 
	to connect to the DeNet decentralized network and mine DFile tokens.`,
	Run: func(cmd *cobra.Command, args []string) {
		accs := account.GetAllAccounts()

		if len(accs) == 0 {
			accountCreateCmd.Run(accountCreateCmd, []string{})
		} else {
			accountListCmd.Run(accountListCmd, []string{})
			accountCheckCmd.Run(accountCheckCmd, []string{})
		}

	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
}
