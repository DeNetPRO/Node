package cmd

import (
	"dfile-secondary-node/common"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
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
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&cfgFileName, "config", "c", "", "config file name")
	err := runCmd.MarkFlagRequired("config")
	if err != nil {
		log.Fatal(err)
	}

}

func initConfig() {
	fmt.Println(cfgFileName)
	configsDir, err := common.GetConfigsDirectory()
	if err != nil {
		log.Fatal(err)
	}
	viper.AddConfigPath(configsDir)
	viper.SetConfigName(cfgFileName)
	viper.SetConfigType("json")
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using configuration file: ", cfgFileName)
	}
}
