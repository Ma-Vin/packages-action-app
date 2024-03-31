package main

import (
	"log"

	"github.com/ma-vin/packages-action/config"
	"github.com/ma-vin/packages-action/service"
)

var (
	version    string
	gitSha     string
	branchName string
)

// Main funtion to execute the actions process
func main() {
	printVersion()
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

// prints the version, git hash and branch name if set by ldflags
func printVersion() {
	if version != "" {
		log.Printf("Version: %s", version)
	}
	if gitSha != "" {
		log.Printf("GitSha:  %s", gitSha)
	}
	if branchName != "" {
		log.Printf("Branch:  %s", branchName)
	}
}
