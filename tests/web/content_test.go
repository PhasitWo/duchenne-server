package web_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PhasitWo/duchenne-server/handlers/web"
	"github.com/PhasitWo/duchenne-server/model"
	"github.com/PhasitWo/duchenne-server/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateContent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("bindingError", func(t *testing.T) {
		input := model.CreateContentRequest{
			Title:       "", // error require
			Body:        "hello",
			IsPublished: false,
			Order:       15,
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		webH := web.WebHandler{}

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/", webH.CreateContent)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 400, recorder.Code)
	})
	t.Run("internalError", func(t *testing.T) {
		input := model.CreateContentRequest{
			Title:       "hi",
			Body:        "hello",
			IsPublished: false,
			Order:       15,
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().CreateContent(model.Content{
			Title:       input.Title,
			Body:        input.Body,
			IsPublished: input.IsPublished,
			Order:       input.Order,
		}).Return(-1, errors.New("err"))

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/", webH.CreateContent)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("success", func(t *testing.T) {
		input := model.CreateContentRequest{
			Title:       "hi",
			Body:        "hello",
			IsPublished: false,
			Order:       15,
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().CreateContent(model.Content{
			Title:       input.Title,
			Body:        input.Body,
			IsPublished: input.IsPublished,
			Order:       input.Order,
		}).Return(15, nil)

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/", webH.CreateContent)
		router.ServeHTTP(recorder, req)

		expectRespBody, err := json.Marshal(gin.H{"id": 15})
		assert.NoError(t, err)

		assert.Equal(t, 201, recorder.Code)
		assert.Equal(t, expectRespBody, recorder.Body.Bytes())
	})
}

func TestUpdateContent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("atoiError", func(t *testing.T) {
		// setup mock
		webH := web.WebHandler{}

		req := httptest.NewRequest(http.MethodPut, "/asd", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT("/:id", webH.UpdateContent)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 400, recorder.Code)
	})
	t.Run("bindingError", func(t *testing.T) {
		input := model.CreateContentRequest{
			Title:       "", // error require
			Body:        "hello",
			IsPublished: false,
			Order:       15,
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		webH := web.WebHandler{}

		req := httptest.NewRequest(http.MethodPut, "/87", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT("/:id", webH.UpdateContent)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 400, recorder.Code)
	})
	t.Run("internalError", func(t *testing.T) {
		input := model.CreateContentRequest{
			Title:       "Hi",
			Body:        "hello",
			IsPublished: false,
			Order:       15,
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().UpdateContent(model.Content{
			ID:          87,
			Title:       input.Title,
			Body:        input.Body,
			IsPublished: input.IsPublished,
			Order:       input.Order,
		}).Return(errors.New("err"))

		req := httptest.NewRequest(http.MethodPut, "/87", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT("/:id", webH.UpdateContent)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("success", func(t *testing.T) {
		input := model.CreateContentRequest{
			Title:       "Hi",
			Body:        "hello",
			IsPublished: false,
			Order:       15,
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().UpdateContent(model.Content{
			ID:          87,
			Title:       input.Title,
			Body:        input.Body,
			IsPublished: input.IsPublished,
			Order:       input.Order,
		}).Return(nil)

		req := httptest.NewRequest(http.MethodPut, "/87", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT("/:id", webH.UpdateContent)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 200, recorder.Code)
	})
}

func TestDeleteContent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("atoiError", func(t *testing.T) {
		// setup mock
		webH := web.WebHandler{}

		req := httptest.NewRequest(http.MethodDelete, "/asd", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.DELETE("/:id", webH.DeleteContent)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 400, recorder.Code)
	})
	t.Run("internalError", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().DeleteContent(65).Return(errors.New("err"))

		req := httptest.NewRequest(http.MethodDelete, "/65", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.DELETE("/:id", webH.DeleteContent)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("success", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().DeleteContent(65).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/65", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.DELETE("/:id", webH.DeleteContent)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 204, recorder.Code)
	})
}
