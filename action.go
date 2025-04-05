package main

import (
	"github.com/ma-vin/packages-action/config"
	"github.com/ma-vin/packages-action/service"
	"github.com/ma-vin/typewriter/logger"
)

var (
	version    string
	gitSha     string
	branchName string
)

// Main funtion to execute the actions process
func main() {
	printVersion()
	logger.Information("Start packages action")
	initAll()

	var loadedConfig, err = config.ReadConfiguration()

	checkError(err)

	err = service.DeleteVersions(loadedConfig)
	checkError(err)

	logger.Information("Packages action done")
}

func checkError(err error) {
	if err != nil {
		logger.Fatalf("Packages action failed: %s", err)
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
		logger.Informationf("Version: %s", version)
	}
	if gitSha != "" {
		logger.Informationf("GitSha:  %s", gitSha)
	}
	if branchName != "" {
		logger.Informationf("Branch:  %s", branchName)
	}
}
