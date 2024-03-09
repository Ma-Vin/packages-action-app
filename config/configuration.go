package config

import (
	"log"
	"os"
	"strings"
)

// type of possible supported packages
type PackageType int

const (
	// package type for maven
	MAVEN PackageType = iota
	// package type for a not supported or unknown type
	UNKNOWN PackageType = iota
)

// structure to hold configuration of the action
type Config struct {
	// organization whose packages are to handle (Either this or user has to be set)
	Organization string
	// user whose packages are to handle  (Either this or organization has to be set)
	User string
	// package type whiche are to handle (not nil)
	PackageType PackageType
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
func mapToPackageType(toMap string) PackageType {
	switch strings.ToLower(toMap) {
	case "maven":
		return MAVEN
	default:
		return UNKNOWN
	}
}

// maps a given to package type to a readable string
func MapFromPackageType(toMap PackageType) string {
	switch toMap {
	case MAVEN:
		return "maven"
	case UNKNOWN:
		return "unkown"
	default:
		return ""
	}
}

// Checks whether a given configuration is valid or not
func isValid(config *Config) bool {
	return (config.Organization != "" || config.User != "") && config.PackageType != UNKNOWN && config.PackageName != "" && config.GithubToken != ""
}

// prints a given configuration to the standard output
func printConfig(config *Config) {
	log.Println("Read configuration", config.Organization)
	log.Println("  Organization:    ", config.Organization)
	log.Println("  User:            ", config.User)
	log.Println("  PackageType:     ", MapFromPackageType(config.PackageType))
	log.Println("  PackageName:     ", config.PackageName)
	log.Println("  DeleteSnapshots: ", config.DeleteSnapshots)
	if config.GithubToken != "" {
		log.Println("  GithubToken:      ***")
	} else {
		log.Println("  GithubToken:")
	}
}
