package service

import (
	"errors"
	"testing"

	"github.com/ma-vin/packages-action/config"
	"github.com/ma-vin/packages-action/service/github_model"
	"github.com/ma-vin/packages-action/testutil"
)

var candidatesConf config.Config
var candidateVersionOne github_model.Version
var candidateVersionTwo github_model.Version
var candidateVersionThreee github_model.Version
var candidatePacakge github_model.UserPackage

func initCandidateTest() {
	candidatesConf = config.Config{User: "DummyUser", PackageType: "maven", PackageName: "DummyPackage", VersionNameToDelete: "", DeleteSnapshots: false, NumberOfMajorVersionsToKeep: -1, NumberOfMinorVersionsToKeep: -1, NumberOfPatchVersionsToKeep: -1}

	candidateVersionOne = github_model.Version{Id: 2, Name: "1.0.0", Description: "First Version", CreatedAt: "2024-03-12T20:00:00Z", UpdatedAt: "2024-03-13T16:00:00Z"}
	candidateVersionTwo = github_model.Version{Id: 3, Name: "1.1.0", Description: "Second Version", CreatedAt: "2024-03-13T20:00:00Z", UpdatedAt: "2024-03-14T16:00:00Z"}
	candidateVersionThreee = github_model.Version{Id: 4, Name: "1.1.1", Description: "Third Version", CreatedAt: "2024-03-14T20:00:00Z", UpdatedAt: "2024-03-20:00:00Z"}

	candidatePacakge = github_model.UserPackage{Id: 1, Name: "DummyPackage", CreatedAt: "2024-03-12T20:00:00Z", UpdatedAt: "2024-03-20:00:00Z"}

	VersionsGetExecutor = func(config *config.Config) (*[]github_model.Version, error) {
		return &[]github_model.Version{candidateVersionOne, candidateVersionTwo, candidateVersionThreee}, nil
	}

	PackageGetExecutor = func(config *config.Config) (*github_model.UserPackage, error) {
		return &candidatePacakge, nil
	}
}

func TestDetermineCandidatesVersionName(t *testing.T) {
	initCandidateTest()

	candidatesConf.VersionNameToDelete = "1.1.0"

	candidates, err := DetermineCandidates(&candidatesConf)

	testutil.AssertNil(err, t, "err")
	testutil.AssertNotNil(candidates, t, "candidates")
	testutil.AssertEquals(1, len(*candidates), t, "len candidates")
	testutil.AssertEquals(3, (*candidates)[0].Id, t, "id")
	testutil.AssertEquals(VERSION_CANDIDATE, (*candidates)[0].Type, t, "type")
}

func TestDetermineCandidatesSnapshot(t *testing.T) {
	initCandidateTest()

	candidatesConf.DeleteSnapshots = true
	candidateVersionTwo.Name = "1.1.0-SNAPSHOT"

	candidates, err := DetermineCandidates(&candidatesConf)

	testutil.AssertNil(err, t, "err")
	testutil.AssertNotNil(candidates, t, "candidates")
	testutil.AssertEquals(1, len(*candidates), t, "len candidates")
	testutil.AssertEquals(3, (*candidates)[0].Id, t, "id")
	testutil.AssertEquals(VERSION_CANDIDATE, (*candidates)[0].Type, t, "type")
}

func TestDetermineCandidatesMajorButAllTheSame(t *testing.T) {
	initCandidateTest()

	candidatesConf.NumberOfMajorVersionsToKeep = 1

	candidates, err := DetermineCandidates(&candidatesConf)

	testutil.AssertNil(err, t, "err")
	testutil.AssertNotNil(candidates, t, "candidates")
	testutil.AssertEquals(0, len(*candidates), t, "len candidates")
}

func TestDetermineCandidatesMajorAllDifferent(t *testing.T) {
	initCandidateTest()

	candidatesConf.NumberOfMajorVersionsToKeep = 1
	candidateVersionTwo.Name = "2.0.0"
	candidateVersionThreee.Name = "3.0.0"

	candidates, err := DetermineCandidates(&candidatesConf)

	testutil.AssertNil(err, t, "err")
	testutil.AssertNotNil(candidates, t, "candidates")
	testutil.AssertEquals(2, len(*candidates), t, "len candidates")
	testutil.AssertNotEquals(4, (*candidates)[0].Id, t, "id 1. entry candidates")
	testutil.AssertNotEquals(4, (*candidates)[1].Id, t, "id 2. entry candidates")
	testutil.AssertEquals(VERSION_CANDIDATE, (*candidates)[0].Type, t, "type 1. entry candidates")
	testutil.AssertEquals(VERSION_CANDIDATE, (*candidates)[1].Type, t, "type 2. entry candidates")
}

func TestDetermineCandidatesMajorAllDifferentWithSnapshot(t *testing.T) {
	initCandidateTest()

	candidatesConf.NumberOfMajorVersionsToKeep = 1
	candidateVersionTwo.Name = "2.0.0-SNAPSHOT"
	candidateVersionThreee.Name = "3.0.0"

	candidates, err := DetermineCandidates(&candidatesConf)

	testutil.AssertNil(err, t, "err")
	testutil.AssertNotNil(candidates, t, "candidates")
	testutil.AssertEquals(1, len(*candidates), t, "len candidates")
	testutil.AssertEquals(2, (*candidates)[0].Id, t, "id 1. entry candidates")
	testutil.AssertEquals(VERSION_CANDIDATE, (*candidates)[0].Type, t, "type 1. entry candidates")
}

func TestDetermineCandidatesMinorButAllTheSame(t *testing.T) {
	initCandidateTest()

	candidatesConf.NumberOfMinorVersionsToKeep = 1
	candidateVersionTwo.Name = "1.0.1"
	candidateVersionThreee.Name = "1.0.2"

	candidates, err := DetermineCandidates(&candidatesConf)

	testutil.AssertNil(err, t, "err")
	testutil.AssertNotNil(candidates, t, "candidates")
	testutil.AssertEquals(0, len(*candidates), t, "len candidates")
}

func TestDetermineCandidatesMinorAllDifferent(t *testing.T) {
	initCandidateTest()

	candidatesConf.NumberOfMinorVersionsToKeep = 1
	candidateVersionTwo.Name = "1.1.0"
	candidateVersionThreee.Name = "1.2.0"

	candidates, err := DetermineCandidates(&candidatesConf)

	testutil.AssertNil(err, t, "err")
	testutil.AssertNotNil(candidates, t, "candidates")
	testutil.AssertEquals(2, len(*candidates), t, "len candidates")
	testutil.AssertNotEquals(4, (*candidates)[0].Id, t, "id 1. entry candidates")
	testutil.AssertNotEquals(4, (*candidates)[1].Id, t, "id 2. entry candidates")
	testutil.AssertEquals(VERSION_CANDIDATE, (*candidates)[0].Type, t, "type 1. entry candidates")
	testutil.AssertEquals(VERSION_CANDIDATE, (*candidates)[1].Type, t, "type 2. entry candidates")
}

func TestDetermineCandidatesMinorAllDifferentWithSnapshot(t *testing.T) {
	initCandidateTest()

	candidatesConf.NumberOfMinorVersionsToKeep = 1
	candidateVersionTwo.Name = "1.1.0-SNAPSHOT"
	candidateVersionThreee.Name = "1.2.0"

	candidates, err := DetermineCandidates(&candidatesConf)

	testutil.AssertNil(err, t, "err")
	testutil.AssertNotNil(candidates, t, "candidates")
	testutil.AssertEquals(1, len(*candidates), t, "len candidates")
	testutil.AssertEquals(2, (*candidates)[0].Id, t, "id 1. entry candidates")
	testutil.AssertEquals(VERSION_CANDIDATE, (*candidates)[0].Type, t, "type 1. entry candidates")
}

func TestDetermineCandidatesPatchButAllAtDifferentMinor(t *testing.T) {
	initCandidateTest()

	candidatesConf.NumberOfPatchVersionsToKeep = 1
	candidateVersionTwo.Name = "1.1.1"
	candidateVersionThreee.Name = "1.2.2"

	candidates, err := DetermineCandidates(&candidatesConf)

	testutil.AssertNil(err, t, "err")
	testutil.AssertNotNil(candidates, t, "candidates")
	testutil.AssertEquals(0, len(*candidates), t, "len candidates")
}

func TestDetermineCandidatesPatchButAllAtDifferentMajor(t *testing.T) {
	initCandidateTest()

	candidatesConf.NumberOfPatchVersionsToKeep = 1
	candidateVersionTwo.Name = "2.0.1"
	candidateVersionThreee.Name = "3.0.2"

	candidates, err := DetermineCandidates(&candidatesConf)

	testutil.AssertNil(err, t, "err")
	testutil.AssertNotNil(candidates, t, "candidates")
	testutil.AssertEquals(0, len(*candidates), t, "len candidates")
}

func TestDetermineCandidatesPatchAllDifferent(t *testing.T) {
	initCandidateTest()

	candidatesConf.NumberOfPatchVersionsToKeep = 1
	candidateVersionTwo.Name = "1.0.1"
	candidateVersionThreee.Name = "1.0.2"

	candidates, err := DetermineCandidates(&candidatesConf)

	testutil.AssertNil(err, t, "err")
	testutil.AssertNotNil(candidates, t, "candidates")
	testutil.AssertEquals(2, len(*candidates), t, "len candidates")
	testutil.AssertNotEquals(4, (*candidates)[0].Id, t, "id 1. entry candidates")
	testutil.AssertNotEquals(4, (*candidates)[1].Id, t, "id 2. entry candidates")
	testutil.AssertEquals(VERSION_CANDIDATE, (*candidates)[0].Type, t, "type 1. entry candidates")
	testutil.AssertEquals(VERSION_CANDIDATE, (*candidates)[1].Type, t, "type 2. entry candidates")
}

func TestDetermineCandidatesPatchAllDifferentWithSnapshot(t *testing.T) {
	initCandidateTest()

	candidatesConf.NumberOfPatchVersionsToKeep = 1
	candidateVersionTwo.Name = "1.0.1-SNAPSHOT"
	candidateVersionThreee.Name = "1.0.2"

	candidates, err := DetermineCandidates(&candidatesConf)

	testutil.AssertNil(err, t, "err")
	testutil.AssertNotNil(candidates, t, "candidates")
	testutil.AssertEquals(1, len(*candidates), t, "len candidates")
	testutil.AssertEquals(2, (*candidates)[0].Id, t, "id 1. entry candidates")
	testutil.AssertEquals(VERSION_CANDIDATE, (*candidates)[0].Type, t, "type 1. entry candidates")
}

func TestDetermineCandidatesToManyVersionParts(t *testing.T) {
	initCandidateTest()

	candidatesConf.NumberOfMajorVersionsToKeep = 1
	candidateVersionTwo.Name = "2.0.0.1"
	candidateVersionThreee.Name = "3.0.0"

	candidates, err := DetermineCandidates(&candidatesConf)

	testutil.AssertNotNil(err, t, "err")
	testutil.AssertEquals("there are more items than 'major.minor.patch' or 'major.minor.patch-SNAPSHOT' at version name '2.0.0.1' with id 3", err.Error(), t, "err message")
	testutil.AssertNil(candidates, t, "candidates")
}

func TestDetermineCandidatesMajorAllDifferentMissingMinorPatch(t *testing.T) {
	initCandidateTest()

	candidatesConf.NumberOfMajorVersionsToKeep = 1
	candidateVersionTwo.Name = "2.0"
	candidateVersionThreee.Name = "3"

	candidates, err := DetermineCandidates(&candidatesConf)

	testutil.AssertNil(err, t, "err")
	testutil.AssertNotNil(candidates, t, "candidates")
	testutil.AssertEquals(2, len(*candidates), t, "len candidates")
	testutil.AssertNotEquals(4, (*candidates)[0].Id, t, "id 1. entry candidates")
	testutil.AssertNotEquals(4, (*candidates)[1].Id, t, "id 2. entry candidates")
	testutil.AssertEquals(VERSION_CANDIDATE, (*candidates)[0].Type, t, "type 1. entry candidates")
	testutil.AssertEquals(VERSION_CANDIDATE, (*candidates)[1].Type, t, "type 2. entry candidates")
}

func TestDetermineCandidatesMajorInvalidMajor(t *testing.T) {
	initCandidateTest()

	candidatesConf.NumberOfMajorVersionsToKeep = 1
	candidateVersionTwo.Name = "2.0.0"
	candidateVersionThreee.Name = "3a.0.0"

	candidates, err := DetermineCandidates(&candidatesConf)

	testutil.AssertNotNil(err, t, "err")
	testutil.AssertEquals("failed to format major version to int at version name '3a.0.0' with id 4: strconv.Atoi: parsing \"3a\": invalid syntax", err.Error(), t, "err message")
	testutil.AssertNil(candidates, t, "candidates")
}

func TestDetermineCandidatesMajorInvalidMinor(t *testing.T) {
	initCandidateTest()

	candidatesConf.NumberOfMajorVersionsToKeep = 1
	candidateVersionTwo.Name = "2.0.0"
	candidateVersionThreee.Name = "3.b.0"

	candidates, err := DetermineCandidates(&candidatesConf)

	testutil.AssertNotNil(err, t, "err")
	testutil.AssertEquals("failed to format minor version to int at version name '3.b.0' with id 4: strconv.Atoi: parsing \"b\": invalid syntax", err.Error(), t, "err message")
	testutil.AssertNil(candidates, t, "candidates")
}

func TestDetermineCandidatesMajorInvalidPatch(t *testing.T) {
	initCandidateTest()

	candidatesConf.NumberOfMajorVersionsToKeep = 1
	candidateVersionTwo.Name = "2.0.0"
	candidateVersionThreee.Name = "3.0.c"

	candidates, err := DetermineCandidates(&candidatesConf)

	testutil.AssertNotNil(err, t, "err")
	testutil.AssertEquals("failed to format patch version to int at version name '3.0.c' with id 4: strconv.Atoi: parsing \"c\": invalid syntax", err.Error(), t, "err message")
	testutil.AssertNil(candidates, t, "candidates")
}

func TestDetermineCandidatesGetVersionsWithError(t *testing.T) {
	initCandidateTest()

	VersionsGetExecutor = func(config *config.Config) (*[]github_model.Version, error) {
		return nil, errors.New("TestError")
	}

	candidatesConf.DeleteSnapshots = true

	candidates, err := DetermineCandidates(&candidatesConf)

	testutil.AssertNotNil(err, t, "err")
	testutil.AssertEquals("TestError", err.Error(), t, "err message")
	testutil.AssertNil(candidates, t, "candidates")
}

func TestDetermineCandidatesDeletePackage(t *testing.T) {
	initCandidateTest()

	candidatesConf.DeleteSnapshots = true
	candidateVersionOne.Name = "1.0.0-SNAPSHOT"
	candidateVersionTwo.Name = "2.0.0-SNAPSHOT"
	candidateVersionThreee.Name = "3.0.0-SNAPSHOT"

	candidates, err := DetermineCandidates(&candidatesConf)

	testutil.AssertNil(err, t, "err")
	testutil.AssertNotNil(candidates, t, "candidates")
	testutil.AssertEquals(1, len(*candidates), t, "len candidates")
	testutil.AssertEquals(1, (*candidates)[0].Id, t, "id")
	testutil.AssertEquals(PACKAGE_CANDIDATE, (*candidates)[0].Type, t, "type")
}

func TestDetermineCandidatesGetPackageWithError(t *testing.T) {
	initCandidateTest()

	PackageGetExecutor = func(config *config.Config) (*github_model.UserPackage, error) {
		return nil, errors.New("TestError")
	}

	candidatesConf.DeleteSnapshots = true
	candidateVersionOne.Name = "1.0.0-SNAPSHOT"
	candidateVersionTwo.Name = "2.0.0-SNAPSHOT"
	candidateVersionThreee.Name = "3.0.0-SNAPSHOT"

	candidates, err := DetermineCandidates(&candidatesConf)

	testutil.AssertNotNil(err, t, "err")
	testutil.AssertEquals("TestError", err.Error(), t, "err message")
	testutil.AssertNil(candidates, t, "candidates")
}
