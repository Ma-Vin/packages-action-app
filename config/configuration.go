package config

import (
	"log"
	"os"
	"strings"
)

const (
	// package type for maven
	MAVEN string = "maven"
	// package type for a not supported or unknown type
	UNKNOWN string = "unkown"
)

// structure to hold configuration of the action
type Config struct {
	// organization whose packages are to handle (Either this or user has to be set)
	Organization string
	// user whose packages are to handle  (Either this or organization has to be set)
	User string
	// package type whiche are to handle (not nil)
	PackageType string
	// name of the package (not nil)
	PackageName string
	// indicator whether to delete snapshots or not
	DeleteSnapshots bool
	// token which is to use to authenticate against github rest api (not nil)
	GithubToken string
}

/*
Reads the configuration from environment variables:
  - ORGANIZATION
  - USER
  - PACKAGE_TYPE
  - PACKAGE_NAME
  - DELETE_SNAPSHOTS
  - GITHUB_TOKEN
*/
func ReadConfiguration() *Config {
	var config Config
	config.Organization = getTrimEnv("ORGANIZATION")
	config.User = getTrimEnv("USER")
	config.PackageType = mapToPackageType(getTrimEnv("PACKAGE_TYPE"))
	config.PackageName = getTrimEnv("PACKAGE_NAME")
	config.DeleteSnapshots = strings.EqualFold("true", getTrimEnv("DELETE_SNAPSHOTS"))
	config.GithubToken = getTrimEnv("GITHUB_TOKEN")

	printConfig(&config)

	if isValid(&config) {
		return &config
	}
	log.Println("Invalid configuration!")
	return nil
}

// determines an environment variable and return it as trimmed string
func getTrimEnv(envName string) string {
	return strings.TrimSpace(os.Getenv(envName))
}

// maps a given string to a package type
func mapToPackageType(toMap string) string {
	switch strings.ToLower(toMap) {
	case MAVEN:
		return MAVEN
	default:
		return UNKNOWN
	}
}

// Checks whether a given configuration is valid or not
func isValid(config *Config) bool {
	if (config.Organization != "" && config.User != "") || (config.Organization == "" && config.User == "") {
		log.Println("Fill exactly one: Either user or organization")
		return false
	}
	if config.PackageType == UNKNOWN {
		log.Println("The packagetype is unknown")
		return false
	}
	if config.PackageName == "" {
		log.Println("Missing package name")
		return false
	}
	if config.GithubToken == "" {
		log.Println("Missing GitHub token")
		return false
	}
	return true
}

// prints a given configuration to the standard output
func printConfig(config *Config) {
	log.Println("Read configuration", config.Organization)
	log.Println("  Organization:    ", config.Organization)
	log.Println("  User:            ", config.User)
	log.Println("  PackageType:     ", config.PackageType)
	log.Println("  PackageName:     ", config.PackageName)
	log.Println("  DeleteSnapshots: ", config.DeleteSnapshots)
	if config.GithubToken != "" {
		log.Println("  GithubToken:      ***")
	} else {
		log.Println("  GithubToken:")
	}
}
