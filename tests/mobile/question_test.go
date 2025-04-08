package mobile_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/PhasitWo/duchenne-server/handlers/mobile"
	"github.com/PhasitWo/duchenne-server/model"
	"github.com/PhasitWo/duchenne-server/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestGetAllPatientQuestion(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("noPatientIdFromAuthMiddleware", func(t *testing.T) {
		// setup mock
		mobileH := mobile.MobileHandler{}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/", mobileH.GetAllPatientQuestion)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("internalError", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		mobileH := mobile.MobileHandler{Repo: repo}

		repo.EXPECT().GetAllQuestion(30, 0, mock.Anything).Return([]model.QuestionTopic{}, errors.New("err"))

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/", func(ctx *gin.Context) { ctx.Set("patientId", 1) }, mobileH.GetAllPatientQuestion)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("success", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		mobileH := mobile.MobileHandler{Repo: repo}

		repo.EXPECT().GetAllQuestion(30, 0, mock.Anything).Return([]model.QuestionTopic{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/", func(ctx *gin.Context) { ctx.Set("patientId", 1) }, mobileH.GetAllPatientQuestion)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 200, recorder.Code)
	})
}

func TestGetOneQuestion(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("noPatientIdFromAuthMiddleware", func(t *testing.T) {
		// setup mock
		mobileH := mobile.MobileHandler{}

		req := httptest.NewRequest(http.MethodGet, "/10", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/:id", mobileH.GetQuestion)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("notFound", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		mobileH := mobile.MobileHandler{Repo: repo}

		mockErr := fmt.Errorf("wrap : %w", gorm.ErrRecordNotFound)
		repo.EXPECT().GetQuestion("10").Return(model.SafeQuestion{}, mockErr)

		req := httptest.NewRequest(http.MethodGet, "/10", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/:id", func(ctx *gin.Context) { ctx.Set("patientId", 1) }, mobileH.GetQuestion)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 404, recorder.Code)
	})
	t.Run("internalError", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		mobileH := mobile.MobileHandler{Repo: repo}

		mockErr := fmt.Errorf("wrap : %w", errors.New("err"))
		repo.EXPECT().GetQuestion("10").Return(model.SafeQuestion{}, mockErr)

		req := httptest.NewRequest(http.MethodGet, "/10", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/:id", func(ctx *gin.Context) { ctx.Set("patientId", 1) }, mobileH.GetQuestion)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("unauthorized", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		mobileH := mobile.MobileHandler{Repo: repo}

		repo.EXPECT().GetQuestion("10").Return(model.SafeQuestion{Question: model.Question{Patient: model.Patient{ID: 99}}}, nil)

		req := httptest.NewRequest(http.MethodGet, "/10", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/:id", func(ctx *gin.Context) { ctx.Set("patientId", 1) }, mobileH.GetQuestion)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 401, recorder.Code)
	})
	t.Run("success", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		mobileH := mobile.MobileHandler{Repo: repo}

		repo.EXPECT().GetQuestion("10").Return(model.SafeQuestion{Question: model.Question{Patient: model.Patient{ID: 1}}}, nil)

		req := httptest.NewRequest(http.MethodGet, "/10", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/:id", func(ctx *gin.Context) { ctx.Set("patientId", 1) }, mobileH.GetQuestion)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 200, recorder.Code)
	})
}

func TestCreateQuestion(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("noPatientIdFromAuthMiddleware", func(t *testing.T) {
		// setup mock
		mobileH := mobile.MobileHandler{}

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/", mobileH.CreateQuestion)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("bindingErrorRequiredTopic", func(t *testing.T) {
		input := model.CreateQuestionRequest{
			Topic:    "", // error require
			Question: "im curious",
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		mobileH := mobile.MobileHandler{}

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/", func(ctx *gin.Context) { ctx.Set("patientId", 1) }, mobileH.CreateQuestion)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 400, recorder.Code)
	})
	t.Run("bindingErrorRequiredQuestion", func(t *testing.T) {
		input := model.CreateQuestionRequest{
			Topic:    "hello",
			Question: "", // error require
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		mobileH := mobile.MobileHandler{}

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/", func(ctx *gin.Context) { ctx.Set("patientId", 1) }, mobileH.CreateQuestion)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 400, recorder.Code)
	})
	t.Run("exceedTopicLength", func(t *testing.T) {
		input := model.CreateQuestionRequest{
			Topic:    strings.Repeat("x", 51),
			Question: "im curious",
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		mobileH := mobile.MobileHandler{}

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/", func(ctx *gin.Context) { ctx.Set("patientId", 1) }, mobileH.CreateQuestion)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 422, recorder.Code)
	})
	t.Run("exceedQuestionLength", func(t *testing.T) {
		input := model.CreateQuestionRequest{
			Topic:    "hello",
			Question: strings.Repeat("x", 701),
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		mobileH := mobile.MobileHandler{}

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/", func(ctx *gin.Context) { ctx.Set("patientId", 1) }, mobileH.CreateQuestion)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 422, recorder.Code)
	})
	t.Run("internalError", func(t *testing.T) {
		input := model.CreateQuestionRequest{
			Topic:    "hello",
			Question: "im curious",
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		mobileH := mobile.MobileHandler{Repo: repo}

		repo.EXPECT().CreateQuestion(1, input.Topic, input.Question, mock.Anything).Return(-1, errors.New("err"))

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/", func(ctx *gin.Context) { ctx.Set("patientId", 1) }, mobileH.CreateQuestion)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("success", func(t *testing.T) {
		input := model.CreateQuestionRequest{
			Topic:    "hello",
			Question: "im curious",
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		mobileH := mobile.MobileHandler{Repo: repo}

		repo.EXPECT().CreateQuestion(1, input.Topic, input.Question, mock.Anything).Return(55, nil)

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/", func(ctx *gin.Context) { ctx.Set("patientId", 1) }, mobileH.CreateQuestion)
		router.ServeHTTP(recorder, req)

		expectRespBody, err := json.Marshal(gin.H{"id": 55})
		assert.NoError(t, err)

		assert.Equal(t, 201, recorder.Code)
		assert.Equal(t, expectRespBody, recorder.Body.Bytes())
	})
}

func TestDeleteQuestion(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("noPatientIdFromAuthMiddleware", func(t *testing.T) {
		// setup mock
		mobileH := mobile.MobileHandler{}

		req := httptest.NewRequest(http.MethodDelete, "/15", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.DELETE("/:id", mobileH.DeleteQuestion)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("notFound", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		mobileH := mobile.MobileHandler{Repo: repo}

		mockErr := fmt.Errorf("wrap : %w", gorm.ErrRecordNotFound)
		repo.EXPECT().GetQuestion("15").Return(model.SafeQuestion{}, mockErr)

		req := httptest.NewRequest(http.MethodDelete, "/15", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.DELETE("/:id", func(ctx *gin.Context) { ctx.Set("patientId", 1) }, mobileH.DeleteQuestion)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 404, recorder.Code)
	})
	t.Run("getQuestionInternalError", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		mobileH := mobile.MobileHandler{Repo: repo}

		mockErr := fmt.Errorf("wrap : %w", errors.New("err"))
		repo.EXPECT().GetQuestion("15").Return(model.SafeQuestion{}, mockErr)

		req := httptest.NewRequest(http.MethodDelete, "/15", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.DELETE("/:id", func(ctx *gin.Context) { ctx.Set("patientId", 1) }, mobileH.DeleteQuestion)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("unauthorized", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		mobileH := mobile.MobileHandler{Repo: repo}

		repo.EXPECT().GetQuestion("15").Return(model.SafeQuestion{Question: model.Question{Patient: model.Patient{ID: 99}}}, nil)

		req := httptest.NewRequest(http.MethodDelete, "/15", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.DELETE("/:id", func(ctx *gin.Context) { ctx.Set("patientId", 1) }, mobileH.DeleteQuestion)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 401, recorder.Code)
	})
	t.Run("deleteInternalError", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		mobileH := mobile.MobileHandler{Repo: repo}

		repo.EXPECT().GetQuestion("15").Return(model.SafeQuestion{Question: model.Question{Patient: model.Patient{ID: 1}}}, nil)
		repo.EXPECT().DeleteQuestion("15").Return(errors.New("err"))

		req := httptest.NewRequest(http.MethodDelete, "/15", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.DELETE("/:id", func(ctx *gin.Context) { ctx.Set("patientId", 1) }, mobileH.DeleteQuestion)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("success", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		mobileH := mobile.MobileHandler{Repo: repo}

		repo.EXPECT().GetQuestion("15").Return(model.SafeQuestion{Question: model.Question{Patient: model.Patient{ID: 1}}}, nil)
		repo.EXPECT().DeleteQuestion("15").Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/15", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.DELETE("/:id", func(ctx *gin.Context) { ctx.Set("patientId", 1) }, mobileH.DeleteQuestion)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 204, recorder.Code)
	})
}
