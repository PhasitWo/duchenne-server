package web_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/PhasitWo/duchenne-server/handlers/web"
	"github.com/PhasitWo/duchenne-server/model"
	"github.com/PhasitWo/duchenne-server/notification"
	"github.com/PhasitWo/duchenne-server/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

/* new method when using expecter struct
type-safe method replace .On
RunAndReturn to dynamically set a return value based on the input to the mock's call
*/

// $ gotest -v ./tests/web

func TestGetAllAppointment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("success", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().GetAllAppointment(9999, 0).Return([]model.SafeAppointment{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/", webH.GetAllAppointment)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 200, recorder.Code)
	})
	t.Run("invalidTypeQueryParam", func(t *testing.T) {
		// setup mock
		webH := web.WebHandler{}

		req := httptest.NewRequest(http.MethodGet, "/?type=brabra", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/", webH.GetAllAppointment)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 400, recorder.Code)
	})
	t.Run("invalidLimitQueryParam", func(t *testing.T) {
		// setup mock
		webH := web.WebHandler{}

		req := httptest.NewRequest(http.MethodGet, "/?limit=brabra", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/", webH.GetAllAppointment)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 400, recorder.Code)
	})
	t.Run("internalError", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().GetAllAppointment(9999, 0).Return([]model.SafeAppointment{}, errors.New("some internal error"))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/", webH.GetAllAppointment)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
}

func TestGetOneAppointment(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		id := "1"
		apm := model.SafeAppointment{
			Appointment: model.Appointment{
				ID: 1,
			},
		}
		// setup mock
		repo := repository.NewMockRepo(t)
		repo.EXPECT().GetAppointment(id).Return(apm, nil).Once()
		webH := web.WebHandler{Repo: repo}

		rr := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(rr)
		c.Params = gin.Params{gin.Param{Key: "id", Value: id}}

		webH.GetAppointment(c)

		respBody, err := json.Marshal(apm) // expected respBody
		assert.NoError(t, err)
		assert.Equal(t, 200, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
	})
	t.Run("notFound", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}
		
		mockErr := fmt.Errorf("wrap : %w", gorm.ErrRecordNotFound)
		repo.EXPECT().GetAppointment("1").Return(model.SafeAppointment{}, mockErr).Once()

		req := httptest.NewRequest(http.MethodGet, "/1", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/:id", webH.GetAppointment)
		router.ServeHTTP(recorder, req) // work around for ctx.Status() not setting status code when testing

		assert.Equal(t, 404, recorder.Code)
	})
	t.Run("internalError", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		mockErr := fmt.Errorf("wrap : %w", errors.New("some internal error"))
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().GetAppointment("1").Return(model.SafeAppointment{}, mockErr).Once()

		req := httptest.NewRequest(http.MethodGet, "/1", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/:id", webH.GetAppointment)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
}

func TestCreateAppointment(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		apmReq := model.CreateAppointmentRequest{
			Date:      int(time.Now().Add(60 * time.Minute).Unix()),
			PatientId: 1,
			DoctorId:  1,
			Approve:   false,
		}
		input, err := json.Marshal(&apmReq)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		noti := notification.NewMockService(t)
		webH := web.WebHandler{Repo: repo, NotiService: noti}

		repo.EXPECT().CreateAppointment(mock.Anything).Return(1, nil).Once()
		noti.EXPECT().SendNotiByPatientId(mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe() // go routine

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(input))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/", webH.CreateAppointment)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 201, recorder.Code)
	})
	t.Run("bindingError", func(t *testing.T) {
		apmReq := model.CreateAppointmentRequest{
			Date:      int(time.Now().Add(60 * time.Minute).Unix()),
			PatientId: 0, // error require
			DoctorId:  0, // error require
			Approve:   false,
		}
		input, err := json.Marshal(&apmReq)
		assert.NoError(t, err)
		// setup mock
		webH := web.WebHandler{}

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(input))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/", webH.CreateAppointment)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 400, recorder.Code)
	})
	t.Run("invalidDate", func(t *testing.T) {
		apmReq := model.CreateAppointmentRequest{
			Date:      1, // Date < now
			PatientId: 1,
			DoctorId:  1,
			Approve:   false,
		}
		input, err := json.Marshal(&apmReq)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(input))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/", webH.CreateAppointment)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 422, recorder.Code)
	})
	t.Run("internalError", func(t *testing.T) {
		apmReq := model.CreateAppointmentRequest{
			Date:      int(time.Now().Add(60 * time.Minute).Unix()),
			PatientId: 1,
			DoctorId:  1,
			Approve:   false,
		}
		input, err := json.Marshal(&apmReq)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().CreateAppointment(mock.Anything).Return(-1, errors.New("some internal err")).Once()

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(input))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/", webH.CreateAppointment)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
}

func TestUpdateAppointment(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		apmReq := model.CreateAppointmentRequest{
			Date:      int(time.Now().Add(60 * time.Minute).Unix()),
			PatientId: 1,
			DoctorId:  1,
			Approve:   false,
		}
		input, err := json.Marshal(&apmReq)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		noti := notification.NewMockService(t)
		webH := web.WebHandler{Repo: repo, NotiService: noti}

		repo.EXPECT().UpdateAppointment(mock.Anything).Return(nil).Once()
		noti.EXPECT().SendNotiByPatientId(mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe() // go routine

		req := httptest.NewRequest(http.MethodPut, "/1", bytes.NewReader(input))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT("/:id", webH.UpdateAppointment)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 200, recorder.Code)
	})
	t.Run("bindingError", func(t *testing.T) {
		apmReq := model.CreateAppointmentRequest{
			Date:      int(time.Now().Add(60 * time.Minute).Unix()),
			PatientId: 0, // error require
			DoctorId:  0, // error require
			Approve:   false,
		}
		input, err := json.Marshal(&apmReq)
		assert.NoError(t, err)
		// setup mock
		webH := web.WebHandler{}

		req := httptest.NewRequest(http.MethodPut, "/1", bytes.NewReader(input))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT("/:id", webH.UpdateAppointment)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 400, recorder.Code)
	})
	t.Run("invalidDate", func(t *testing.T) {
		apmReq := model.CreateAppointmentRequest{
			Date:      1, // Date < now
			PatientId: 1,
			DoctorId:  1,
			Approve:   false,
		}
		input, err := json.Marshal(&apmReq)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		req := httptest.NewRequest(http.MethodPut, "/1", bytes.NewReader(input))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT("/:id", webH.UpdateAppointment)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 422, recorder.Code)
	})
	t.Run("internalError", func(t *testing.T) {
		apmReq := model.CreateAppointmentRequest{
			Date:      int(time.Now().Add(60 * time.Minute).Unix()),
			PatientId: 1,
			DoctorId:  1,
			Approve:   false,
		}
		input, err := json.Marshal(&apmReq)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().UpdateAppointment(mock.Anything).Return(errors.New("some internal err")).Once()

		req := httptest.NewRequest(http.MethodPut, "/1", bytes.NewReader(input))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT("/:id", webH.UpdateAppointment)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
}

func TestDeleteAppointment(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		apm := model.SafeAppointment{
			Appointment: model.Appointment{
				ID: 1,
				PatientID: 2,
			},
		}
		// setup mock
		repo := repository.NewMockRepo(t)
		noti := notification.NewMockService(t)
		webH := web.WebHandler{Repo: repo, NotiService: noti}

		repo.EXPECT().GetAppointment("1").Return(apm, nil).Once()
		repo.EXPECT().DeleteAppointment("1").Return(nil).Once()
		noti.EXPECT().SendNotiByPatientId(2, mock.Anything, mock.Anything).Return(nil).Maybe() // go routine

		req := httptest.NewRequest(http.MethodDelete, "/1", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.DELETE("/:id", webH.DeleteAppointment)
		router.ServeHTTP(recorder, req) // work around for ctx.Status() not setting status code when testing

		assert.Equal(t, 204, recorder.Code)
	})
	t.Run("notFound", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}
		
		mockErr := fmt.Errorf("wrap : %w", gorm.ErrRecordNotFound)
		repo.EXPECT().GetAppointment("1").Return(model.SafeAppointment{}, mockErr).Once()

		req := httptest.NewRequest(http.MethodDelete, "/1", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.DELETE("/:id", webH.DeleteAppointment)
		router.ServeHTTP(recorder, req) // work around for ctx.Status() not setting status code when testing

		assert.Equal(t, 404, recorder.Code)
	})
	t.Run("getInternalError", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}
		
		mockErr := fmt.Errorf("wrap : %w", errors.New("some internal error"))
		repo.EXPECT().GetAppointment("1").Return(model.SafeAppointment{}, mockErr).Once()

		req := httptest.NewRequest(http.MethodDelete, "/1", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.DELETE("/:id", webH.DeleteAppointment)
		router.ServeHTTP(recorder, req) // work around for ctx.Status() not setting status code when testing

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("deleteInternalError", func(t *testing.T) {
		apm := model.SafeAppointment{
			Appointment: model.Appointment{
				ID: 1,
				PatientID: 2,
			},
		}
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().GetAppointment("1").Return(apm, nil).Once()
		repo.EXPECT().DeleteAppointment("1").Return(errors.New("some internal error")).Once()

		req := httptest.NewRequest(http.MethodDelete, "/1", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.DELETE("/:id", webH.DeleteAppointment)
		router.ServeHTTP(recorder, req) // work around for ctx.Status() not setting status code when testing

		assert.Equal(t, 500, recorder.Code)
	})
}
