package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFileName string
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run is command that allows you to start mining process of DFile tokens",
	Long:  `run is command that allows you to start mining process of DFile tokens`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(viper.AllKeys())
		fmt.Println(viper.GetString("address"))
		//
		//fmt.Println("Account address: ", )
		//dfileAccount := account.DFileAccount{}
		//var password string
		//
		//fmt.Println("Enter password: ")
		//bytePassword, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		//if err != nil {
		//	log.Fatal(err)
		//}
		//
		//password = string(bytePassword)
		//fmt.Println(viper.GetString("address"))
		//err = dfileAccount.LoadAccount(accountAddress, password)
		//if err != nil {
		//	switch err {
		//	case keystore.ErrDecrypt:
		//		fmt.Println(err.Error())
		//	default:
		//		log.Fatal(err)
		//	}
		//}

	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&cfgFileName, "config", "c", "", "config file name")
	err := runCmd.MarkFlagRequired("config")
	if err != nil {
		log.Fatal(err)
	}

}
