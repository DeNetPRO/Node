package cmd

import (
	"dfile-secondary-node/account"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
)

// runCmd represents the run command
var accountAddress string

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run is command that allows you to start mining process of DFile tokens",
	Long: `run is command that allows you to start mining process of DFile tokens`,
	Run: func(cmd *cobra.Command, args []string) {
		dfileAccount := account.DFileAccount{}
		var password string

		fmt.Println("Enter password: ")
		bytePassword, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			log.Fatal(err)
		}

		password = string(bytePassword)
		err = dfileAccount.LoadAccount(accountAddress, password)
		if err != nil {
			switch err {
			case keystore.ErrDecrypt:
				fmt.Println(err.Error())
			default:
				log.Fatal(err)
			}
		}


	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&accountAddress, "account_address", "a", "", "account address")
	err := runCmd.MarkFlagRequired("account_address")
	if err != nil {
		log.Fatal(err)
	}
}
