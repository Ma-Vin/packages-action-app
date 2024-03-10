package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ma-vin/packages-action/config"
	"github.com/ma-vin/packages-action/service/github_model"
)

const gitHubUrl string = "https://api.github.com"
const gitHubUserUrl string = gitHubUrl + "/users"
const gitHubModelVersion string = "2022-11-28"
const gitHubModelJsonType string = "application/vnd.github+json"

type queryParameter struct {
	name  string
	value string
}

// calls GitHub rest api to get all packages of a certain type and user.
// /users/{username}/packages
func GetUserPackages(configuration *config.Config) (*[]github_model.UserPackage, error) {
	url := concatUrl(gitHubUserUrl, configuration.User, "packages")
	response, err := get(url, configuration, []queryParameter{{name: "package_type", value: configuration.PackageType}})

	if err != nil {
		return nil, err
	}

	var userPackages []github_model.UserPackage

	return &userPackages, mapJsonResponse(response, &userPackages)
}

// calls GitHub rest api to get a package of a certain type and user.
// users/{username}/packages/{package_type}/{package_name}
func GetUserPackage(packageName string, configuration *config.Config) (*github_model.UserPackage, error) {
	url := concatUrl(gitHubUserUrl, configuration.User, "packages", configuration.PackageType, packageName)
	response, err := get(url, configuration, nil)

	if err != nil {
		return nil, err
	}

	var userPackage github_model.UserPackage

	return &userPackage, mapJsonResponse(response, &userPackage)
}

// calls GitHub rest api to get a package of a certain type and user
// /users/{username}/packages/{package_type}/{package_name}
func DeleteUserPackage(packageName string, configuration *config.Config) error {
	url := concatUrl(gitHubUserUrl, configuration.User, "packages", configuration.PackageType, packageName)

	response, err := delete(url, configuration, nil)

	if err != nil {
		return err
	}
	if response.StatusCode >= 400 {
		return fmt.Errorf("an error status code occured: %d - %s", response.StatusCode, http.StatusText(response.StatusCode))
	}
	return nil
}

// calls GitHub rest api to get all versions of a certain package, type and user.
// /users/{username}/packages/{package_type}/{package_name}/versions
func GetUserPackageVersions(packageName string, configuration *config.Config) (*[]github_model.Version, error) {
	url := concatUrl(gitHubUserUrl, configuration.User, "packages", configuration.PackageType, packageName, "versions")
	response, err := get(url, configuration, nil)

	if err != nil {
		return nil, err
	}

	var versions []github_model.Version

	return &versions, mapJsonResponse(response, &versions)
}

// calls GitHub rest api to get a version of a certain package, type and user.
// /users/{username}/packages/{package_type}/{package_name}/versions/{package_version_id}
func GetUserPackageVersion(packageName string, versionId int, configuration *config.Config) (*github_model.Version, error) {
	url := concatUrl(gitHubUserUrl, configuration.User, "packages", configuration.PackageType, packageName, "versions", strconv.Itoa(versionId))
	response, err := get(url, configuration, nil)

	if err != nil {
		return nil, err
	}

	var version github_model.Version

	return &version, mapJsonResponse(response, &version)
}

// calls GitHub rest api to get all packages of a certain type and user.
// The result is printed to log
func GetAndPrintUserPackages(configuration *config.Config) (*[]github_model.UserPackage, error) {

	userPackages, err := GetUserPackages(configuration)

	if err != nil {
		return userPackages, err
	}

	log.Println("Number of packages:", len(*userPackages))
	for i, p := range *userPackages {
		log.Println(i+1, p.Name, p.Id)
	}

	return userPackages, nil
}

// calls GitHub rest api to get a package of a certain type and user.
// The result is printed to log
func GetAndPrintUserPackage(packageName string, configuration *config.Config) (*github_model.UserPackage, error) {

	userPackage, err := GetUserPackage(packageName, configuration)

	if err != nil {
		return userPackage, err
	}

	log.Println(userPackage.Name, userPackage.Id)

	return userPackage, nil
}

// calls GitHub rest api to get all versions of a certain package, type and user.
// The result is printed to log
func GetAndPrintUserPackageVersions(packageName string, configuration *config.Config) (*[]github_model.Version, error) {

	versions, err := GetUserPackageVersions(packageName, configuration)

	if err != nil {
		return versions, err
	}

	log.Println("Number of versions:", len(*versions))
	for i, p := range *versions {
		log.Println(i+1, p.Name, p.Id)
	}

	return versions, nil
}

// calls GitHub rest api to get a version of a certain package, type and user
// The result is printed to log
func GetAndPrintUserPackageVersion(packageName string, versionId int, configuration *config.Config) (*github_model.Version, error) {

	version, err := GetUserPackageVersion(packageName, versionId, configuration)

	if err != nil {
		return version, err
	}

	log.Println(version.Name, version.Id)

	return version, nil
}

// maps the the json body of a response to a given target object
func mapJsonResponse(response *http.Response, target any) error {
	if response.StatusCode >= 400 {
		return fmt.Errorf("an error status code occured: %d - %s", response.StatusCode, http.StatusText(response.StatusCode))
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(responseData, target)
	if err != nil {
		return err
	}

	return nil
}

// Executes a get rest call
func get(url string, configuration *config.Config, parameters []queryParameter) (*http.Response, error) {
	return executeRequestWithoutBody(http.MethodGet, url, configuration, parameters)
}

// Executes a delete rest call
func delete(url string, configuration *config.Config, parameters []queryParameter) (*http.Response, error) {
	return executeRequestWithoutBody(http.MethodDelete, url, configuration, parameters)
}

func executeRequestWithoutBody(operation string, url string, configuration *config.Config, parameters []queryParameter) (*http.Response, error) {
	c := http.Client{Timeout: time.Duration(1) * time.Second}

	req, err := http.NewRequest(operation, url, nil)
	if err != nil {
		return nil, err
	}

	addHeader(req, configuration)
	addUrlQueryParameters(req, &parameters)

	return c.Do(req)
}

// adds the default header elements for a call against GitHub rest api
func addHeader(req *http.Request, configuration *config.Config) {
	req.Header.Add("Accept", gitHubModelJsonType)
	req.Header.Add("Authorization", "Bearer "+configuration.GithubToken)
	req.Header.Add("X-GitHub-Api-Version", gitHubModelVersion)
}

// add query parameters to an url of a given request
func addUrlQueryParameters(req *http.Request, parameters *[]queryParameter) {
	q := req.URL.Query()
	for _, p := range *parameters {
		q.Add(p.name, p.value)
	}
	req.URL.RawQuery = q.Encode()
}

func concatUrl(urlParts ...string) string {
	var sb strings.Builder
	for i, urlPart := range urlParts {
		sb.WriteString(urlPart)
		if i+1 < len(urlParts) && !strings.HasSuffix(urlPart, "/") {
			sb.WriteString("/")
		}
	}
	return sb.String()
}
