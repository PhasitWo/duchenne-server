package main

import (
	"net/http"
	"testing"
)

func TestWebAuthMiddleware(t *testing.T) {
	testCases := []testCase{
		{name: "request with no token", requestBody: nil, expected: http.StatusUnauthorized},
		{name: "request with invalid token", authToken: webInvalidAuthToken, requestBody: nil, expected: http.StatusUnauthorized},
		{name: "request with valid token", authToken: webValidAuthToken, requestBody: nil, expected: http.StatusOK}}
	testInternal(t, testCases, "GET", "/web/api/profile")
}

func TestWebGetProfile(t *testing.T) {
	testcases := []testCase{
		{name: "request", authToken: webValidAuthToken, expected: http.StatusOK},
	}
	testInternal(t, testcases, "GET", "/web/api/profile")
}

func TestWebUpdateProfile(t *testing.T) {
	badRequestBody := []byte(`
	{
	"firstNamesdsd" : "updatefn",
	"middleName" : "hahaxdxd",
    "lastName" : "updateln",
    "username" : "spdm",
    "password" : "1234"
	}
	`)
	noMiddleNameRequestBody := []byte(`
	{
	"firstName" : "updatefn",
    "lastName" : "updateln",
    "username" : "spdm",
    "password" : "1234"
	}
	`)
	hasMiddleNameRequestBody := []byte(`
	{
	"firstName" : "updatefn",
	"middleName" : "hahaxdxd",
    "lastName" : "updateln",
    "username" : "spdm",
    "password" : "1234"
	}
	`)
	testcases := []testCase{
		{name: "request with bad input", authToken: webValidAuthToken, requestBody: badRequestBody, expected: http.StatusBadRequest},
		{name: "request with no middle name", authToken: webValidAuthToken, requestBody: noMiddleNameRequestBody, expected: http.StatusOK},
		{name: "request with middle name", authToken: webValidAuthToken, requestBody: hasMiddleNameRequestBody, expected: http.StatusOK},
	}
	testInternal(t, testcases, "POST", "/web/api/profile")
}
