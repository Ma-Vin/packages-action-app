package service

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/ma-vin/packages-action/config"
)

type DetermineCandidatesExecutor func(config *config.Config) (*[]Candidate, error)
type GitHubDeleteVersionRestExecutor func(packageName string, versionId int, config *config.Config) error
type GitHubDeletePackageRestExecutor func(packageName string, config *config.Config) error

var CandidatesExecutor DetermineCandidatesExecutor = initCandidatesExecutor()
var DeleteVersionExecutor GitHubDeleteVersionRestExecutor = initDeleteVersionExecutor()
var DeletePackageExecutor GitHubDeletePackageRestExecutor = initDeletePackageExecutor()

func initCandidatesExecutor() DetermineCandidatesExecutor {
	return func(config *config.Config) (*[]Candidate, error) {
		return DetermineCandidates(config)
	}
}

func initDeleteVersionExecutor() GitHubDeleteVersionRestExecutor {
	return func(packageName string, versionId int, config *config.Config) error {
		return DeleteUserPackageVersion(packageName, versionId, config)
	}
}

func initDeletePackageExecutor() GitHubDeletePackageRestExecutor {
	return func(packageName string, config *config.Config) error {
		return DeleteUserPackage(packageName, config)
	}
}

func InitAllDeletion() {
	CandidatesExecutor = initCandidatesExecutor()
	DeleteVersionExecutor = initDeleteVersionExecutor()
	DeletePackageExecutor = initDeletePackageExecutor()
}

// Deletes versions from Github with concurrency
func DeleteVersions(config *config.Config) error {
	candidates, err := CandidatesExecutor(config)
	if err != nil {
		return err
	}

	logCandidates(candidates)

	if config.DryRun {
		log.Println("Skip deletion because of dryRun")
		return nil
	}

	count := len(*candidates)

	channel := make(chan error, count)
	var wg sync.WaitGroup
	wg.Add(count)

	for _, c := range *candidates {
		go deleteCandidate(&c, config, channel, &wg)
	}

	wg.Wait()
	close(channel)

	withErrors := false
	for err := range channel {
		if err != nil {
			withErrors = true
			log.Println(err.Error())
		}
	}
	if withErrors {
		return errors.New("delete execution with errors")
	}
	return nil
}

// executes the deletion for a candidate and confirms it to wainting group
func deleteCandidate(candidate *Candidate, config *config.Config, channel chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	switch candidate.Type {
	case VERSION_CANDIDATE:
		channel <- DeleteVersionExecutor(config.PackageName, candidate.Id, config)
	case PACKAGE_CANDIDATE:
		channel <- DeletePackageExecutor(config.PackageName, config)
	default:
		channel <- fmt.Errorf("cannot delete candidate '%s' with id %d of unknown type", candidate.Name, candidate.Id)
	}
}

// logs the candidates which will be deleted
func logCandidates(candidates *[]Candidate) {
	log.Println("the following elements will be deleted")
	for i, c := range *candidates {
		log.Printf("  %d. type: %s name: '%s' id: %d created: %s updated: %s description: '%s'", i+1, getCandidateTypeText(&c.Type), c.Name, c.Id, c.CreatedAt, c.UpdatedAt, c.Description)
	}
}

// returns the text for a given candidate type
func getCandidateTypeText(typeValue *int) string {
	switch *typeValue {
	case VERSION_CANDIDATE:
		return "version"
	case PACKAGE_CANDIDATE:
		return "package"
	default:
		return "unknown"
	}
}
