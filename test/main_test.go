package test

import (
	"log"
	"os"
	"testing"

	"git.denetwork.xyz/dfile/dfile-secondary-node/paths"
)

func TestMain(m *testing.M) {
	paths.WorkDirName = "dfile-test"

	err := paths.Init()
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
