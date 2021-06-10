package test

import (
	"dfile-secondary-node/shared"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	err := shared.InitPaths()
	if err != nil {
		log.Fatal("Fatal Error: couldn't locate home directory")
	}
	exitVal := m.Run()

	err = os.RemoveAll(shared.WorkDirPath)
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(exitVal)
}
