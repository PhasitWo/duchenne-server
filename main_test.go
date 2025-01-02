package main

import (
	// "encoding/json"
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"testing"

	"github.com/PhasitWo/duchenne-server/config"
	"github.com/PhasitWo/duchenne-server/handlers/mobile"
	"github.com/PhasitWo/duchenne-server/middleware"
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

func TestMain(m *testing.M) {
	config.LoadConfig()
	db := setupDB()
	defer db.Close()
	// db.Exec("INSERT INTO appointment (id ,create_at, date, patient_id, doctor_id) VALUES (?,?, ?, ?, ?)", 9999, 1111, 2222, 3, 1) // FOR delete appointment
	tx, _ := db.Begin()
	defer tx.Rollback()
	setupDBdata(db, tx) // mockup data
	handler = &mobile.MobileHandler{Repo: repository.New(tx), DBConn: db}
	router = setupTestRouter()
	attachTestHandler(router, handler)
	m.Run()
}

func attachTestHandler(r *gin.Engine, m *mobile.MobileHandler) {
	mobile := r.Group("/mobile")
	{
		mobileAuth := mobile.Group("/auth")
		{
			mobileAuth.POST("/login", m.Login)
			mobileAuth.POST("/signup", m.Signup)
			mobileAuth.POST("/logout", middleware.MobileAuthMiddleware, m.Logout)
		}
		mobileProtected := mobile.Group("/api")
		mobileProtected.Use(middleware.MobileAuthMiddleware)
		{
			mobileProtected.GET("/test", m.Test)
			mobileProtected.GET("/profile", m.GetProfile)
			mobileProtected.GET("/appointment", m.GetAllPatientAppointment)
			mobileProtected.GET("/appointment/:id", m.GetAppointment)
			mobileProtected.POST("/appointment", m.CreateAppointment)
			mobileProtected.DELETE("/appointment/:id", m.DeleteAppointment)
			mobileProtected.GET("/question", m.GetAllPatientQuestion)
			mobileProtected.GET("/question/:id", m.GetQuestion)
			mobileProtected.POST("/question", m.CreateQuestion)
			mobileProtected.DELETE("/question/:id", m.DeleteQuestion)
			mobileProtected.GET("/doctor", m.GetAllDoctor)
			mobileProtected.GET("/device", m.GetAllDevice)
			mobileProtected.POST("/device", m.CreateDevice)
		}
	}
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	r := gin.Default()
	return r
}

var selfAppointmentId int
var otherPatientAppointmentId int
var toBeDeletedAppointmentId int
var selfQuestionId int
var otherPatientQuestionId int
var toBeDeletedQuestionId int

func setupDBdata(db *sql.DB ,tx *sql.Tx) {
	// verified account
	db.Exec(`insert into patient (id, hn, first_name, middle_name, last_name, email, phone, verified) 
	values (1234 , "mt1", "fnmt1", "mnmt1", "lnmt1", "mt@test.com", "9193929", 1)`)
	// unverified account
	db.Exec(`insert into patient (id, hn, first_name, middle_name, last_name, email, phone, verified) 
	values (4321 , "mt2", "fnmt2", NULL, "lnmt2", NULL, NULL, 0)`)
	db.Exec(`insert into patient (id, hn, first_name, middle_name, last_name, email, phone, verified) 
	values (5551 , "mt3", "fnmt3", "mnmt3", "lnmt3", NULL, NULL, 0)`)
	// appointment
	r := repository.New(tx)
	selfAppointmentId, _ = r.CreateAppointment(555, 555, 3, 1)
	toBeDeletedAppointmentId, _ = r.CreateAppointment(555, 555, 3, 1)
	otherPatientAppointmentId, _ = r.CreateAppointment(555, 555, 4, 1)
	// question
	selfQuestionId, _ = r.CreateQuestion(3, "to be deleted", "asdasdsa", 555)
	toBeDeletedQuestionId, _ = r.CreateQuestion(3, "to be deleted", "asdasdsa", 555)
	otherPatientQuestionId, _ = r.CreateQuestion(4, "haha xdxd", "asdasdsa", 555)
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

func TestAuthMiddleware(t *testing.T) {
	testCases := []testCase{
		{name: "request with no token", requestBody: nil, expected: http.StatusUnauthorized},
		{name: "request with invalid token", authToken: invalidAuthToken, requestBody: nil, expected: http.StatusUnauthorized},
		{name: "request with valid token", authToken: validAuthToken, requestBody: nil, expected: http.StatusOK}}
	testInternal(t, testCases, "GET", "/api/profile")
}

func TestLogin(t *testing.T) {
	validInput := []byte(`{"hn" : "mt1","firstName" : "fnmt1","lastName" : "lnmt1","deviceName": "main_test","expoToken": "dummy-expo-token"}`)
	badInput := []byte(`{"firstName" : "mt1","lastName" : "lnmt1","deviceName": "main_test","expoToken": "dummy-expo-token"}`)
	invalidInput := []byte(`{"hn" : "mt1","firstName" : "fn3123","lastName" : "lnmt1","deviceName": "main_test","expoToken": "dummy-expo-token"}`)
	unverifiedAccInput := []byte(`{"hn" : "mt2","firstName" : "fnmt2","lastName" : "lnmt2","deviceName": "main_test","expoToken": "dummy-expo-token"}`)
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
	validInput := []byte(`{"hn" : "mt3","firstName" : "fnmt3", "middleName" : "mnmt3","lastName" : "lnmt3","phone": "0000000","email": "test@tmail.com"}`)
	badInput := []byte(`{"firstName" : "fnmt3", "middleName" : "mnmt3","lastName" : "lnmt3","phone": "0000000","email": "test@tmail.com"}`)
	mnNotRequireMnInput := []byte(`{"hn" : "mt2","firstName" : "fnmt2", "middleName" : "mnmt2","lastName" : "lnmt2","phone": "0000000","email": "test@tmail.com"}`)
	noMnRequireMnInput := []byte(`{"hn" : "mt3","firstName" : "fnmt3","lastName" : "lnmt3","phone": "0000000","email": "test@tmail.com"}`)
	nonExistentAccInput := []byte(`{"hn" : "test28aa","firstName" : "fn28", "middleName" : "mn28","lastName" : "ln28","phone": "0000000","email": "test@tmail.com"}`)
	AlreadyVerifiedAccInput := []byte(`{"hn" : "mt1","firstName" : "fnmt1", "middleName" : "mnmt1","lastName" : "lnmt1","phone": "0000000","email": "test@tmail.com"}`)
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
	testInternal(t, testCaseSet1, "GET", fmt.Sprintf("/api/appointment/%d", selfAppointmentId))
	testCaseSet2 := []testCase{
		{name: "request to other patient's appointment", authToken: validAuthToken, requestBody: nil, expected: http.StatusUnauthorized},
	}
	testInternal(t, testCaseSet2, "GET", fmt.Sprintf("/api/appointment/%d", otherPatientAppointmentId))
	testCaseSet3 := []testCase{
		{name: "request to nonexistent appointment", authToken: validAuthToken, requestBody: nil, expected: http.StatusNotFound},
	}
	testInternal(t, testCaseSet3, "GET", "/api/appointment/999999")
}

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
		{name: "request with valid input", authToken: validAuthToken, requestBody: validInput, expected: http.StatusCreated},
	}
	testInternal(t, testCases, "POST", "/api/appointment")
}

// request to other patient's appointment
func TestDeleteAppointment(t *testing.T) {
	testInternal(
		t,
		[]testCase{{name: "request to own appointment", authToken: validAuthToken, expected: http.StatusNoContent}},
		"DELETE",
		fmt.Sprintf("/api/appointment/%d", toBeDeletedAppointmentId),
	)
	testInternal(
		t,
		[]testCase{{name: "request to other patient's appointment", authToken: validAuthToken, expected: http.StatusUnauthorized}},
		"DELETE",
		fmt.Sprintf("/api/appointment/%d", otherPatientAppointmentId),
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
	testInternal(t, testCaseSet1, "GET", fmt.Sprintf("/api/question/%d", selfQuestionId))
	testCaseSet2 := []testCase{
		{name: "request to other patient's question", authToken: validAuthToken, requestBody: nil, expected: http.StatusUnauthorized},
	}
	testInternal(t, testCaseSet2, "GET", fmt.Sprintf("/api/question/%d", otherPatientQuestionId))
	testCaseSet3 := []testCase{
		{name: "request to nonexistent question", authToken: validAuthToken, requestBody: nil, expected: http.StatusNotFound},
	}
	testInternal(t, testCaseSet3, "GET", "/api/question/99999")
}

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
		{name: "request with valid input", authToken: validAuthToken, requestBody: validInput, expected: http.StatusCreated},
	}
	testInternal(t, testCases, "POST", "/api/question")
}

func TestDeleteQuestion(t *testing.T) {
	testInternal(
		t,
		[]testCase{{name: "request to own question", authToken: validAuthToken, expected: http.StatusNoContent}},
		"DELETE",
		fmt.Sprintf("/api/question/%d", toBeDeletedQuestionId),
	)
	testInternal(
		t,
		[]testCase{{name: "request to other patient's question", authToken: validAuthToken, expected: http.StatusUnauthorized}},
		"DELETE",
		fmt.Sprintf("/api/question/%d", otherPatientQuestionId),
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

func TestGetAllDevice(t *testing.T) {
	testInternal(
		t,
		[]testCase{{name: "request to get active devices for push notifications", authToken: validAuthToken, expected: http.StatusOK}},
		"GET",
		"/api/device",
	)
}

func TestCreateDevice(t *testing.T) {
	req := []byte(`{"deviceName" : "thunder client","expoToken" : "test insert new device"}`)
	tc := []testCase{{name: "request to insert new device for push notifications", authToken: validAuthToken, requestBody: req, expected: http.StatusOK}}
	testInternal(
		t,
		tc,
		"POST",
		"/api/device",
	)
}
