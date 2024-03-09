package service

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
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

// calls GitHub rest api to get all packages of a certain type and user
func GetUserPackages(configuration *config.Config) (*[]github_model.UserPackage, error) {

	response, err := get(gitHubUserUrl+"/"+configuration.User+"/packages", configuration, []queryParameter{{name: "package_type", value: configuration.PackageType}})

	if err != nil {
		return nil, err
	}
	if response.StatusCode >= 400 {
		log.Println("An error status code occured: ", response.StatusCode)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var packages []github_model.UserPackage

	err = json.Unmarshal(responseData, &packages)
	if err != nil {
		return nil, err
	}

	return &packages, nil
}

// calls GitHub rest api to get all packages of a certain type and user.
// The result ist printed to log
func GetAndPrintUserPackages(configuration *config.Config) (*[]github_model.UserPackage, error) {

	packages, err := GetUserPackages(configuration)

	if err != nil {
		return packages, err
	}

	log.Println("Number of packages:", len(*packages))
	for i, p := range *packages {
		log.Println(i+1, p.Name, p.Id)
	}

	return packages, nil
}

// Executes a get rest call
func get(url string, configuration *config.Config, parameters []queryParameter) (*http.Response, error) {
	c := http.Client{Timeout: time.Duration(1) * time.Second}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", gitHubModelJsonType)
	req.Header.Add("Authorization", "Bearer "+configuration.GithubToken)
	req.Header.Add("X-GitHub-Api-Version", gitHubModelVersion)

	q := req.URL.Query()
	for _, p := range parameters {
		q.Add(p.name, p.value)
	}
	req.URL.RawQuery = q.Encode()

	return c.Do(req)
}
