package cmd

import (
	"git.denetwork.xyz/DeNet/dfile-secondary-node/account"
	"github.com/spf13/cobra"
)

// RootCmd is an entry point for executing CLI commands.
var rootCmd = &cobra.Command{
	Use:   "DeNet-Node",
	Short: "DeNet-Node is a decentralized network node",
	Long:  `DeNet-Node is a CLI application that grants access to your machines unused space`,
	Run: func(cmd *cobra.Command, args []string) {
		accs := account.List()
		if len(accs) == 0 {
			accountImportCmd.Run(accountImportCmd, []string{})
		} else {
			accountLoginCmd.Run(accountLoginCmd, []string{})
		}
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
