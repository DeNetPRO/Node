package test

import (
	"dfile-secondary-node/paths"
	"dfile-secondary-node/shared"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	paths.WorkDirName = "dfile-test"

	err := shared.InitPaths()
	if err != nil {
		log.Fatal("Fatal Error: couldn't locate home directory")
	}
	exitVal := m.Run()

	err = os.RemoveAll(paths.WorkDirPath)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(exitVal)
}
