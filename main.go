package main

import (
	"log"

	"git.denetwork.xyz/DeNet/dfile-secondary-node/cmd"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/logger"
	"git.denetwork.xyz/DeNet/dfile-secondary-node/paths"
)

func main() {
	err := paths.Init()
	if err != nil {
		logger.Log(logger.MarkLocation("main->", err))
		log.Fatal("Fatal Error: couldn't locate home directory")
	}

	// upnp.Init()

	cmd.Execute()
}
