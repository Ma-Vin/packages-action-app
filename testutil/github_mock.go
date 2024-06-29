package testutil

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/ma-vin/packages-action/service/github_model"
)

const gitHubModelJsonType string = "application/vnd.github+json"

var server *httptest.Server
var versionsData *[]github_model.Version
var packageData *github_model.UserPackage

var GetUserPackageVersionsCounter int
var DeleteUserPackageVersionCounter int
var GetUserPackageCounter int
var DeleteUserPackageCounter int
var GetAllUserPackagesCounter int

func CreateAndStartMock(userName string, packageType string, packageName string, versions *[]github_model.Version, userPackage *github_model.UserPackage) string {
	versionsData = versions
	packageData = userPackage

	mux := http.NewServeMux()

	getVersionsUrl := fmt.Sprintf("/users/%s/packages/%s/%s/versions", userName, packageType, packageName)
	mux.HandleFunc(getVersionsUrl, getUserPackageVersionsHandler)
	GetUserPackageVersionsCounter = 0

	deleteVersionUrl := fmt.Sprintf("/users/%s/packages/%s/%s/versions/{id}", userName, packageType, packageName)
	mux.HandleFunc(deleteVersionUrl, deleteUserPackageVersionHandler)
	DeleteUserPackageVersionCounter = 0

	getPackageUrl := fmt.Sprintf("/users/%s/packages/%s/%s", userName, packageType, packageName)
	mux.HandleFunc(getPackageUrl, getOrDeleteUserPackageHandler)
	GetUserPackageCounter = 0
	DeleteUserPackageCounter = 0

	getAllPackagesUrl := fmt.Sprintf("/users/%s/packages", userName)
	mux.HandleFunc(getAllPackagesUrl, getAllUserPackagesHandler)
	GetAllUserPackagesCounter = 0

	server = httptest.NewServer(mux)
	log.Println("Mock - server started")

	return server.URL
}

func StopMock() {
	server.Close()
	log.Println("Mock - server stopped")
}

func getUserPackageVersionsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Mock - getUserPackageVersionsHandler %s '%s'", r.Method, r.URL)
	if r.Method == http.MethodGet {
		GetUserPackageVersionsCounter++
		w.Header().Set("Content-Type", gitHubModelJsonType)
		json.NewEncoder(w).Encode(versionsData)
	} else {
		w.WriteHeader(500)
	}
}

func deleteUserPackageVersionHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Mock - deleteUserPackageVersionHandler %s '%s'", r.Method, r.URL)
	if r.Method == http.MethodDelete {
		DeleteUserPackageVersionCounter++
		w.Header().Set("Content-Type", gitHubModelJsonType)
		w.WriteHeader(204)
	} else {
		w.WriteHeader(500)
	}
}

func getOrDeleteUserPackageHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Mock - getOrDeleteUserPackageHandler %s '%s'", r.Method, r.URL)
	switch r.Method {
	case http.MethodGet:
		GetUserPackageCounter++
		w.Header().Set("Content-Type", gitHubModelJsonType)
		json.NewEncoder(w).Encode(packageData)
	case http.MethodDelete:
		DeleteUserPackageCounter++
		w.Header().Set("Content-Type", gitHubModelJsonType)
		w.WriteHeader(204)
	default:
		w.WriteHeader(500)
	}
}

func getAllUserPackagesHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Mock - getAllUserPackagesHandler %s '%s'", r.Method, r.URL)
	if r.Method == http.MethodGet {
		var data *[]github_model.UserPackage
		if packageData != nil {
			data = &[]github_model.UserPackage{*packageData}
		} else {
			data = &[]github_model.UserPackage{}
		}

		GetAllUserPackagesCounter++
		w.Header().Set("Content-Type", gitHubModelJsonType)
		json.NewEncoder(w).Encode(data)
	} else {
		w.WriteHeader(500)
	}
}
