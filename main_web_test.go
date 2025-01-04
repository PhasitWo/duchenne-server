package main

import (
	"net/http"
	"strconv"
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
	dupEntryRequestBody := []byte(`
	{
	"firstName" : "updatefn",
	"middleName" : "hahaxdxd",
    "lastName" : "updateln",
    "username" : "main_test_duplicate",
    "password" : "1234"
	}
	`)
	testcases := []testCase{
		{name: "request with bad input", authToken: webValidAuthToken, requestBody: badRequestBody, expected: http.StatusBadRequest},
		{name: "request with duplicate username", authToken: webValidAuthToken, requestBody: dupEntryRequestBody, expected: http.StatusConflict},
		{name: "request with no middle name", authToken: webValidAuthToken, requestBody: noMiddleNameRequestBody, expected: http.StatusOK},
		{name: "request with middle name", authToken: webValidAuthToken, requestBody: hasMiddleNameRequestBody, expected: http.StatusOK},
	}
	testInternal(t, testcases, "PUT", "/web/api/profile")
}

func TestWebGetAllDoctor(t *testing.T) {
	testInternal(
		t,
		[]testCase{{name: "request", authToken: webValidAuthToken, expected: http.StatusOK}},
		"GET",
		"/web/api/doctor",
	)
}

func TestWebGetOneDoctor(t *testing.T) {
	testInternal(
		t,
		[]testCase{
			{name: "request", authToken: webValidAuthToken, expected: http.StatusOK},
		},
		"GET",
		"/web/api/doctor/"+strconv.Itoa(existing2DoctorId),
	)
	testInternal(
		t,
		[]testCase{
			{name: "request to nonexist doctor", authToken: webValidAuthToken, expected: http.StatusNotFound},
		},
		"GET",
		"/web/api/doctor/9999",
	)
}

func TestWebCreateDoctor(t *testing.T) {
	badRequest := []byte(`{
	"firstNamess" : "testfn",
	"middleName" : "testmn",
	"lastName": "testln",
	"username" : "xdxd",
  	"password" : "1234",
  	"role": "admin"}`)
	dupEntryRequest := []byte(`{
	"firstName" : "testfn",
	"middleName" : "testmn",
	"lastName": "testln",
	"username" : "main_test_duplicate",
  	"password" : "1234",
  	"role": "admin"}`)
	validRequest := []byte(`{
	"firstName" : "testfn",
	"middleName" : "testmn",
	"lastName": "testln",
	"username" : "xdxd",
  	"password" : "1234",
  	"role": "admin"}`)
	validRequestNoMn := []byte(`{
	"firstName" : "testfn",
	"lastName": "testln",
	"username" : "asdsadadsad",
  	"password" : "1234",
  	"role": "admin"}`)
	testInternal(
		t,
		[]testCase{
			{name: "request with bad input", authToken: webValidAuthToken, requestBody: badRequest, expected: http.StatusBadRequest},
			{name: "request with duplicate username", authToken: webValidAuthToken, requestBody: dupEntryRequest, expected: http.StatusConflict},
			{name: "request with valid input", authToken: webValidAuthToken, requestBody: validRequest, expected: http.StatusCreated},
			{name: "request with valid input but no middlename", authToken: webValidAuthToken, requestBody: validRequestNoMn, expected: http.StatusCreated},
		},
		"POST",
		"/web/api/doctor",
	)
}

func TestWebUpdateDoctor(t *testing.T) {
	badRequest := []byte(`{
	"firstNamess" : "testfn",
	"middleName" : "testmn",
	"lastName": "testln",
	"username" : "xdxd",
  	"password" : "1234",
  	"role": "admin"}`)
	dupEntryRequest := []byte(`{
	"firstName" : "testfn",
	"middleName" : "testmn",
	"lastName": "testln",
	"username" : "main_test_duplicate",
  	"password" : "1234",
  	"role": "admin"}`)
	validRequest := []byte(`{
	"firstName" : "testfn",
	"middleName" : "testmn",
	"lastName": "testln",
	"username" : "asdlkj999913123jlkj",
  	"password" : "1234",
  	"role": "admin"}`)
	validRequestNoMn := []byte(`{
	"firstName" : "testfn",
	"lastName": "testln",
	"username" : "asdlkj999913123jlkj",
  	"password" : "1234",
  	"role": "admin"}`)
	testInternal(
		t,
		[]testCase{
			{name: "request with bad url", authToken: webValidAuthToken, requestBody: dupEntryRequest, expected: http.StatusBadRequest},
		},
		"PUT",
		"/web/api/doctor/asd",
	)
	testInternal(
		t,
		[]testCase{
			{name: "request with bad input", authToken: webValidAuthToken, requestBody: badRequest, expected: http.StatusBadRequest},
			{name: "request with duplicate username", authToken: webValidAuthToken, requestBody: dupEntryRequest, expected: http.StatusConflict},
		},
		"PUT",
		"/web/api/doctor/"+strconv.Itoa(existing2DoctorId),
	)
	testInternal(
		t,
		[]testCase{
			{name: "request to nonexist doctor", authToken: webValidAuthToken, requestBody: validRequest, expected: http.StatusNotFound},
		},
		"PUT",
		"/web/api/doctor/99999",
	)
	testInternal(
		t,
		[]testCase{
			{name: "request with valid input", authToken: webValidAuthToken, requestBody: validRequest, expected: http.StatusOK},
			{name: "request with valid input no middlename", authToken: webValidAuthToken, requestBody: validRequestNoMn, expected: http.StatusOK},
		},
		"PUT",
		"/web/api/doctor/"+strconv.Itoa(existing1DoctorId),
	)
}
func TestWebDeleteDoctor(t *testing.T) {
	testInternal(
		t,
		[]testCase{
			{name: "request with bad url", authToken: webValidAuthToken, expected: http.StatusBadRequest},
		},
		"DELETE",
		"/web/api/doctor/asd",
	)
	testInternal(
		t,
		[]testCase{
			{name: "request to nonexist doctor", authToken: webValidAuthToken, expected: http.StatusNoContent},
		},
		"DELETE",
		"/web/api/doctor/9999",
	)
	testInternal(
		t,
		[]testCase{
			{name: "request to existing doctor", authToken: webValidAuthToken, expected: http.StatusNoContent},
		},
		"DELETE",
		"/web/api/doctor/"+strconv.Itoa(toBeDeletedDoctorId),
	)
}

func TestWebGetAllPatient(t *testing.T) {
	testInternal(
		t,
		[]testCase{{name: "request", authToken: webValidAuthToken, expected: http.StatusOK}},
		"GET",
		"/web/api/patient",
	)
}

func TestWebGetOnePatient(t *testing.T) {
	testInternal(
		t,
		[]testCase{
			{name: "request", authToken: webValidAuthToken, expected: http.StatusOK},
		},
		"GET",
		"/web/api/patient/"+strconv.Itoa(existing1PatientId),
	)
	testInternal(
		t,
		[]testCase{
			{name: "request to nonexist doctor", authToken: webValidAuthToken, expected: http.StatusNotFound},
		},
		"GET",
		"/web/api/patient/9999",
	)
}

func TestWebCreatePatient(t *testing.T) {
	badRequest := []byte(`{
	"hnss" : "asd",
	"firstName" : "testfn",
	"middleName" : "testmn",
	"lastName": "testln",
	"email" : "xdxd@tmail.com",
  	"phone" : "0900001122",
  	"verified": true}`)
	dupEntryRequest := []byte(`{
	"hn" : "mt1",
	"firstName" : "testfn",
	"middleName" : "testmn",
	"lastName": "testln",
	"email" : "xdxd@tmail.com",
  	"phone" : "0900001122",
  	"verified": true}`)
	validRequest := []byte(`{
	"hn" : "webtest1",
	"firstName" : "testfn",
	"middleName" : "testmn",
	"lastName": "testln",
	"email" : "xdxd@tmail.com",
  	"phone" : "0900001122",
  	"verified": true}`)
	validRequestNoVerified := []byte(`{
	"hn" : "webtest2",
	"firstName" : "testfn",
	"middleName" : "testmn",
	"lastName": "testln",
	"email" : "xdxd@tmail.com",
  	"phone" : "0900001122"
  	}`)
	validRequestNoMn := []byte(`{
	"hn" : "webtest3",
	"firstName" : "testfn",
	"lastName": "testln",
	"email" : "xdxd@tmail.com",
  	"phone" : "0900001122"
  	}`)
	testInternal(
		t,
		[]testCase{
			{name: "request with bad input", authToken: webValidAuthToken, requestBody: badRequest, expected: http.StatusBadRequest},
			{name: "request with duplicate username", authToken: webValidAuthToken, requestBody: dupEntryRequest, expected: http.StatusConflict},
			{name: "request with valid input", authToken: webValidAuthToken, requestBody: validRequest, expected: http.StatusCreated},
			{name: "request with valid input but not verified", authToken: webValidAuthToken, requestBody: validRequestNoVerified, expected: http.StatusCreated},
			{name: "request with valid input but no middlename", authToken: webValidAuthToken, requestBody: validRequestNoMn, expected: http.StatusCreated},
		},
		"POST",
		"/web/api/patient",
	)
}
