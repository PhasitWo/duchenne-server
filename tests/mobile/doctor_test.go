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

func TestGetAllDoctor(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("internalError", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		mobileH := mobile.MobileHandler{Repo: repo}

		repo.EXPECT().GetAllDoctor().Return([]model.TrimDoctor{}, errors.New("err"))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/", mobileH.GetAllDoctor)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("success", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		mobileH := mobile.MobileHandler{Repo: repo}

		repo.EXPECT().GetAllDoctor().Return([]model.TrimDoctor{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/", mobileH.GetAllDoctor)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 200, recorder.Code)
	})
}
