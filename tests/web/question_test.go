package web_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"encoding/json"
	
	"github.com/PhasitWo/duchenne-server/handlers/web"
	"github.com/PhasitWo/duchenne-server/model"
	"github.com/PhasitWo/duchenne-server/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func testInternalQueryParam(param string) func(t *testing.T) {
	return func(t *testing.T) {
		// setup mock
		webH := web.WebHandler{}

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/?%s=asdsad", param), nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/", webH.GetAllQuestion)
		router.ServeHTTP(recorder, req)

		expectRespBody, err := json.Marshal(gin.H{"error": fmt.Sprintf("cannot parse %s value", param)})
		assert.NoError(t, err)

		assert.Equal(t, 400, recorder.Code)
		assert.Equal(t, expectRespBody, recorder.Body.Bytes())
	}
}

func TestGetAllQuestion(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("limitError", testInternalQueryParam("limit"))
	t.Run("offsetError", testInternalQueryParam("offset"))
	t.Run("doctorIdError", testInternalQueryParam("doctorId"))
	t.Run("patientIdError", testInternalQueryParam("patientId"))
	t.Run("invalidType", func(t *testing.T) {
		// setup mock
		webH := web.WebHandler{}

		req := httptest.NewRequest(http.MethodGet, "/?type=asasd", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/", webH.GetAllQuestion)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 400, recorder.Code)
	})
	t.Run("internalError", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().GetAllQuestion(9999, 0).Return([]model.QuestionTopic{}, errors.New("some internal error"))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/", webH.GetAllQuestion)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("success", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().GetAllQuestion(9999, 0).Return([]model.QuestionTopic{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/", webH.GetAllQuestion)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 200, recorder.Code)
	})
}
