package cmd

import (
	"log"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/account"
	blckChain "git.denetwork.xyz/DeNet/dfile-secondary-node/blockchain_provider"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/cleaner"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/paths"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/rpcserver"
	tstpkg "git.denetwork.xyz/DeNet/dfile-secondary-node/tst_pkg"
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
