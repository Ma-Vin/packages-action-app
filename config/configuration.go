package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	// package type for maven
	MAVEN string = "maven"
	// package type for a not supported or unknown type
	UNKNOWN string = "unkown"

	ENV_NAME_ORGANIZATION           string = "ORGANIZATION"
	ENV_NAME_USER                   string = "USER"
	ENV_NAME_PACKAGE_TYPE           string = "PACKAGE_TYPE"
	ENV_NAME_PACKAGE_NAME           string = "PACKAGE_NAME"
	ENV_NAME_VERSION_NAME_TO_DELETE string = "VERSION_NAME_TO_DELETE"
	ENV_NAME_DELETE_SNAPSHOTS       string = "DELETE_SNAPSHOTS"
	ENV_NAME_NUMBER_MAJOR_TO_KEEP   string = "NUMBER_MAJOR_TO_KEEP"
	ENV_NAME_NUMBER_MINOR_TO_KEEP   string = "NUMBER_MINOR_TO_KEEP"
	ENV_NAME_NUMBER_PATCH_TO_KEEP   string = "NUMBER_PATCH_TO_KEEP"
	ENV_NAME_GITHUB_TOKEN           string = "GITHUB_TOKEN"
	ENV_NAME_DRY_RUN                string = "DRY_RUN"
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
	// name of a version to delete
	VersionNameToDelete string
	// indicator whether to delete snapshots or not
	DeleteSnapshots bool
	// Number major versions to keep
	NumberOfMajorVersionsToKeep int
	// Number minor versions to keep
	NumberOfMinorVersionsToKeep int
	// Number patch versions to keep
	NumberOfPatchVersionsToKeep int
	// token which is to use to authenticate against github rest api (not nil)
	GithubToken string
	// Indicator wether to run application without deletion or not. Default true (No deletetion)
	DryRun bool
}

/*
Reads the configuration from environment variables:
  - ORGANIZATION
  - USER
  - PACKAGE_TYPE
  - PACKAGE_NAME
  - VERSION_NAME_TO_DELETE
  - DELETE_SNAPSHOTS
  - NUMBER_MAJOR_TO_KEEP
  - NUMBER_MINOR_TO_KEEP
  - NUMBER_PATCH_TO_KEEP
  - GITHUB_TOKEN
*/
func ReadConfiguration() *Config {
	var config Config
	config.Organization = getTrimEnv(ENV_NAME_ORGANIZATION)
	config.User = getTrimEnv(ENV_NAME_USER)
	config.PackageType = mapToPackageType(getTrimEnv(ENV_NAME_PACKAGE_TYPE))
	config.PackageName = getTrimEnv(ENV_NAME_PACKAGE_NAME)
	config.VersionNameToDelete = getTrimEnv(ENV_NAME_VERSION_NAME_TO_DELETE)
	config.DeleteSnapshots = getBoolEnv(ENV_NAME_DELETE_SNAPSHOTS)
	config.NumberOfMajorVersionsToKeep = getIntEnv(ENV_NAME_NUMBER_MAJOR_TO_KEEP)
	config.NumberOfMinorVersionsToKeep = getIntEnv(ENV_NAME_NUMBER_MINOR_TO_KEEP)
	config.NumberOfPatchVersionsToKeep = getIntEnv(ENV_NAME_NUMBER_PATCH_TO_KEEP)
	config.GithubToken = getTrimEnv(ENV_NAME_GITHUB_TOKEN)
	config.DryRun = getBoolEnvDefault(ENV_NAME_DRY_RUN, true)

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

// determines an environment variable and return it as bool
func getBoolEnv(envName string) bool {
	return strings.EqualFold("true", getTrimEnv(envName))
}

// determines an environment variable and return it as bool if present, other wise the given default value
func getBoolEnvDefault(envName string, defaultValue bool) bool {
	value := getTrimEnv(envName)
	if value == "" {
		return defaultValue
	}
	return strings.EqualFold("true", value)
}

// determines an environment variable and return it as int.
func getIntEnv(envName string) int {
	envValue := getTrimEnv(envName)
	if envValue != "" {
		res, err := strconv.Atoi(os.Getenv(envName))
		if err != nil {
			log.Println(err)
			return -1
		}
		if res <= 0 {
			log.Println("Only positiv values are allowed but found ", res, " for environment variable ", envName)
			return -1
		}
		return res
	}
	return -1
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
	if config.VersionNameToDelete == "" && !config.DeleteSnapshots &&
		config.NumberOfMajorVersionsToKeep <= 0 && config.NumberOfMinorVersionsToKeep <= 0 && config.NumberOfPatchVersionsToKeep <= 0 {
		log.Println("Nothing configured to delete: set a conrete version name, snapshot deletion or major, minor or patch to keep")
		return false
	}
	return true
}

// prints a given configuration to the standard output
func printConfig(config *Config) {
	log.Println("Read configuration", config.Organization)
	log.Println("  Organization:        ", config.Organization)
	log.Println("  User:                ", config.User)
	log.Println("  PackageType:         ", config.PackageType)
	log.Println("  PackageName:         ", config.PackageName)
	log.Println("  VersionNameToDelete: ", config.VersionNameToDelete)
	log.Println("  DeleteSnapshots:     ", config.DeleteSnapshots)
	printPositiv("  MajorVersionsToKeep: ", config.NumberOfMajorVersionsToKeep)
	printPositiv("  MinorVersionsToKeep: ", config.NumberOfMinorVersionsToKeep)
	printPositiv("  PatchVersionsToKeep: ", config.NumberOfPatchVersionsToKeep)
	if config.GithubToken != "" {
		log.Println("  GithubToken:          ***")
	} else {
		log.Println("  GithubToken:")
	}
	log.Println("  DryRun:              ", config.DryRun)
}

func printPositiv(text string, value int) {
	if 0 < value {
		log.Println(text, value)
		return
	}
	log.Println(text)
}
