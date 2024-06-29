package service

import (
	"errors"
	"testing"

	"github.com/ma-vin/packages-action/config"
	"github.com/ma-vin/packages-action/testutil"
)

var deletionConf config.Config
var deletionVersionCandidate Candidate
var deletionPackageCandidate Candidate
var deletionCandidates *[]Candidate
var deletionCandidatesError error
var deleteVersionError error
var deletePackageError error
var countGetCandidatesExecuted int
var countDeleteVersionExecuted int
var countDeletePackageExecuted int

func initDeletionTest() {
	deletionConf = config.Config{}
	deletionVersionCandidate = Candidate{Id: 2, Name: "1.0.0", Description: "First Version", CreatedAt: "2024-03-12T20:00:00Z", UpdatedAt: "2024-03-13T16:00:00Z", Type: VERSION_CANDIDATE}
	deletionPackageCandidate = Candidate{Id: 1, Name: "DummyPackage", Description: "DummyPackage", CreatedAt: "2024-03-12T20:00:00Z", UpdatedAt: "2024-03-17T20:00:00Z", Type: PACKAGE_CANDIDATE}

	initDeletionExecutor()
}

func initDeletionExecutor() {
	countGetCandidatesExecuted = 0
	countDeleteVersionExecuted = 0
	countDeletePackageExecuted = 0

	deletionCandidates = nil

	deletionCandidatesError = nil
	deleteVersionError = nil
	deletePackageError = nil

	CandidatesExecutor = func(config *config.Config) (*[]Candidate, error) {
		countGetCandidatesExecuted++
		return deletionCandidates, deletionCandidatesError
	}
	DeleteVersionExecutor = func(packageName string, versionId int, config *config.Config) error {
		countDeleteVersionExecuted++
		return deleteVersionError
	}
	DeletePackageExecutor = func(packageName string, config *config.Config) error {
		countDeletePackageExecuted++
		return deletePackageError
	}
}

func TestDeleteVersionsSuccessfulVersionType(t *testing.T) {
	initDeletionTest()

	deletionCandidates = &[]Candidate{deletionVersionCandidate}

	err := DeleteVersions(&deletionConf)
	testutil.AssertNil(err, t, "err")
	testutil.AssertEquals(1, countGetCandidatesExecuted, t, "get candidates executed")
	testutil.AssertEquals(1, countDeleteVersionExecuted, t, "delete version executed")
	testutil.AssertEquals(0, countDeletePackageExecuted, t, "delete package executed")
}

func TestDeleteVersionsSuccessfulPackageType(t *testing.T) {
	initDeletionTest()

	deletionCandidates = &[]Candidate{deletionPackageCandidate}

	err := DeleteVersions(&deletionConf)
	testutil.AssertNil(err, t, "err")
	testutil.AssertEquals(1, countGetCandidatesExecuted, t, "get candidates executed")
	testutil.AssertEquals(0, countDeleteVersionExecuted, t, "delete version executed")
	testutil.AssertEquals(1, countDeletePackageExecuted, t, "delete package executed")
}

func TestDeleteVersionsNoCandidates(t *testing.T) {
	initDeletionTest()

	deletionCandidates = &[]Candidate{}

	err := DeleteVersions(&deletionConf)
	testutil.AssertNil(err, t, "err")
	testutil.AssertEquals(1, countGetCandidatesExecuted, t, "get candidates executed")
	testutil.AssertEquals(0, countDeleteVersionExecuted, t, "delete version executed")
	testutil.AssertEquals(0, countDeletePackageExecuted, t, "delete package executed")
}

func TestDeleteVersionsUnknwonType(t *testing.T) {
	initDeletionTest()

	deletionCandidates = &[]Candidate{{Id: 3, Name: "Unknwon", Description: "Unknwon", CreatedAt: "2024-03-12T20:00:00Z", UpdatedAt: "2024-03-17T20:00:00Z", Type: PACKAGE_CANDIDATE + VERSION_CANDIDATE + 1}}

	err := DeleteVersions(&deletionConf)
	testutil.AssertNotNil(err, t, "err")
	testutil.AssertEquals("delete execution with errors", err.Error(), t, "error message")
	testutil.AssertEquals(1, countGetCandidatesExecuted, t, "get candidates executed")
	testutil.AssertEquals(0, countDeleteVersionExecuted, t, "delete version executed")
	testutil.AssertEquals(0, countDeletePackageExecuted, t, "delete package executed")
}

func TestDeleteVersionsFailedGetCandidates(t *testing.T) {
	initDeletionTest()

	deletionCandidatesError = errors.New("testError")

	err := DeleteVersions(&deletionConf)
	testutil.AssertNotNil(err, t, "err")
	testutil.AssertEquals("testError", err.Error(), t, "error message")
	testutil.AssertEquals(1, countGetCandidatesExecuted, t, "get candidates executed")
	testutil.AssertEquals(0, countDeleteVersionExecuted, t, "delete version executed")
	testutil.AssertEquals(0, countDeletePackageExecuted, t, "delete package executed")
}

func TestDeleteVersionsFailedDeleteVersion(t *testing.T) {
	initDeletionTest()

	deletionCandidates = &[]Candidate{deletionVersionCandidate}
	deleteVersionError = errors.New("testError")

	err := DeleteVersions(&deletionConf)
	testutil.AssertNotNil(err, t, "err")
	testutil.AssertEquals("delete execution with errors", err.Error(), t, "error message")
	testutil.AssertEquals(1, countGetCandidatesExecuted, t, "get candidates executed")
	testutil.AssertEquals(1, countDeleteVersionExecuted, t, "delete version executed")
	testutil.AssertEquals(0, countDeletePackageExecuted, t, "delete package executed")
}

func TestDeleteVersionsFailedDeletePackage(t *testing.T) {
	initDeletionTest()

	deletionCandidates = &[]Candidate{deletionPackageCandidate}
	deletePackageError = errors.New("testError")

	err := DeleteVersions(&deletionConf)
	testutil.AssertNotNil(err, t, "err")
	testutil.AssertEquals("delete execution with errors", err.Error(), t, "error message")
	testutil.AssertEquals(1, countGetCandidatesExecuted, t, "get candidates executed")
	testutil.AssertEquals(0, countDeleteVersionExecuted, t, "delete version executed")
	testutil.AssertEquals(1, countDeletePackageExecuted, t, "delete package executed")
}

func TestDeleteVersionsSuccessfulDryRun(t *testing.T) {
	initDeletionTest()

	deletionConf.DryRun = true
	deletionCandidates = &[]Candidate{deletionVersionCandidate, deletionPackageCandidate}

	err := DeleteVersions(&deletionConf)
	testutil.AssertNil(err, t, "err")
	testutil.AssertEquals(1, countGetCandidatesExecuted, t, "get candidates executed")
	testutil.AssertEquals(0, countDeleteVersionExecuted, t, "delete version executed")
	testutil.AssertEquals(0, countDeletePackageExecuted, t, "delete package executed")
}

func TestDeleteVersionsSuccessfulMultiple(t *testing.T) {
	initDeletionTest()

	deletionVersionCandidateTwo := Candidate{Id: 4, Name: "2.0.0", Description: "Second Version", CreatedAt: "2024-03-17T20:00:00Z", UpdatedAt: "2024-03-17T20:00:00Z", Type: VERSION_CANDIDATE}
	deletionPackageCandidateTwo := Candidate{Id: 3, Name: "OtherDummyPackage", Description: "DummyPackage", CreatedAt: "2024-03-17T20:00:00Z", UpdatedAt: "2024-03-17T20:00:00Z", Type: PACKAGE_CANDIDATE}

	deletionCandidates = &[]Candidate{deletionVersionCandidate, deletionVersionCandidateTwo, deletionPackageCandidate, deletionPackageCandidateTwo}

	err := DeleteVersions(&deletionConf)
	testutil.AssertNil(err, t, "err")
	testutil.AssertEquals(1, countGetCandidatesExecuted, t, "get candidates executed")
	testutil.AssertEquals(2, countDeleteVersionExecuted, t, "delete version executed")
	testutil.AssertEquals(2, countDeletePackageExecuted, t, "delete package executed")
}
