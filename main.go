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
	"dfile-secondary-node/cmd"
	"dfile-secondary-node/logger"
	"dfile-secondary-node/shared"
	"dfile-secondary-node/upnp"
	"log"
)

func main() {
	err := shared.InitPaths()
	if err != nil {
		logger.LogError("main->", err)
		log.Fatal("Fatal Error: couldn't locate home directory")
	}

	upnp.InitIGD()

	cmd.Execute()
}
