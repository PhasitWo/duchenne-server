package common_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PhasitWo/duchenne-server/handlers/common"
	"github.com/PhasitWo/duchenne-server/model"
	"github.com/PhasitWo/duchenne-server/repository"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetAllContent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("limitParamError", func(t *testing.T) {
		// setup mock
		commonH := common.CommonHandler{}

		req := httptest.NewRequest(http.MethodGet, "/?limit=asd", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/", commonH.GetAllContent)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 400, recorder.Code)
	})
	t.Run("offsetParamError", func(t *testing.T) {
		// setup mock
		commonH := common.CommonHandler{}

		req := httptest.NewRequest(http.MethodGet, "/?offset=asd", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/", commonH.GetAllContent)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 400, recorder.Code)
	})
	t.Run("internalError", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		commonH := common.CommonHandler{Repo: repo}

		repo.EXPECT().GetAllContent(9999, 0).Return([]model.Content{}, errors.New("err"))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/", commonH.GetAllContent)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("success", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		commonH := common.CommonHandler{Repo: repo}

		repo.EXPECT().GetAllContent(9999, 0).Return([]model.Content{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/", commonH.GetAllContent)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 200, recorder.Code)
	})
}

func TestGetOneContent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("notFound", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		commonH := common.CommonHandler{Repo: repo}

		mockErr := fmt.Errorf("wrap : %w", gorm.ErrRecordNotFound)
		repo.EXPECT().GetContent("15").Return(model.Content{}, mockErr)

		req := httptest.NewRequest(http.MethodGet, "/15", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/:id", commonH.GetOneContent)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 404, recorder.Code)
	})
	t.Run("internalError", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		commonH := common.CommonHandler{Repo: repo}

		mockErr := fmt.Errorf("wrap : %w", errors.New("err"))
		repo.EXPECT().GetContent("15").Return(model.Content{}, mockErr)

		req := httptest.NewRequest(http.MethodGet, "/15", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/:id", commonH.GetOneContent)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("success", func(t *testing.T) {
		content := model.Content{
			ID:          15,
			Title:       "hello",
			Body:        "HI",
			Order:       1,
			IsPublished: true,
		}
		// setup mock
		repo := repository.NewMockRepo(t)
		commonH := common.CommonHandler{Repo: repo}

		repo.EXPECT().GetContent("15").Return(content, nil)

		req := httptest.NewRequest(http.MethodGet, "/15", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/:id", commonH.GetOneContent)
		router.ServeHTTP(recorder, req)

		expectRespBody, err := json.Marshal(&content)
		assert.NoError(t, err)

		assert.Equal(t, 200, recorder.Code)
		assert.Equal(t, expectRespBody, recorder.Body.Bytes())
	})
}
