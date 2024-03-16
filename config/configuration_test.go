package config

import (
	"os"
	"testing"

	"github.com/ma-vin/packages-action/testutil"
)

func unsetEnv() {
	os.Unsetenv(ENV_NAME_ORGANIZATION)
	os.Unsetenv(ENV_NAME_USER)
	os.Unsetenv(ENV_NAME_PACKAGE_TYPE)
	os.Unsetenv(ENV_NAME_PACKAGE_NAME)
	os.Unsetenv(ENV_NAME_VERSION_NAME_TO_DELETE)
	os.Unsetenv(ENV_NAME_DELETE_SNAPSHOTS)
	os.Unsetenv(ENV_NAME_NUMBER_MAJOR_TO_KEEP)
	os.Unsetenv(ENV_NAME_NUMBER_MINOR_TO_KEEP)
	os.Unsetenv(ENV_NAME_NUMBER_PATCH_TO_KEEP)
	os.Unsetenv(ENV_NAME_GITHUB_TOKEN)
	os.Unsetenv(ENV_NAME_DRY_RUN)
}

func TestReadConfigurationUserAndOrganization(t *testing.T) {
	unsetEnv()

	os.Setenv(ENV_NAME_ORGANIZATION, "Ma-Vin-Org")
	os.Setenv(ENV_NAME_USER, "Ma-Vin")
	os.Setenv(ENV_NAME_PACKAGE_TYPE, MAVEN)
	os.Setenv(ENV_NAME_PACKAGE_NAME, "packages-action-app")
	os.Setenv(ENV_NAME_DELETE_SNAPSHOTS, "TRUE")
	os.Setenv(ENV_NAME_GITHUB_TOKEN, "abcdef123")

	conf := ReadConfiguration()

	testutil.AssertNil(conf, t, "conf")
}

func TestReadConfigurationMissingUserAndOrganization(t *testing.T) {
	unsetEnv()

	os.Setenv(ENV_NAME_PACKAGE_TYPE, MAVEN)
	os.Setenv(ENV_NAME_PACKAGE_NAME, "packages-action-app")
	os.Setenv(ENV_NAME_DELETE_SNAPSHOTS, "TRUE")
	os.Setenv(ENV_NAME_GITHUB_TOKEN, "abcdef123")

	conf := ReadConfiguration()

	testutil.AssertNil(conf, t, "conf")
}

func TestReadConfigurationUser(t *testing.T) {
	unsetEnv()

	os.Setenv(ENV_NAME_USER, "Ma-Vin")
	os.Setenv(ENV_NAME_PACKAGE_TYPE, MAVEN)
	os.Setenv(ENV_NAME_PACKAGE_NAME, "packages-action-app")
	os.Setenv(ENV_NAME_DELETE_SNAPSHOTS, "TRUE")
	os.Setenv(ENV_NAME_GITHUB_TOKEN, "abcdef123")

	conf := ReadConfiguration()

	testutil.AssertNotNil(conf, t, "conf")
	testutil.AssertEquals("", conf.Organization, t, "organization")
	testutil.AssertEquals("Ma-Vin", conf.User, t, "user")
	testutil.AssertEquals("maven", conf.PackageType, t, "package type")
	testutil.AssertEquals("packages-action-app", conf.PackageName, t, "package name")
	testutil.AssertEquals("", conf.VersionNameToDelete, t, "version name to delete")
	testutil.AssertEquals(true, conf.DeleteSnapshots, t, "delete snapshots")
	testutil.AssertEquals(-1, conf.NumberOfMajorVersionsToKeep, t, "number of major versions")
	testutil.AssertEquals(-1, conf.NumberOfMinorVersionsToKeep, t, "number of minor versions")
	testutil.AssertEquals(-1, conf.NumberOfPatchVersionsToKeep, t, "number of patch versions")
	testutil.AssertEquals("abcdef123", conf.GithubToken, t, "github token")
	testutil.AssertEquals(true, conf.DryRun, t, "dry run")
}

func TestReadConfigurationOrganization(t *testing.T) {
	unsetEnv()

	os.Setenv(ENV_NAME_ORGANIZATION, "Ma-Vin-Org")
	os.Setenv(ENV_NAME_PACKAGE_TYPE, MAVEN)
	os.Setenv(ENV_NAME_PACKAGE_NAME, "packages-action-app")
	os.Setenv(ENV_NAME_DELETE_SNAPSHOTS, "TRUE")
	os.Setenv(ENV_NAME_GITHUB_TOKEN, "abcdef123")

	conf := ReadConfiguration()

	testutil.AssertNotNil(conf, t, "conf")
	testutil.AssertEquals("Ma-Vin-Org", conf.Organization, t, "organization")
	testutil.AssertEquals("", conf.User, t, "user")
	testutil.AssertEquals("maven", conf.PackageType, t, "package type")
	testutil.AssertEquals("packages-action-app", conf.PackageName, t, "package name")
	testutil.AssertEquals("", conf.VersionNameToDelete, t, "version name to delete")
	testutil.AssertEquals(true, conf.DeleteSnapshots, t, "delete snapshots")
	testutil.AssertEquals(-1, conf.NumberOfMajorVersionsToKeep, t, "number of major versions")
	testutil.AssertEquals(-1, conf.NumberOfMinorVersionsToKeep, t, "number of minor versions")
	testutil.AssertEquals(-1, conf.NumberOfPatchVersionsToKeep, t, "number of patch versions")
	testutil.AssertEquals("abcdef123", conf.GithubToken, t, "github token")
	testutil.AssertEquals(true, conf.DryRun, t, "dry run")
}

func TestReadConfigurationUnkownPackageType(t *testing.T) {
	unsetEnv()

	os.Setenv(ENV_NAME_USER, "Ma-Vin")
	os.Setenv(ENV_NAME_PACKAGE_TYPE, "abc")
	os.Setenv(ENV_NAME_PACKAGE_NAME, "packages-action-app")
	os.Setenv(ENV_NAME_DELETE_SNAPSHOTS, "TRUE")
	os.Setenv(ENV_NAME_GITHUB_TOKEN, "abcdef123")

	conf := ReadConfiguration()

	testutil.AssertNil(conf, t, "conf")
}

func TestReadConfigurationMissingPackageName(t *testing.T) {
	unsetEnv()

	os.Setenv(ENV_NAME_USER, "Ma-Vin")
	os.Setenv(ENV_NAME_PACKAGE_TYPE, MAVEN)
	os.Setenv(ENV_NAME_DELETE_SNAPSHOTS, "TRUE")
	os.Setenv(ENV_NAME_GITHUB_TOKEN, "abcdef123")

	conf := ReadConfiguration()

	testutil.AssertNil(conf, t, "conf")
}
func TestReadConfigurationMissingToken(t *testing.T) {
	unsetEnv()

	os.Setenv(ENV_NAME_USER, "Ma-Vin")
	os.Setenv(ENV_NAME_PACKAGE_TYPE, MAVEN)
	os.Setenv(ENV_NAME_PACKAGE_NAME, "packages-action-app")
	os.Setenv(ENV_NAME_DELETE_SNAPSHOTS, "TRUE")

	conf := ReadConfiguration()

	testutil.AssertNil(conf, t, "conf")
}

func TestReadConfigurationNothingToDelete(t *testing.T) {
	unsetEnv()

	os.Setenv(ENV_NAME_USER, "Ma-Vin")
	os.Setenv(ENV_NAME_PACKAGE_TYPE, MAVEN)
	os.Setenv(ENV_NAME_PACKAGE_NAME, "packages-action-app")
	os.Setenv(ENV_NAME_GITHUB_TOKEN, "abcdef123")

	conf := ReadConfiguration()

	testutil.AssertNil(conf, t, "conf")
}

func TestReadConfigurationDeleteSNaphostNegative(t *testing.T) {
	unsetEnv()

	os.Setenv(ENV_NAME_USER, "Ma-Vin")
	os.Setenv(ENV_NAME_PACKAGE_TYPE, MAVEN)
	os.Setenv(ENV_NAME_PACKAGE_NAME, "packages-action-app")
	os.Setenv(ENV_NAME_DELETE_SNAPSHOTS, "FALSE")
	os.Setenv(ENV_NAME_GITHUB_TOKEN, "abcdef123")

	conf := ReadConfiguration()

	testutil.AssertNil(conf, t, "conf")
}

func TestReadConfigurationUserAllSet(t *testing.T) {
	unsetEnv()

	os.Setenv(ENV_NAME_USER, "Ma-Vin")
	os.Setenv(ENV_NAME_PACKAGE_TYPE, MAVEN)
	os.Setenv(ENV_NAME_PACKAGE_NAME, "packages-action-app")
	os.Setenv(ENV_NAME_VERSION_NAME_TO_DELETE, "1.2.3-SNAPSHOT")
	os.Setenv(ENV_NAME_DELETE_SNAPSHOTS, "TRUE")
	os.Setenv(ENV_NAME_NUMBER_MAJOR_TO_KEEP, "3")
	os.Setenv(ENV_NAME_NUMBER_MINOR_TO_KEEP, "2")
	os.Setenv(ENV_NAME_NUMBER_PATCH_TO_KEEP, "1")
	os.Setenv(ENV_NAME_GITHUB_TOKEN, "abcdef123")
	os.Setenv(ENV_NAME_DRY_RUN, "false")

	conf := ReadConfiguration()

	testutil.AssertNotNil(conf, t, "conf")
	testutil.AssertEquals("", conf.Organization, t, "organization")
	testutil.AssertEquals("Ma-Vin", conf.User, t, "user")
	testutil.AssertEquals("maven", conf.PackageType, t, "package type")
	testutil.AssertEquals("packages-action-app", conf.PackageName, t, "package name")
	testutil.AssertEquals("1.2.3-SNAPSHOT", conf.VersionNameToDelete, t, "version name to delete")
	testutil.AssertEquals(true, conf.DeleteSnapshots, t, "delete snapshots")
	testutil.AssertEquals(3, conf.NumberOfMajorVersionsToKeep, t, "number of major versions")
	testutil.AssertEquals(2, conf.NumberOfMinorVersionsToKeep, t, "number of minor versions")
	testutil.AssertEquals(1, conf.NumberOfPatchVersionsToKeep, t, "number of patch versions")
	testutil.AssertEquals("abcdef123", conf.GithubToken, t, "github token")
	testutil.AssertEquals(false, conf.DryRun, t, "dry run")
}

func TestReadConfigurationInvalidInt(t *testing.T) {
	unsetEnv()

	os.Setenv(ENV_NAME_USER, "Ma-Vin")
	os.Setenv(ENV_NAME_PACKAGE_TYPE, MAVEN)
	os.Setenv(ENV_NAME_PACKAGE_NAME, "packages-action-app")
	os.Setenv(ENV_NAME_VERSION_NAME_TO_DELETE, "1.2.3-SNAPSHOT")
	os.Setenv(ENV_NAME_DELETE_SNAPSHOTS, "TRUE")
	os.Setenv(ENV_NAME_NUMBER_MAJOR_TO_KEEP, "-3")
	os.Setenv(ENV_NAME_NUMBER_MINOR_TO_KEEP, "2abc")
	os.Setenv(ENV_NAME_NUMBER_PATCH_TO_KEEP, "1")
	os.Setenv(ENV_NAME_GITHUB_TOKEN, "abcdef123")

	conf := ReadConfiguration()

	testutil.AssertNotNil(conf, t, "conf")
	testutil.AssertEquals("", conf.Organization, t, "organization")
	testutil.AssertEquals("Ma-Vin", conf.User, t, "user")
	testutil.AssertEquals("maven", conf.PackageType, t, "package type")
	testutil.AssertEquals("packages-action-app", conf.PackageName, t, "package name")
	testutil.AssertEquals("1.2.3-SNAPSHOT", conf.VersionNameToDelete, t, "version name to delete")
	testutil.AssertEquals(true, conf.DeleteSnapshots, t, "delete snapshots")
	testutil.AssertEquals(-1, conf.NumberOfMajorVersionsToKeep, t, "number of major versions")
	testutil.AssertEquals(-1, conf.NumberOfMinorVersionsToKeep, t, "number of minor versions")
	testutil.AssertEquals(1, conf.NumberOfPatchVersionsToKeep, t, "number of patch versions")
	testutil.AssertEquals("abcdef123", conf.GithubToken, t, "github token")
	testutil.AssertEquals(true, conf.DryRun, t, "dry run")
}
