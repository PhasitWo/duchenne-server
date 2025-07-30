package web_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PhasitWo/duchenne-server/handlers/web"
	"github.com/PhasitWo/duchenne-server/model"
	"github.com/PhasitWo/duchenne-server/services/notification"
	"github.com/PhasitWo/duchenne-server/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
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

func TestGetOneQuestion(t *testing.T) {
	t.Run("notFound", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		mockErr := fmt.Errorf("wrap : %w", gorm.ErrRecordNotFound)
		repo.EXPECT().GetQuestion("1").Return(model.SafeQuestion{}, mockErr)

		req := httptest.NewRequest(http.MethodGet, "/1", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/:id", webH.GetQuestion)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 404, recorder.Code)
	})
	t.Run("internalError", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		mockErr := fmt.Errorf("wrap : %w", errors.New("some internal error"))
		repo.EXPECT().GetQuestion("1").Return(model.SafeQuestion{}, mockErr)

		req := httptest.NewRequest(http.MethodGet, "/1", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/:id", webH.GetQuestion)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("success", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().GetQuestion("1").Return(model.SafeQuestion{}, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/1", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/:id", webH.GetQuestion)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 200, recorder.Code)
	})
}

func TestAnswerQuestion(t *testing.T) {
	t.Run("atoiError", func(t *testing.T) {
		// setup mock
		webH := web.WebHandler{}

		req := httptest.NewRequest(http.MethodPost, "/asd", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/:id", func(ctx *gin.Context) { ctx.Set("doctorId", 1) }, webH.AnswerQuestion)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 400, recorder.Code)
	})
	t.Run("notFound", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		mockErr := fmt.Errorf("wrap : %w", gorm.ErrRecordNotFound)
		repo.EXPECT().GetQuestion("1").Return(model.SafeQuestion{}, mockErr)

		req := httptest.NewRequest(http.MethodPost, "/1", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/:id", webH.AnswerQuestion)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 404, recorder.Code)
	})
	t.Run("getQuestionInternalError", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		mockErr := fmt.Errorf("wrap : %w", errors.New("some internal error"))
		repo.EXPECT().GetQuestion("1").Return(model.SafeQuestion{}, mockErr)

		req := httptest.NewRequest(http.MethodPost, "/1", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/:id", webH.AnswerQuestion)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("repliedQuestion", func(t *testing.T) {
		answerAt := 123
		question := model.SafeQuestion{
			Question: model.Question{ID: 1, AnswerAt: &answerAt},
		}
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().GetQuestion("1").Return(question, nil)

		req := httptest.NewRequest(http.MethodPost, "/1", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/:id", webH.AnswerQuestion)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 409, recorder.Code)
	})
	t.Run("bindingError", func(t *testing.T) {
		question := model.SafeQuestion{
			Question: model.Question{ID: 1, AnswerAt: nil},
		}
		input := model.QuestionAnswerRequest{
			Answer: "", // error require
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().GetQuestion("1").Return(question, nil)

		req := httptest.NewRequest(http.MethodPost, "/1", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/:id", webH.AnswerQuestion)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 400, recorder.Code)
	})
	t.Run("noDoctorIdFromAuthMiddleware", func(t *testing.T) {
		question := model.SafeQuestion{
			Question: model.Question{ID: 1, AnswerAt: nil},
		}
		input := model.QuestionAnswerRequest{
			Answer: "hello",
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().GetQuestion("1").Return(question, nil)

		req := httptest.NewRequest(http.MethodPost, "/1", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/:id", webH.AnswerQuestion)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("internalError", func(t *testing.T) {
		question := model.SafeQuestion{
			Question: model.Question{ID: 1, AnswerAt: nil},
		}
		input := model.QuestionAnswerRequest{
			Answer: "hello",
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().GetQuestion("1").Return(question, nil).Once()
		repo.EXPECT().UpdateQuestionAnswer(question.ID, input.Answer, 1).Return(errors.New("error")).Once()

		req := httptest.NewRequest(http.MethodPost, "/1", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/:id", func(ctx *gin.Context) { ctx.Set("doctorId", 1) }, webH.AnswerQuestion)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("success", func(t *testing.T) {
		question := model.SafeQuestion{
			Question: model.Question{ID: 1, AnswerAt: nil},
		}
		input := model.QuestionAnswerRequest{
			Answer: "hello",
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		noti := notification.NewMockService(t)
		webH := web.WebHandler{Repo: repo, NotiService: noti}

		repo.EXPECT().GetQuestion("1").Return(question, nil).Once()
		repo.EXPECT().UpdateQuestionAnswer(question.ID, input.Answer, 1).Return(nil).Once()
		noti.EXPECT().SendNotiByPatientId(mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

		req := httptest.NewRequest(http.MethodPost, "/1", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/:id", func(ctx *gin.Context) { ctx.Set("doctorId", 1) }, webH.AnswerQuestion)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 200, recorder.Code)
	})
}
