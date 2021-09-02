package test

import (
	"log"
	"os"
	"testing"

	"git.denetwork.xyz/dfile/dfile-secondary-node/paths"
	"git.denetwork.xyz/dfile/dfile-secondary-node/shared"
)

func TestMain(m *testing.M) {
	shared.TestMode = true

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
