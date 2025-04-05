package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ma-vin/packages-action/config"
	"github.com/ma-vin/packages-action/service/github_model"
)

const gitHubModelVersion string = "2022-11-28"
const gitHubModelJsonType string = "application/vnd.github+json"

const users_url_part string = "users"
const packages_url_part string = "packages"
const versions_url_part string = "versions"

type queryParameter struct {
	name  string
	value string
}

type ClientExecutor func(*http.Client, *http.Request) (*http.Response, error)

var ClientRestExecutor ClientExecutor = initClientExector()

func initClientExector() ClientExecutor {
	return func(c *http.Client, req *http.Request) (*http.Response, error) {
		return c.Do(req)
	}
}

func InitAllGitHubRest() {
	ClientRestExecutor = initClientExector()
}

// calls GitHub rest api to get all packages of a certain type and user.
// /users/{username}/packages
func GetUserPackages(configuration *config.Config) (*[]github_model.UserPackage, error) {
	url := concatUrl(configuration.GitHubRestUrl, users_url_part, configuration.User, packages_url_part)
	response, err := get(url, configuration, []queryParameter{{name: "package_type", value: configuration.PackageType}})

	if err != nil {
		return nil, err
	}

	var userPackages []github_model.UserPackage
	err = mapJsonResponse(response, &userPackages, configuration)

	if err != nil {
		return nil, err
	}

	return &userPackages, nil
}

// calls GitHub rest api to get a package of a certain type and user.
// users/{username}/packages/{package_type}/{package_name}
func GetUserPackage(packageName string, configuration *config.Config) (*github_model.UserPackage, error) {
	url := concatUrl(configuration.GitHubRestUrl, users_url_part, configuration.User, packages_url_part, configuration.PackageType, packageName)
	response, err := get(url, configuration, nil)

	if err != nil {
		return nil, err
	}

	var userPackage github_model.UserPackage
	err = mapJsonResponse(response, &userPackage, configuration)

	if err != nil {
		return nil, err
	}
	return &userPackage, nil
}

// calls GitHub rest api to get a package of a certain type and user
// /users/{username}/packages/{package_type}/{package_name}
func DeleteUserPackage(packageName string, configuration *config.Config) error {
	url := concatUrl(configuration.GitHubRestUrl, users_url_part, configuration.User, packages_url_part, configuration.PackageType, packageName)

	response, err := delete(url, configuration, nil)

	if err != nil {
		return err
	}
	return checkResponseStatusCode(response, configuration)
}

// calls GitHub rest api to get all versions of a certain package, type and user.
// /users/{username}/packages/{package_type}/{package_name}/versions
func GetUserPackageVersions(packageName string, configuration *config.Config) (*[]github_model.Version, error) {
	url := concatUrl(configuration.GitHubRestUrl, users_url_part, configuration.User, packages_url_part, configuration.PackageType, packageName, versions_url_part)
	response, err := get(url, configuration, nil)

	if err != nil {
		return nil, err
	}

	var versions []github_model.Version
	err = mapJsonResponse(response, &versions, configuration)

	if err != nil {
		return nil, err
	}
	return &versions, nil
}

// calls GitHub rest api to get a version of a certain package, type and user.
// /users/{username}/packages/{package_type}/{package_name}/versions/{package_version_id}
func GetUserPackageVersion(packageName string, versionId int, configuration *config.Config) (*github_model.Version, error) {
	url := concatUrl(configuration.GitHubRestUrl, users_url_part, configuration.User, packages_url_part, configuration.PackageType, packageName, versions_url_part, strconv.Itoa(versionId))
	response, err := get(url, configuration, nil)

	if err != nil {
		return nil, err
	}

	var version github_model.Version
	err = mapJsonResponse(response, &version, configuration)

	if err != nil {
		return nil, err
	}
	return &version, nil
}

// calls GitHub rest api to get a package of a certain type and user
// /users/{username}/packages/{package_type}/{package_name}/versions/{package_version_id}
func DeleteUserPackageVersion(packageName string, versionId int, configuration *config.Config) error {
	url := concatUrl(configuration.GitHubRestUrl, users_url_part, configuration.User, packages_url_part, configuration.PackageType, packageName, versions_url_part, strconv.Itoa(versionId))

	response, err := delete(url, configuration, nil)

	if err != nil {
		return err
	}
	return checkResponseStatusCode(response, configuration)
}

// maps the the json body of a response to a given target object
func mapJsonResponse(response *http.Response, target any, configuration *config.Config) error {
	err := checkResponseStatusCode(response, configuration)
	if err != nil {
		return err
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

// creates the client, request, adds header elemets and url query parameters before sending. TLS is not configured explicitly since tls.Config uses TLS1.2 as MinVersion
func executeRequestWithoutBody(operation string, url string, configuration *config.Config, parameters []queryParameter) (*http.Response, error) {
	c := http.Client{Timeout: time.Duration(configuration.Timeout) * time.Second}

	req, err := http.NewRequest(operation, url, nil)
	if err != nil {
		return nil, err
	}

	addHeader(req, configuration)
	addUrlQueryParameters(req, &parameters)

	return ClientRestExecutor(&c, req)
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

// checks if the response has a failure status code and creates in this case an error
func checkResponseStatusCode(response *http.Response, configuration *config.Config) error {
	if response.StatusCode >= 400 {
		logHeader(&response.Header, "response header", configuration)
		if response.Request != nil {
			logHeader(&response.Request.Header, "request header", configuration)
			return fmt.Errorf("an error status code occured at %s '%s': %d - %s", response.Request.Method, response.Request.URL, response.StatusCode, http.StatusText(response.StatusCode))
		}
		return fmt.Errorf("an error status code occured: %d - %s", response.StatusCode, http.StatusText(response.StatusCode))
	}
	return nil
}

func logHeader(response *http.Header, headerName string, configuration *config.Config) {
	if !configuration.Debug {
		return
	}
	fmt.Println(headerName)
	for name, values := range *response {
		if name == "Authorization" {
			fmt.Println(name, "***")
			continue
		}
		for _, value := range values {
			fmt.Println(name, value)
		}
	}
}
