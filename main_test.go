package main

import (
	// "encoding/json"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"testing"

	"github.com/PhasitWo/duchenne-server/endpoint/mobile"
	"github.com/PhasitWo/duchenne-server/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name        string
	authToken   string
	requestBody []byte
	expected    interface{}
}

var handler *mobile.MobileHandler
var router *gin.Engine


const basePath = "/mobile"

// test3 fn3 ln3
const validAuthToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXRpZW50SWQiOjMsImRldmljZUlkIjoxMywiZXhwIjoxNzQxNDYxOTk3fQ.fo_5VWqPmyfRL0TccGdYMTh8enuFUu61-4G_aqoDnJY"
const invalidAuthToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXRpZW50SWQiOjMsImRldmljZUlkIjoxNCwiZXhwIjoxNzMzNzQ0OTM3fQ.WQnEQwt8AJuCOF22DeKCCh3hpmM-9hX8PgNZOZG8NFs"

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	r := gin.Default()
	return r
}

func TestMain(m *testing.M) {
	db := setupDB()
	defer db.Close()
	db.Exec("INSERT INTO appointment (id ,create_at, date, patient_id, doctor_id) VALUES (?,?, ?, ?, ?)", 9999, 1111, 2222, 3, 1) // FOR delete appointment
	tx, _ := db.Begin()
	defer tx.Rollback()
	handler = &mobile.MobileHandler{Repo: repository.New(tx), DBConn: db}
	router = setupTestRouter()
	attachHandler(router, handler)
	m.Run()
}

func testInternal(t *testing.T, testCases []testCase, method string, url string) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(method, basePath+url, bytes.NewBuffer(tc.requestBody))
			if tc.authToken != "" {
				req.Header.Set("Authorization", tc.authToken)
			}
			router.ServeHTTP(w, req)
			assert.Equal(t, tc.expected, w.Code)
		})
	}
}

func testInternalNoTx(t *testing.T, testCases []testCase, method string, url string) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(method, basePath+url, bytes.NewBuffer(tc.requestBody))
			if tc.authToken != "" {
				req.Header.Set("Authorization", tc.authToken)
			}
			router.ServeHTTP(w, req)
			assert.Equal(t, tc.expected, w.Code)
		})
	}
}

func TestAuthMiddleware(t *testing.T) {
	testCases := []testCase{
		{name: "request with no token", requestBody: nil, expected: http.StatusUnauthorized},
		{name: "request with invalid token", authToken: invalidAuthToken, requestBody: nil, expected: http.StatusUnauthorized},
		{name: "request with valid token", authToken: validAuthToken, requestBody: nil, expected: http.StatusOK}}
	testInternal(t, testCases, "GET", "/api/profile")
}

func TestLogin(t *testing.T) {
	validInput := []byte(`{"hn" : "test3","firstName" : "fn3","lastName" : "ln3","deviceName": "main_test","expoToken": "dummy-expo-token"}`)
	badInput := []byte(`{"firstName" : "fn3","lastName" : "ln3","deviceName": "main_test","expoToken": "dummy-expo-token"}`)
	invalidInput := []byte(`{"hn" : "test3","firstName" : "fn3123","lastName" : "ln3","deviceName": "main_test","expoToken": "dummy-expo-token"}`)
	unverifiedAccInput := []byte(`{"hn" : "test30","firstName" : "fn30","lastName" : "ln30","deviceName": "main_test","expoToken": "dummy-expo-token"}`)
	nonExistentAccInput := []byte(`{"hn" : "test30111","firstName" : "fn30","lastName" : "ln30","deviceName": "main_test","expoToken": "dummy-expo-token"}`)
	testCases := []testCase{
		{name: "request with valid input", requestBody: validInput, expected: http.StatusOK},
		{name: "request with bad input", requestBody: badInput, expected: http.StatusBadRequest},
		{name: "request with invalid input", requestBody: invalidInput, expected: http.StatusUnauthorized},
		{name: "request with unverified account input", requestBody: unverifiedAccInput, expected: http.StatusForbidden},
		{name: "request with nonexistent account", requestBody: nonExistentAccInput, expected: http.StatusNotFound}}
	testInternal(t, testCases, "POST", "/auth/login")
}

func TestSignup(t *testing.T) {
	validInput := []byte(`{"hn" : "test28","firstName" : "fn28", "middleName" : "mn28","lastName" : "ln28","phone": "0000000","email": "test@tmail.com"}`)
	badInput := []byte(`{"firstName" : "fn28", "middleName" : "mn28","lastName" : "ln28","phone": "0000000","email": "test@tmail.com"}`)
	mnNotRequireMnInput := []byte(`{"hn" : "test30","firstName" : "fn30", "middleName" : "something","lastName" : "ln30","phone": "0000000","email": "test@tmail.com"}`)
	noMnRequireMnInput := []byte(`{"hn" : "test28","firstName" : "fn28", "lastName" : "ln28","phone": "0000000","email": "test@tmail.com"}`)
	nonExistentAccInput := []byte(`{"hn" : "test28aa","firstName" : "fn28", "middleName" : "mn28","lastName" : "ln28","phone": "0000000","email": "test@tmail.com"}`)
	AlreadyVerifiedAccInput := []byte(`{"hn" : "test3","firstName" : "fn3", "middleName" : "","lastName" : "ln3","phone": "0000000","email": "test@tmail.com"}`)
	testCases := []testCase{
		{name: "request with bad input", requestBody: badInput, expected: http.StatusBadRequest},
		{name: "request with middleName but account not require middlename", requestBody: mnNotRequireMnInput, expected: http.StatusUnauthorized},
		{name: "request with no middleName but account require middlename", requestBody: noMnRequireMnInput, expected: http.StatusUnauthorized},
		{name: "request with nonexistent account", requestBody: nonExistentAccInput, expected: http.StatusNotFound},
		{name: "request with already verified account input", requestBody: AlreadyVerifiedAccInput, expected: http.StatusConflict},
		{name: "request with valid input", requestBody: validInput, expected: http.StatusOK}}
	testInternal(t, testCases, "POST", "/auth/signup")
}

func TestGetProfile(t *testing.T) {
	testCases := []testCase{
		{name: "request", authToken: validAuthToken, requestBody: nil, expected: http.StatusOK},
	}
	testInternal(t, testCases, "GET", "/api/profile")
}

func TestGetAllAppointment(t *testing.T) {
	testCases := []testCase{
		{name: "request", authToken: validAuthToken, requestBody: nil, expected: http.StatusOK},
	}
	testInternal(t, testCases, "GET", "/api/appointment")
}

func TestGetOneAppointment(t *testing.T) {
	testCaseSet1 := []testCase{
		{name: "request to own appointment", authToken: validAuthToken, requestBody: nil, expected: http.StatusOK},
	}
	testInternal(t, testCaseSet1, "GET", "/api/appointment/31")
	testCaseSet2 := []testCase{
		{name: "request to other patient's appointment", authToken: validAuthToken, requestBody: nil, expected: http.StatusUnauthorized},
	}
	testInternal(t, testCaseSet2, "GET", "/api/appointment/35")
	testCaseSet3 := []testCase{
		{name: "request to nonexistent appointment", authToken: validAuthToken, requestBody: nil, expected: http.StatusNotFound},
	}
	testInternal(t, testCaseSet3, "GET", "/api/appointment/100")
}

var insertedAppointmentId int

func TestCreateAppointment(t *testing.T) {
	validInput := []byte(`{ "date" : 1766120265, "doctorId" : 1}`)
	badDateInput := []byte(`{ "date" : 1733745465, "doctorId" : 1}`)
	badDoctorIdInput := []byte(`{ "date" : 1766120265, "doctorId" : 999}`)
	missingDoctorIdInput := []byte(`{ "date" : 1766120265`)
	missingDateInput := []byte(`{ "doctorId" : 1}`)

	testCases := []testCase{
		{name: "request with bad date input", authToken: validAuthToken, requestBody: badDateInput, expected: http.StatusUnprocessableEntity},
		{name: "request with bad doctorId input", authToken: validAuthToken, requestBody: badDoctorIdInput, expected: http.StatusInternalServerError},
		{name: "request with missing date input", authToken: validAuthToken, requestBody: missingDateInput, expected: http.StatusBadRequest},
		{name: "request with missing doctorId input", authToken: validAuthToken, requestBody: missingDoctorIdInput, expected: http.StatusBadRequest},
	}
	testInternal(t, testCases, "POST", "/api/appointment")

	t.Run("request with valid input", func(t *testing.T) {
		var resp map[string]interface{}
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", basePath+"/api/appointment", bytes.NewBuffer(validInput))
		req.Header.Set("Authorization", validAuthToken)
		router.ServeHTTP(w, req)
		if assert.Equal(t, http.StatusCreated, w.Code) {
			json.Unmarshal([]byte(w.Body.String()), &resp)
			insertedAppointmentId = int(resp["id"].(float64))
		}
	})
}

// request to other patient's appointment
func TestDeleteAppointment(t *testing.T) {
	testInternal(
		t,
		[]testCase{{name: "request to own appointment", authToken: validAuthToken, expected: http.StatusNoContent}},
		"DELETE",
		fmt.Sprintf("/api/appointment/%d", 9999),
	)
	testInternal(
		t,
		[]testCase{{name: "request to other patient's appointment", authToken: validAuthToken, expected: http.StatusUnauthorized}},
		"DELETE",
		"/api/appointment/35",
	)
	testInternal(
		t,
		[]testCase{{name: "request to nonexistent appointment", authToken: validAuthToken, expected: http.StatusNotFound}},
		"DELETE",
		"/api/appointment/8888",
	)
}

func TestGetAllQuestion(t *testing.T) {
	testCases := []testCase{
		{name: "request", authToken: validAuthToken, requestBody: nil, expected: http.StatusOK},
	}
	testInternal(t, testCases, "GET", "/api/question")
}

func TestGetOneQuestion(t *testing.T) {
	testCaseSet1 := []testCase{
		{name: "request to own question", authToken: validAuthToken, requestBody: nil, expected: http.StatusOK},
	}
	testInternal(t, testCaseSet1, "GET", "/api/question/25")
	testCaseSet2 := []testCase{
		{name: "request to other patient's question", authToken: validAuthToken, requestBody: nil, expected: http.StatusUnauthorized},
	}
	testInternal(t, testCaseSet2, "GET", "/api/question/30")
	testCaseSet3 := []testCase{
		{name: "request to nonexistent question", authToken: validAuthToken, requestBody: nil, expected: http.StatusNotFound},
	}
	testInternal(t, testCaseSet3, "GET", "/api/question/100")
}

var insertedQuestionId int

func TestCreateQuestion(t *testing.T) {
	validInput := []byte(`{ "topic" : "my topic", "question" : "my question"}`)
	exceedLimitTopicInput := []byte(fmt.Sprintf(`{ "topic" : "%s", "question" : "my question"}`, strings.Repeat("tests", 11)))
	exceedLimitQuestionInput := []byte(fmt.Sprintf(`{ "topic" : "sda", "question" : "%s"}`, strings.Repeat("tests", 141)))
	emptyTopicInput := []byte(`{ "topic" : "", "question" : "my question"}`)
	emptyQuestionInput := []byte(`{ "topic" : "my topic", "question" : ""}`)
	missingTopicInput := []byte(`{  "question" : "my question"}`)
	missingQuestionInput := []byte(`{ "topic" : "my topic"}`)

	testCases := []testCase{
		{name: "request with exceeding limit topic input", authToken: validAuthToken, requestBody: exceedLimitTopicInput, expected: http.StatusUnprocessableEntity},
		{name: "request with exceeding limit question input", authToken: validAuthToken, requestBody: exceedLimitQuestionInput, expected: http.StatusUnprocessableEntity},
		{name: "request with empty topic input", authToken: validAuthToken, requestBody: emptyTopicInput, expected: http.StatusBadRequest},
		{name: "request with empty question input", authToken: validAuthToken, requestBody: emptyQuestionInput, expected: http.StatusBadRequest},
		{name: "request with missing topic input", authToken: validAuthToken, requestBody: missingTopicInput, expected: http.StatusBadRequest},
		{name: "request with missing question input", authToken: validAuthToken, requestBody: missingQuestionInput, expected: http.StatusBadRequest},
	}
	testInternal(t, testCases, "POST", "/api/question")
	t.Run("request with valid input", func(t *testing.T) {
		var resp map[string]interface{}
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", basePath+"/api/question", bytes.NewBuffer(validInput))
		req.Header.Set("Authorization", validAuthToken)
		router.ServeHTTP(w, req)
		if assert.Equal(t, http.StatusCreated, w.Code) {
			json.Unmarshal([]byte(w.Body.String()), &resp)
			insertedQuestionId = int(resp["id"].(float64))
		}
	})
}

func TestDeleteQuestion(t *testing.T) {
	testInternal(
		t,
		[]testCase{{name: "request to own question", authToken: validAuthToken, expected: http.StatusNoContent}},
		"DELETE",
		fmt.Sprintf("/api/question/%d", insertedQuestionId),
	)
	testInternal(
		t,
		[]testCase{{name: "request to other patient's question", authToken: validAuthToken, expected: http.StatusUnauthorized}},
		"DELETE",
		"/api/question/30",
	)
	testInternal(
		t,
		[]testCase{{name: "request to nonexistent question", authToken: validAuthToken, expected: http.StatusNotFound}},
		"DELETE",
		"/api/question/9999",
	)
}

func TestGetAllDoctor(t *testing.T) {
	testInternal(
		t,
		[]testCase{{name: "request", authToken: validAuthToken, expected: http.StatusOK}},
		"GET",
		"/api/doctor",
	)
}
