package main

import (
	"log"

	"github.com/ma-vin/packages-action/config"
	"github.com/ma-vin/packages-action/service"
)

// Main funtion to execute the actions process
func main() {
	log.Println("Start packages action")
	initAll()

	var loadedConfig, err = config.ReadConfiguration()

	checkError(err)

	err = service.DeleteVersions(loadedConfig)
	checkError(err)

	log.Println("Packages action done")
}

func checkError(err error) {
	if err != nil {
		log.Fatalf("Packages action failed: %s", err)
	}
}

func initAll() {
	service.InitAllCandidates()
	service.InitAllDeletion()
	service.InitAllGitHubRest()
}
