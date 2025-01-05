package main

import (
	// "encoding/json"
	"bytes"
	"database/sql"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/PhasitWo/duchenne-server/config"
	"github.com/PhasitWo/duchenne-server/handlers/mobile"
	"github.com/PhasitWo/duchenne-server/handlers/web"
	"github.com/PhasitWo/duchenne-server/model"
	"github.com/PhasitWo/duchenne-server/repository"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name        string
	authToken   string
	requestBody []byte
	expected    interface{}
}

// test3 fn3 ln3
var router *gin.Engine
var mobileValidAuthToken string
var mobileInvalidAuthToken string
var webValidAuthToken string
var webInvalidAuthToken string

func TestMain(m *testing.M) {
	config.LoadConfig()
	config.AppConfig.ENABLE_REDIS = false // disable redis
	db := setupDB()
	defer db.Close()
	// db.Exec("INSERT INTO appointment (id ,create_at, date, patient_id, doctor_id) VALUES (?,?, ?, ?, ?)", 9999, 1111, 2222, 3, 1) // FOR delete appointment
	tx, _ := db.Begin()
	defer tx.Rollback()
	setupDBdata(db, tx) // mockup data
	mobileHandler := &mobile.MobileHandler{Repo: repository.New(tx), DBConn: db}
	webHandler := &web.WebHandler{Repo: repository.New(tx), DBConn: db}
	router = setupTestRouter()
	attachHandler(router, mobileHandler, webHandler, nil)
	// token
	viper.SetConfigFile("test.yaml")
	viper.ReadInConfig()
	mobileValidAuthToken = viper.GetString("mobileTest.validAuthToken")
	mobileInvalidAuthToken = viper.GetString("mobileTest.invalidAuthToken")
	webValidAuthToken = viper.GetString("webTest.validAuthToken")
	webInvalidAuthToken = viper.GetString("webTest.invalidAuthToken")
	m.Run()
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
var existing1DoctorId int
var existing2DoctorId int
var toBeDeletedDoctorId int
var existing1PatientId int = 1234

func setupDBdata(db *sql.DB, tx *sql.Tx) {
	// verified patient account
	db.Exec(`insert into patient (id, hn, first_name, middle_name, last_name, email, phone, verified) 
	values (1234 , "mt1", "fnmt1", "mnmt1", "lnmt1", "mt@test.com", "9193929", 1)`)
	// unverified patient account
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
	// doctor
	// to test duplicate username
	existing1DoctorId, _ = r.CreateDoctor(model.Doctor{FirstName: "main_test", LastName: "main_test", Username: "main_test_duplicate", Password: "1234", Role: model.ADMIN})
	existing2DoctorId, _ = r.CreateDoctor(model.Doctor{FirstName: "main_test", LastName: "main_test", Username: "main_test", Password: "1234", Role: model.ADMIN})
	toBeDeletedDoctorId, _ = r.CreateDoctor(model.Doctor{FirstName: "main_test", LastName: "main_test", Username: "main_test_delete", Password: "1234", Role: model.ADMIN})
}

func testInternal(t *testing.T, testCases []testCase, method string, url string) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(method, url, bytes.NewBuffer(tc.requestBody))
			if tc.authToken != "" {
				req.Header.Set("Authorization", tc.authToken)
			}
			router.ServeHTTP(w, req)
			assert.Equal(t, tc.expected, w.Code)
		})
	}
}
