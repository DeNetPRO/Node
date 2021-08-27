/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"log"
	"os"

	"git.denetwork.xyz/dfile/dfile-secondary-node/cmd"
	"git.denetwork.xyz/dfile/dfile-secondary-node/logger"
	"git.denetwork.xyz/dfile/dfile-secondary-node/paths"
	"git.denetwork.xyz/dfile/dfile-secondary-node/shared"
	"git.denetwork.xyz/dfile/dfile-secondary-node/upnp"
)

var testMode = "test"

func main() {
	mode := os.Getenv("MODE")
	if mode == testMode {
		shared.TestMode = true
	}

	err := paths.Init()
	if err != nil {
		logger.Log(logger.CreateDetails("main->", err))
		log.Fatal("Fatal Error: couldn't locate home directory")
	}

	upnp.Init()

	cmd.Execute()
}
