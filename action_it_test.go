package main

import (
	"os"
	"testing"

	"github.com/ma-vin/packages-action/config"
	"github.com/ma-vin/packages-action/service/github_model"
	"github.com/ma-vin/packages-action/testutil"
)

func unsetEnv() {
	os.Unsetenv(config.ENV_NAME_GITHUB_REST_API_URL)
	os.Unsetenv(config.ENV_NAME_ORGANIZATION)
	os.Unsetenv(config.ENV_NAME_USER)
	os.Unsetenv(config.ENV_NAME_PACKAGE_TYPE)
	os.Unsetenv(config.ENV_NAME_PACKAGE_NAME)
	os.Unsetenv(config.ENV_NAME_VERSION_NAME_TO_DELETE)
	os.Unsetenv(config.ENV_NAME_DELETE_SNAPSHOTS)
	os.Unsetenv(config.ENV_NAME_NUMBER_MAJOR_TO_KEEP)
	os.Unsetenv(config.ENV_NAME_NUMBER_MINOR_TO_KEEP)
	os.Unsetenv(config.ENV_NAME_NUMBER_PATCH_TO_KEEP)
	os.Unsetenv(config.ENV_NAME_GITHUB_TOKEN)
	os.Unsetenv(config.ENV_NAME_DRY_RUN)
}

func createTestPackage() *github_model.UserPackage {
	return &github_model.UserPackage{Id: 1, Name: "DummyPackage", CreatedAt: "2024-03-12T20:00:00Z", UpdatedAt: "2024-03-20:00:00Z"}
}

func createTestVersions(snapshots bool) *[]github_model.Version {
	nameOne := "1.0.0"
	nameTwo := "2.1.0"
	nameThree := "3.0.1"
	if snapshots {
		nameOne = nameOne + "-SNAPSHOT"
		nameTwo = nameTwo + "-SNAPSHOT"
		nameThree = nameThree + "-SNAPSHOT"
	}

	candidateVersionOne := github_model.Version{Id: 2, Name: nameOne, Description: "First Version", CreatedAt: "2024-03-12T20:00:00Z", UpdatedAt: "2024-03-13T16:00:00Z"}
	candidateVersionTwo := github_model.Version{Id: 3, Name: nameTwo, Description: "Second Version", CreatedAt: "2024-03-13T20:00:00Z", UpdatedAt: "2024-03-14T16:00:00Z"}
	candidateVersionThree := github_model.Version{Id: 4, Name: nameThree, Description: "Third Version", CreatedAt: "2024-03-14T20:00:00Z", UpdatedAt: "2024-03-20:00:00Z"}

	return &[]github_model.Version{candidateVersionOne, candidateVersionTwo, candidateVersionThree}
}

func TestMainDeleteVersionsDryRun(t *testing.T) {
	unsetEnv()

	mockServerUrl := testutil.CreateAndStartMock("Ma-Vin", config.MAVEN, "DummyPackage", createTestVersions(false), createTestPackage())
	defer testutil.StopMock()

	os.Setenv(config.ENV_NAME_GITHUB_REST_API_URL, mockServerUrl)
	os.Setenv(config.ENV_NAME_USER, "Ma-Vin")
	os.Setenv(config.ENV_NAME_PACKAGE_TYPE, config.MAVEN)
	os.Setenv(config.ENV_NAME_PACKAGE_NAME, "DummyPackage")
	os.Setenv(config.ENV_NAME_NUMBER_MAJOR_TO_KEEP, "1")
	os.Setenv(config.ENV_NAME_GITHUB_TOKEN, "abcdef123")

	main()

	testutil.AssertEquals(1, testutil.GetUserPackageVersionsCounter, t, "Count of GetUserPackageVersions")
	testutil.AssertEquals(0, testutil.DeleteUserPackageVersionCounter, t, "Count of DeleteUserPackageVersion")
	testutil.AssertEquals(0, testutil.GetUserPackageCounter, t, "Count of GetUserPackage")
	testutil.AssertEquals(0, testutil.DeleteUserPackageCounter, t, "Count of DeleteUserPackage")
}

func TestMainDeleteVersionsRealRun(t *testing.T) {
	unsetEnv()

	mockServerUrl := testutil.CreateAndStartMock("Ma-Vin", config.MAVEN, "DummyPackage", createTestVersions(false), createTestPackage())
	defer testutil.StopMock()

	os.Setenv(config.ENV_NAME_GITHUB_REST_API_URL, mockServerUrl)
	os.Setenv(config.ENV_NAME_USER, "Ma-Vin")
	os.Setenv(config.ENV_NAME_PACKAGE_TYPE, config.MAVEN)
	os.Setenv(config.ENV_NAME_PACKAGE_NAME, "DummyPackage")
	os.Setenv(config.ENV_NAME_NUMBER_MAJOR_TO_KEEP, "1")
	os.Setenv(config.ENV_NAME_GITHUB_TOKEN, "abcdef123")
	os.Setenv(config.ENV_NAME_DRY_RUN, "false")

	main()

	testutil.AssertEquals(1, testutil.GetUserPackageVersionsCounter, t, "Count of GetUserPackageVersions")
	testutil.AssertEquals(2, testutil.DeleteUserPackageVersionCounter, t, "Count of DeleteUserPackageVersion")
	testutil.AssertEquals(0, testutil.GetUserPackageCounter, t, "Count of GetUserPackage")
	testutil.AssertEquals(0, testutil.DeleteUserPackageCounter, t, "Count of DeleteUserPackage")
}

func TestMainDeleteAllVersionsDryRun(t *testing.T) {
	unsetEnv()

	mockServerUrl := testutil.CreateAndStartMock("Ma-Vin", config.MAVEN, "DummyPackage", createTestVersions(true), createTestPackage())
	defer testutil.StopMock()

	os.Setenv(config.ENV_NAME_GITHUB_REST_API_URL, mockServerUrl)
	os.Setenv(config.ENV_NAME_USER, "Ma-Vin")
	os.Setenv(config.ENV_NAME_PACKAGE_TYPE, config.MAVEN)
	os.Setenv(config.ENV_NAME_PACKAGE_NAME, "DummyPackage")
	os.Setenv(config.ENV_NAME_DELETE_SNAPSHOTS, "true")
	os.Setenv(config.ENV_NAME_GITHUB_TOKEN, "abcdef123")

	main()

	testutil.AssertEquals(1, testutil.GetUserPackageVersionsCounter, t, "Count of GetUserPackageVersions")
	testutil.AssertEquals(0, testutil.DeleteUserPackageVersionCounter, t, "Count of DeleteUserPackageVersion")
	testutil.AssertEquals(1, testutil.GetUserPackageCounter, t, "Count of GetUserPackage")
	testutil.AssertEquals(0, testutil.DeleteUserPackageCounter, t, "Count of DeleteUserPackage")
}

func TestMainDeleteAllVersionsRealRun(t *testing.T) {
	unsetEnv()

	mockServerUrl := testutil.CreateAndStartMock("Ma-Vin", config.MAVEN, "DummyPackage", createTestVersions(true), createTestPackage())
	defer testutil.StopMock()

	os.Setenv(config.ENV_NAME_GITHUB_REST_API_URL, mockServerUrl)
	os.Setenv(config.ENV_NAME_USER, "Ma-Vin")
	os.Setenv(config.ENV_NAME_PACKAGE_TYPE, config.MAVEN)
	os.Setenv(config.ENV_NAME_PACKAGE_NAME, "DummyPackage")
	os.Setenv(config.ENV_NAME_DELETE_SNAPSHOTS, "true")
	os.Setenv(config.ENV_NAME_GITHUB_TOKEN, "abcdef123")
	os.Setenv(config.ENV_NAME_DRY_RUN, "false")

	main()

	testutil.AssertEquals(1, testutil.GetUserPackageVersionsCounter, t, "Count of GetUserPackageVersions")
	testutil.AssertEquals(0, testutil.DeleteUserPackageVersionCounter, t, "Count of DeleteUserPackageVersion")
	testutil.AssertEquals(1, testutil.GetUserPackageCounter, t, "Count of GetUserPackage")
	testutil.AssertEquals(1, testutil.DeleteUserPackageCounter, t, "Count of DeleteUserPackage")
}
