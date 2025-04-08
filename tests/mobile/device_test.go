package mobile_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PhasitWo/duchenne-server/handlers/mobile"
	"github.com/PhasitWo/duchenne-server/model"
	"github.com/PhasitWo/duchenne-server/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetAllDevice(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("noPatientIdFromAuthMiddleware", func(t *testing.T) {
		// setup mock
		mobileH := mobile.MobileHandler{}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/", mobileH.GetAllDevice)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("internalError", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		mobileH := mobile.MobileHandler{Repo: repo}

		repo.EXPECT().GetAllDevice(repository.Criteria{QueryCriteria: repository.PATIENTID, Value: 1}).Return([]model.Device{}, errors.New("err"))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/", func(ctx *gin.Context) { ctx.Set("patientId", 1) }, mobileH.GetAllDevice)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("success", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		mobileH := mobile.MobileHandler{Repo: repo}

		repo.EXPECT().GetAllDevice(repository.Criteria{QueryCriteria: repository.PATIENTID, Value: 1}).Return([]model.Device{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/", func(ctx *gin.Context) { ctx.Set("patientId", 1) }, mobileH.GetAllDevice)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 200, recorder.Code)
	})
}

// create device same as login