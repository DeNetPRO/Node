package cmd

import (
	"log"

	"github.com/DeNetPRO/src/account"
	blckChain "github.com/DeNetPRO/src/blockchain_provider"
	"github.com/DeNetPRO/src/cleaner"
	"github.com/DeNetPRO/src/paths"
	"github.com/DeNetPRO/src/rpcserver"
	tstpkg "github.com/DeNetPRO/src/tst_pkg"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

var debug = &cobra.Command{
	Use:   "debug",
	Short: "debug",
	Long:  "debug",
	Run: func(cmd *cobra.Command, args []string) {

		tstpkg.TestModeOn()

		paths.Init()

		addr, nodeConfig, err := account.Import()
		if err != nil {
			log.Fatal(err)
		}

		go blckChain.StartMakingProofs(common.HexToAddress(addr), tstpkg.Data().Password, nodeConfig)
		go cleaner.Start()

		rpcserver.Start(tstpkg.TestConfig().HTTPPort)
		if err != nil {
			log.Fatal(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(debug)
}
