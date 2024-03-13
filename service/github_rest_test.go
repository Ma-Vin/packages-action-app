package service

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/ma-vin/packages-action/config"
	"github.com/ma-vin/packages-action/service/github_model"
	"github.com/ma-vin/packages-action/testutil"
)

const packageResponse = `{
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

var conf = config.Config{User: "DummyUser", PackageType: "maven", PackageName: "DummyPackage"}

type mockCloser struct {
	mockIoErrorReader
}

type mockIoErrorReader struct {
}

func (mockCloser) Close() error { return nil }
func (mockIoErrorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("IoTestError")
}

func createResponse(body *string, httpStatus int) *http.Response {
	var res http.Response

	res.Body = io.NopCloser(bytes.NewReader([]byte(*body)))
	res.StatusCode = httpStatus

	return &res
}

func createResponseWithIoError(httpStatus int) *http.Response {
	var res http.Response
	var reader mockIoErrorReader
	res.Body = mockCloser{reader}
	res.StatusCode = httpStatus

	return &res
}

func createDefaultPackageResponse() *http.Response {
	var body = packageResponse
	return createResponse(&body, 200)
}

func createDefaultPackagesArrayResponse() *http.Response {
	var body = fmt.Sprintf("[%s]", packageResponse)
	return createResponse(&body, 200)
}

func TestGetUserPackageSuccessful(t *testing.T) {
	SetClientExecutor(func(c *http.Client, req *http.Request) (*http.Response, error) {
		return createDefaultPackageResponse(), nil
	})

	userPackage, err := GetUserPackage(conf.PackageName, &conf)

	testutil.AssertNotNil(userPackage, t, "userPackage")
	testutil.AssertEquals(123456, userPackage.Id, t, "package id")
	testutil.AssertEquals(github_model.MAVEN, userPackage.PackageType, t, "package type")
	testutil.AssertNil(err, t, "err")
}

func TestGetUserPackageWithError(t *testing.T) {
	SetClientExecutor(func(c *http.Client, req *http.Request) (*http.Response, error) {
		return nil, errors.New("SomeTestError")
	})

	userPackage, err := GetUserPackage(conf.PackageName, &conf)

	testutil.AssertNil(userPackage, t, "userPackage")
	testutil.AssertNotNil(err, t, "err")
	testutil.AssertEquals("SomeTestError", err.Error(), t, "error message")
}

func TestGetUserPackageWithErrorHttpStatus(t *testing.T) {
	SetClientExecutor(func(c *http.Client, req *http.Request) (*http.Response, error) {
		res := createDefaultPackageResponse()
		res.StatusCode = 400
		return res, nil
	})

	userPackage, err := GetUserPackage(conf.PackageName, &conf)

	testutil.AssertNil(userPackage, t, "userPackage")
	testutil.AssertNotNil(err, t, "err")
	testutil.AssertEquals("an error status code occured: 400 - Bad Request", err.Error(), t, "error message")
}

func TestGetUserPackageWithBodyIoError(t *testing.T) {
	SetClientExecutor(func(c *http.Client, req *http.Request) (*http.Response, error) {
		res := createResponseWithIoError(200)
		return res, nil
	})

	userPackage, err := GetUserPackage(conf.PackageName, &conf)

	testutil.AssertNil(userPackage, t, "userPackage")
	testutil.AssertNotNil(err, t, "err")
	testutil.AssertEquals("IoTestError", err.Error(), t, "error message")
}

func TestGetUserPackageWithInvalidJsonError(t *testing.T) {
	SetClientExecutor(func(c *http.Client, req *http.Request) (*http.Response, error) {
		var body = ""
		res := createResponse(&body, 200)
		return res, nil
	})

	userPackage, err := GetUserPackages(&conf)

	testutil.AssertNil(userPackage, t, "userPackage")
	testutil.AssertNotNil(err, t, "err")
	testutil.AssertEquals("unexpected end of JSON input", err.Error(), t, "error message")
}

func TestGetUserPackagesArraySuccessful(t *testing.T) {
	SetClientExecutor(func(c *http.Client, req *http.Request) (*http.Response, error) {
		return createDefaultPackagesArrayResponse(), nil
	})

	userPackage, err := GetUserPackages(&conf)

	testutil.AssertNotNil(userPackage, t, "userPackage")

	testutil.AssertEquals(1, len(*userPackage), t, "package id")
	testutil.AssertEquals(123456, (*userPackage)[0].Id, t, "package id")
	testutil.AssertEquals(github_model.MAVEN, (*userPackage)[0].PackageType, t, "package type")
	testutil.AssertNil(err, t, "err")
}

func TestGetUserPackagesArrayWithError(t *testing.T) {
	SetClientExecutor(func(c *http.Client, req *http.Request) (*http.Response, error) {
		return nil, errors.New("SomeTestError")
	})

	userPackage, err := GetUserPackages(&conf)

	testutil.AssertNil(userPackage, t, "userPackage")
	testutil.AssertNotNil(err, t, "err")
	testutil.AssertEquals("SomeTestError", err.Error(), t, "error message")
}

func TestGetUserPackagesArrayWithErrorHttpStatus(t *testing.T) {
	SetClientExecutor(func(c *http.Client, req *http.Request) (*http.Response, error) {
		res := createDefaultPackageResponse()
		res.StatusCode = 400
		return res, nil
	})

	userPackage, err := GetUserPackages(&conf)

	testutil.AssertNil(userPackage, t, "userPackage")
	testutil.AssertNotNil(err, t, "err")
	testutil.AssertEquals("an error status code occured: 400 - Bad Request", err.Error(), t, "error message")
}

func TestGetUserPackagesArrayWithBodyIoError(t *testing.T) {
	SetClientExecutor(func(c *http.Client, req *http.Request) (*http.Response, error) {
		res := createResponseWithIoError(200)
		return res, nil
	})

	userPackage, err := GetUserPackages(&conf)

	testutil.AssertNil(userPackage, t, "userPackage")
	testutil.AssertNotNil(err, t, "err")
	testutil.AssertEquals("IoTestError", err.Error(), t, "error message")
}

func TestGetUserPackagesArrayWithInvalidJsonError(t *testing.T) {
	SetClientExecutor(func(c *http.Client, req *http.Request) (*http.Response, error) {
		var body = ""
		res := createResponse(&body, 200)
		return res, nil
	})

	userPackage, err := GetUserPackages(&conf)

	testutil.AssertNil(userPackage, t, "userPackage")
	testutil.AssertNotNil(err, t, "err")
	testutil.AssertEquals("unexpected end of JSON input", err.Error(), t, "error message")
}
