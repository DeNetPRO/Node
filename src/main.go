package main

import (
	"log"

	"github.com/DeNetPRO/src/cmd"
	"github.com/DeNetPRO/src/logger"
	"github.com/DeNetPRO/src/paths"
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
