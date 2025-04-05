package service

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ma-vin/packages-action/config"
	"github.com/ma-vin/packages-action/service/github_model"
	"github.com/ma-vin/typewriter/logger"
)

const (
	VERSION_CANDIDATE int = iota
	PACKAGE_CANDIDATE int = iota
)

type Candidate struct {
	Name        string
	Id          int
	Description string
	CreatedAt   string
	UpdatedAt   string
	Type        int
}

type GitHubGetVersionsRestExecutor func(config *config.Config) (*[]github_model.Version, error)
type GitHubGetPackageRestExecutor func(config *config.Config) (*github_model.UserPackage, error)
type GitHubGetAllPackagesRestExecutor func(config *config.Config) (*[]github_model.UserPackage, error)

var VersionsGetExecutor GitHubGetVersionsRestExecutor = initVersionsGetExecutor()
var PackageGetExecutor GitHubGetPackageRestExecutor = initPackageGetExecutor()
var AllPackagesGetExecutor GitHubGetAllPackagesRestExecutor = initAllPackagesGetExecutor()

func initVersionsGetExecutor() GitHubGetVersionsRestExecutor {
	return func(config *config.Config) (*[]github_model.Version, error) {
		return GetUserPackageVersions(config.PackageName, config)
	}
}

func initPackageGetExecutor() GitHubGetPackageRestExecutor {
	return func(config *config.Config) (*github_model.UserPackage, error) {
		return GetUserPackage(config.PackageName, config)
	}
}

func initAllPackagesGetExecutor() GitHubGetAllPackagesRestExecutor {
	return func(config *config.Config) (*[]github_model.UserPackage, error) {
		return GetUserPackages(config)
	}
}

func InitAllCandidates() {
	VersionsGetExecutor = initVersionsGetExecutor()
	PackageGetExecutor = initPackageGetExecutor()
	AllPackagesGetExecutor = initAllPackagesGetExecutor()
}

// Determine all candidates to delete. A candidate can be either a version or a package
// If a package would be empty after version deletion, the package is to be deleted
func DetermineCandidates(config *config.Config) (*[]Candidate, error) {
	existence, err := checkPackageExistence(config)
	if err != nil {
		return nil, err
	}
	if !existence {
		logger.Warningf("There does not exists a package with name %s of type %s at user %s: skip deletion", config.PackageName, config.PackageType, config.User)
		return &[]Candidate{}, nil
	}

	candidates, deletePackage, err := determineRelevantVersions(config)
	if err != nil {
		return nil, err
	}

	if !deletePackage {
		return candidates, nil
	}

	candidate, err := determineRelevantPackage(config)
	if err != nil {
		return nil, err
	}
	return &[]Candidate{*candidate}, nil
}

// Checks whether there exists the package for the user
func checkPackageExistence(config *config.Config) (bool, error) {
	packages, err := AllPackagesGetExecutor(config)
	if err != nil {
		return false, err
	}
	if packages == nil || len(*packages) == 0 {
		return false, nil
	}
	for _, p := range *packages {
		if p.Name == config.PackageName {
			return true, nil
		}
	}
	return false, nil
}

// Determines all relevant versions which can be deleted and an indicator if package would be empty after version deletion
func determineRelevantVersions(config *config.Config) (*[]Candidate, bool, error) {
	versions, err := VersionsGetExecutor(config)
	if err != nil {
		return nil, false, err
	}

	versionNameParts, isSnapshot, err := splitVersionNames(versions)
	if err != nil {
		return nil, false, err
	}

	var res []Candidate
	for i, v := range *versions {
		if isVersionRelevant(&i, versions, versionNameParts, isSnapshot, config) {
			res = append(res, Candidate{v.Name, v.Id, v.Description, v.CreatedAt, v.UpdatedAt, VERSION_CANDIDATE})
		}
	}

	return &res, len(res) > 0 && len(*versions) == len(res), nil
}

// Determine the relevant package which is to be deleted
func determineRelevantPackage(config *config.Config) (*Candidate, error) {
	pack, err := PackageGetExecutor(config)
	if err != nil {
		return nil, err
	}

	return &Candidate{pack.Name, pack.Id, pack.Name, pack.CreatedAt, pack.UpdatedAt, PACKAGE_CANDIDATE}, nil
}

// Checks if a version at a given index is to be deleted or not
func isVersionRelevant(index *int, versions *[]github_model.Version, versionNameParts *[][]int, isSnapshot *[]bool, config *config.Config) bool {
	isIndexSnapshot := (*isSnapshot)[*index]

	versionNameMatch := config.VersionNameToDelete != "" && strings.EqualFold((*versions)[*index].Name, config.VersionNameToDelete)
	snapshotDelete := config.DeleteSnapshots && isIndexSnapshot
	deleteMajor := !isIndexSnapshot && config.NumberOfMajorVersionsToKeep > 0 && countCreaterMajorVersions(index, versionNameParts, isSnapshot) >= config.NumberOfMajorVersionsToKeep
	deleteMinor := !isIndexSnapshot && config.NumberOfMinorVersionsToKeep > 0 && countCreaterMinorVersions(index, versionNameParts, isSnapshot) >= config.NumberOfMinorVersionsToKeep
	deletePatch := !isIndexSnapshot && config.NumberOfPatchVersionsToKeep > 0 && countCreaterPatchVersions(index, versionNameParts, isSnapshot) >= config.NumberOfPatchVersionsToKeep

	return versionNameMatch || snapshotDelete || deleteMajor || deleteMinor || deletePatch
}

// Split the name of given versions into major, minor and patch tripel. In addition an indicator whether a version is a snapshot or not
func splitVersionNames(versions *[]github_model.Version) (*[][]int, *[]bool, error) {
	resSplit := make([][]int, len(*versions))
	resSnapshot := make([]bool, len(*versions))

	for i, v := range *versions {

		nameToSplit, isSnaphot := strings.CutSuffix(strings.ToLower(v.Name), "-snapshot")
		resSnapshot[i] = isSnaphot

		parts := strings.Split(nameToSplit, ".")
		if len(parts) > 3 {
			return nil, nil, fmt.Errorf("there are more items than 'major.minor.patch' or 'major.minor.patch-SNAPSHOT' at version name '%s' with id %d", v.Name, v.Id)
		}
		err := addVersionPart(&parts, 0, "major", &i, &resSplit, &v)
		if err != nil {
			return nil, nil, err
		}
		err = addVersionPart(&parts, 1, "minor", &i, &resSplit, &v)
		if err != nil {
			return nil, nil, err
		}
		err = addVersionPart(&parts, 2, "patch", &i, &resSplit, &v)
		if err != nil {
			return nil, nil, err
		}
	}

	return &resSplit, &resSnapshot, nil
}

// adds the major, minor or patch number to splittedVersions for a given version index. there is none, zero will be set
func addVersionPart(parts *[]string, partIndex int, partIndexName string, versionIndex *int, splittedVersions *[][]int, version *github_model.Version) error {
	if len(*parts) > partIndex {
		major, err := strconv.Atoi((*parts)[partIndex])
		if err != nil {
			return fmt.Errorf("failed to format %s version to int at version name '%s' with id %d: %v", partIndexName, version.Name, version.Id, err)
		}
		(*splittedVersions)[*versionIndex] = append((*splittedVersions)[*versionIndex], major)
	} else {
		(*splittedVersions)[*versionIndex] = append((*splittedVersions)[*versionIndex], 0)
	}
	return nil
}

// Counts the versions which have a greater major version than the one at given index
func countCreaterMajorVersions(index *int, versionNameParts *[][]int, isSnapshot *[]bool) int {
	counter := 0
	for i, parts := range *versionNameParts {
		if !(*isSnapshot)[i] && parts[0] > (*versionNameParts)[*index][0] {
			counter++
		}
	}
	return counter
}

// Counts the versions which have equal major but greater minor version than the one at given index
func countCreaterMinorVersions(index *int, versionNameParts *[][]int, isSnapshot *[]bool) int {
	counter := 0
	for i, parts := range *versionNameParts {
		if !(*isSnapshot)[i] && parts[0] == (*versionNameParts)[*index][0] && parts[1] > (*versionNameParts)[*index][1] {
			counter++
		}
	}
	return counter
}

// Counts the versions which have equal major and minor but greater patch version than the one at given index
func countCreaterPatchVersions(index *int, versionNameParts *[][]int, isSnapshot *[]bool) int {
	counter := 0
	for i, parts := range *versionNameParts {
		if !(*isSnapshot)[i] && parts[0] == (*versionNameParts)[*index][0] && parts[1] == (*versionNameParts)[*index][1] && parts[2] > (*versionNameParts)[*index][2] {
			counter++
		}
	}
	return counter
}
