package mobile_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PhasitWo/duchenne-server/handlers/mobile"
	"github.com/PhasitWo/duchenne-server/model"
	"github.com/PhasitWo/duchenne-server/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("noPatientIdFromAuthMiddleware", func(t *testing.T) {
		// setup mock
		mobileH := mobile.MobileHandler{}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/", mobileH.GetProfile)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("notFound", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		mobileH := mobile.MobileHandler{Repo: repo}

		mockErr := fmt.Errorf("wrap : %w", gorm.ErrRecordNotFound)
		repo.EXPECT().GetPatientById(1).Return(model.Patient{}, mockErr)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/", func(ctx *gin.Context) { ctx.Set("patientId", 1) }, mobileH.GetProfile)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 404, recorder.Code)
	})
	t.Run("internalError", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		mobileH := mobile.MobileHandler{Repo: repo}

		mockErr := fmt.Errorf("wrap : %w", errors.New("err"))
		repo.EXPECT().GetPatientById(1).Return(model.Patient{}, mockErr)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/", func(ctx *gin.Context) { ctx.Set("patientId", 1) }, mobileH.GetProfile)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("success", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		mobileH := mobile.MobileHandler{Repo: repo}

		repo.EXPECT().GetPatientById(1).Return(model.Patient{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/", func(ctx *gin.Context) { ctx.Set("patientId", 1) }, mobileH.GetProfile)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 200, recorder.Code)
	})
}
