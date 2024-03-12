package service

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/ma-vin/packages-action/config"
	"github.com/ma-vin/packages-action/service/github_model"
	"github.com/ma-vin/packages-action/testutil"
)

var conf = config.Config{User: "DummyUser", PackageType: "maven", PackageName: "DummyPackage"}

func createResponse(body *string, httpStatus int) *http.Response {
	var res http.Response

	res.Body = io.NopCloser(bytes.NewReader([]byte(*body)))
	res.StatusCode = httpStatus

	return &res
}

func createDefaultPackageResponse() *http.Response {
	var body = `{
		"id": 123456,
		"name": "DummyPackage",
		"package_type": "maven",
		"owner": {
		  "login": "DummyUser",
		  "id": 654321
		},
		"version_count": 12,
		"visibility": "public",
		"created_at": "2024-03-12T20:00:00Z",
		"updated_at": "2024-03-12T20:00:00Z",
		"repository": {
		  "id": 123654,
		  "node_id": "abcdef",
		  "name": "dummy-repo",
		  "full_name": "test/dummy-repo",
		  "private": false,
		  "owner": {
			"login": "DummyUser",
			"id": 654321,
			"node_id": "fedcba",
			"type": "User",
			"site_admin": true
		  }
		}
	  }`

	return createResponse(&body, 200)
}

func TestGetUserPackages(t *testing.T) {
	SetClientExecutor(func(c *http.Client, req *http.Request) (*http.Response, error) {
		return createDefaultPackageResponse(), nil
	})

	userPackage, err := GetUserPackage(conf.PackageName, &conf)

	testutil.AssertNotNil(userPackage, t, "userPackage")
	testutil.AssertEquals(123456, userPackage.Id, t, "package id")
	testutil.AssertEquals(github_model.MAVEN, userPackage.PackageType, t, "package type")
	testutil.AssertNil(err, t, "err")
}
