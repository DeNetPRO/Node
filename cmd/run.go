package cmd

import (
	"dfile-secondary-node/account"
	"github.com/spf13/cobra"
	"log"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run is command that allows you to start mining process of DFile tokens",
	Long: `run is command that allows you to start mining process of DFile tokens`,
	Run: func(cmd *cobra.Command, args []string) {
		dfileAccount := account.DFileAccount{}
		err := dfileAccount.LoadAccount("0x38EA1c699993327f15b7AF5CD51E6e3DF5da0129", "password")
		if err != nil {
			log.Fatal(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
