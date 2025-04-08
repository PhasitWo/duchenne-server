package mobile_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/PhasitWo/duchenne-server/handlers/mobile"
	"github.com/PhasitWo/duchenne-server/model"
	"github.com/PhasitWo/duchenne-server/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestGetAllPatientAppointment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("noPatientIdFromAuthMiddleware", func(t *testing.T) {
		// setup mock
		mobileH := mobile.MobileHandler{}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/", mobileH.GetAllPatientAppointment)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("internalError", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		mobileH := mobile.MobileHandler{Repo: repo}

		repo.EXPECT().GetAllAppointment(15, 0, mock.Anything).Return([]model.SafeAppointment{}, errors.New("err"))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/", func(ctx *gin.Context) { ctx.Set("patientId", 1) }, mobileH.GetAllPatientAppointment)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("success", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		mobileH := mobile.MobileHandler{Repo: repo}

		repo.EXPECT().GetAllAppointment(15, 0, mock.Anything).Return([]model.SafeAppointment{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/", func(ctx *gin.Context) { ctx.Set("patientId", 1) }, mobileH.GetAllPatientAppointment)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 200, recorder.Code)
	})
}

func TestGetOneAppointment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("noPatientIdFromAuthMiddleware", func(t *testing.T) {
		// setup mock
		mobileH := mobile.MobileHandler{}

		req := httptest.NewRequest(http.MethodGet, "/1", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/:id", mobileH.GetAppointment)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("notFound", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		mobileH := mobile.MobileHandler{Repo: repo}

		mockErr := fmt.Errorf("wrap : %w", gorm.ErrRecordNotFound)
		repo.EXPECT().GetAppointment("1").Return(model.SafeAppointment{}, mockErr)

		req := httptest.NewRequest(http.MethodGet, "/1", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/:id", func(ctx *gin.Context) { ctx.Set("patientId", 1) }, mobileH.GetAppointment)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 404, recorder.Code)
	})
	t.Run("internalError", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		mobileH := mobile.MobileHandler{Repo: repo}

		mockErr := fmt.Errorf("wrap : %w", errors.New("err"))
		repo.EXPECT().GetAppointment("1").Return(model.SafeAppointment{}, mockErr)

		req := httptest.NewRequest(http.MethodGet, "/1", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/:id", func(ctx *gin.Context) { ctx.Set("patientId", 1) }, mobileH.GetAppointment)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("unauthorized", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		mobileH := mobile.MobileHandler{Repo: repo}

		repo.EXPECT().GetAppointment("1").Return(
			model.SafeAppointment{
				Appointment: model.Appointment{Patient: model.Patient{ID: 2}},
			},
			nil,
		)

		req := httptest.NewRequest(http.MethodGet, "/1", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/:id", func(ctx *gin.Context) { ctx.Set("patientId", 1) }, mobileH.GetAppointment)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 401, recorder.Code)
	})
	t.Run("success", func(t *testing.T) {
		apm := model.SafeAppointment{
			Appointment: model.Appointment{Patient: model.Patient{ID: 1}},
		}
		// setup mock
		repo := repository.NewMockRepo(t)
		mobileH := mobile.MobileHandler{Repo: repo}

		repo.EXPECT().GetAppointment("1").Return(apm, nil)

		req := httptest.NewRequest(http.MethodGet, "/1", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/:id", func(ctx *gin.Context) { ctx.Set("patientId", 1) }, mobileH.GetAppointment)
		router.ServeHTTP(recorder, req)

		expectRespBody, err := json.Marshal(&apm)
		assert.NoError(t, err)

		assert.Equal(t, 200, recorder.Code)
		assert.Equal(t, expectRespBody, recorder.Body.Bytes())
	})
}
func TestCreateAppointment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("noPatientIdFromAuthMiddleware", func(t *testing.T) {
		// setup mock
		mobileH := mobile.MobileHandler{}

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/", mobileH.CreateAppointment)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("bindingError", func(t *testing.T) {
		input := model.PatientCreateAppointmentRequest{
			Date:     0, // error require
			DoctorId: 10,
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		mobileH := mobile.MobileHandler{}

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/", func(ctx *gin.Context) { ctx.Set("patientId", 1) }, mobileH.CreateAppointment)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 400, recorder.Code)
	})
	t.Run("invalidDate", func(t *testing.T) {
		input := model.PatientCreateAppointmentRequest{
			Date:     123, // date < now
			DoctorId: 10,
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		mobileH := mobile.MobileHandler{}

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/", func(ctx *gin.Context) { ctx.Set("patientId", 1) }, mobileH.CreateAppointment)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 422, recorder.Code)
	})
	t.Run("internalError", func(t *testing.T) {
		input := model.PatientCreateAppointmentRequest{
			Date:     int(time.Now().Add(2 * time.Hour).Unix()),
			DoctorId: 10,
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		mobileH := mobile.MobileHandler{Repo: repo}

		repo.EXPECT().CreateAppointment(mock.Anything).Return(-1, errors.New("err"))

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/", func(ctx *gin.Context) { ctx.Set("patientId", 1) }, mobileH.CreateAppointment)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("success", func(t *testing.T) {
		input := model.PatientCreateAppointmentRequest{
			Date:     int(time.Now().Add(2 * time.Hour).Unix()),
			DoctorId: 10,
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		mobileH := mobile.MobileHandler{Repo: repo}

		repo.EXPECT().CreateAppointment(
			model.Appointment{
				Date:      input.Date,
				PatientID: 1,
				DoctorID:  input.DoctorId,
				ApproveAt: nil,
			},
		).Return(22, nil)

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/", func(ctx *gin.Context) { ctx.Set("patientId", 1) }, mobileH.CreateAppointment)
		router.ServeHTTP(recorder, req)

		expectRespBody, err := json.Marshal(gin.H{"id": 22})
		assert.NoError(t, err)

		assert.Equal(t, 201, recorder.Code)
		assert.Equal(t, expectRespBody, recorder.Body.Bytes())
	})
}
func TestDeleteAppointment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("noPatientIdFromAuthMiddleware", func(t *testing.T) {
		// setup mock
		mobileH := mobile.MobileHandler{}

		req := httptest.NewRequest(http.MethodDelete, "/1", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.DELETE("/:id", mobileH.DeleteAppointment)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("notFound", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		mobileH := mobile.MobileHandler{Repo: repo}

		mockErr := fmt.Errorf("wrap : %w", gorm.ErrRecordNotFound)
		repo.EXPECT().GetAppointment("10").Return(model.SafeAppointment{}, mockErr)

		req := httptest.NewRequest(http.MethodDelete, "/10", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.DELETE("/:id", func(ctx *gin.Context) { ctx.Set("patientId", 1) }, mobileH.DeleteAppointment)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 404, recorder.Code)
	})
	t.Run("getAppointmentInternalError", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		mobileH := mobile.MobileHandler{Repo: repo}

		mockErr := fmt.Errorf("wrap : %w", errors.New("err"))
		repo.EXPECT().GetAppointment("10").Return(model.SafeAppointment{}, mockErr)

		req := httptest.NewRequest(http.MethodDelete, "/10", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.DELETE("/:id", func(ctx *gin.Context) { ctx.Set("patientId", 1) }, mobileH.DeleteAppointment)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("unauthorized", func(t *testing.T) {
		apm := model.SafeAppointment{Appointment: model.Appointment{Patient: model.Patient{ID: 22}}}
		// setup mock
		repo := repository.NewMockRepo(t)
		mobileH := mobile.MobileHandler{Repo: repo}

		repo.EXPECT().GetAppointment("10").Return(apm, nil)

		req := httptest.NewRequest(http.MethodDelete, "/10", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.DELETE("/:id", func(ctx *gin.Context) { ctx.Set("patientId", 1) }, mobileH.DeleteAppointment)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 401, recorder.Code)
	})
	t.Run("internalError", func(t *testing.T) {
		apm := model.SafeAppointment{Appointment: model.Appointment{Patient: model.Patient{ID: 1}}}
		// setup mock
		repo := repository.NewMockRepo(t)
		mobileH := mobile.MobileHandler{Repo: repo}

		repo.EXPECT().GetAppointment("10").Return(apm, nil)
		repo.EXPECT().DeleteAppointment("10").Return(errors.New("err"))

		req := httptest.NewRequest(http.MethodDelete, "/10", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.DELETE("/:id", func(ctx *gin.Context) { ctx.Set("patientId", 1) }, mobileH.DeleteAppointment)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("success", func(t *testing.T) {
		apm := model.SafeAppointment{Appointment: model.Appointment{Patient: model.Patient{ID: 1}}}
		// setup mock
		repo := repository.NewMockRepo(t)
		mobileH := mobile.MobileHandler{Repo: repo}

		repo.EXPECT().GetAppointment("10").Return(apm, nil)
		repo.EXPECT().DeleteAppointment("10").Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/10", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.DELETE("/:id", func(ctx *gin.Context) { ctx.Set("patientId", 1) }, mobileH.DeleteAppointment)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 204, recorder.Code)
	})
}
