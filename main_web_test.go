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