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

type queryParameter struct {
	name  string
	value string
}

// calls GitHub rest api to get all packages of a certain type and user
func GetAndPrintUserPackages(configuration *config.Config) []github_model.UserPackage {

	response, err := get(gitHubUserUrl+"/"+configuration.User+"/packages", configuration, []queryParameter{{name: "package_type", value: config.MapFromPackageType(configuration.PackageType)}})

	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode >= 400 {
		log.Fatal("StatusCode: ", response.StatusCode)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var packages []github_model.UserPackage

	errUnmarshall := json.Unmarshal(responseData, &packages)
	if errUnmarshall != nil {
		log.Fatal(errUnmarshall)
	}

	log.Println("Number of packages:", len(packages))
	for i, p := range packages {
		log.Println(i, p.Name, p.Id)
	}

	return packages
}

func get(url string, configuration *config.Config, parameters []queryParameter) (resp *http.Response, err error) {
	c := http.Client{Timeout: time.Duration(1) * time.Second}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Accept", `application/vnd.github+json`)
	req.Header.Add("Authorization", "Bearer "+configuration.GithubToken)
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	q := req.URL.Query()
	for _, p := range parameters {
		q.Add(p.name, p.value)
	}
	req.URL.RawQuery = q.Encode()

	return c.Do(req)
}
