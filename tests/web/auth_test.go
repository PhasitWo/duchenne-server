package web_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PhasitWo/duchenne-server/auth"
	"github.com/PhasitWo/duchenne-server/config"
	"github.com/PhasitWo/duchenne-server/handlers/web"
	"github.com/PhasitWo/duchenne-server/model"
	"github.com/PhasitWo/duchenne-server/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config.AppConfig.JWT_KEY = "test_key"
	t.Run("success", func(t *testing.T) {
		doctor := model.Doctor{
			ID: 1,
			Username: "test",
			Password: "admin",
			Role: "root",
		}
		input := gin.H{
			"username": "test",
			"password": "admin",
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().GetDoctorByUsername("test").Return(doctor, nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)
		
		expectToken, err := auth.GenerateDoctorToken(doctor.ID, doctor.Role)
		assert.NoError(t, err)
		expectResp:= gin.H{
			"token" : expectToken,
		}
		expectRespBody, err := json.Marshal(expectResp)
		assert.NoError(t, err)

		router.POST("/", webH.Login)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 200, recorder.Code)
		assert.Equal(t, expectRespBody, recorder.Body.Bytes())
	})
	t.Run("bindingError", func(t *testing.T) {
		input := gin.H{
			"username": "", // error require
			"password": "", // error require
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		webH := web.WebHandler{}

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/", webH.Login)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 400, recorder.Code)
	})
	t.Run("notFound", func(t *testing.T) {
		input := gin.H{
			"username": "test",
			"password": "admin",
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		mockErr := fmt.Errorf("wrap : %w", gorm.ErrRecordNotFound)
		repo.EXPECT().GetDoctorByUsername("test").Return(model.Doctor{}, mockErr).Once()

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)
		
		router.POST("/", webH.Login)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 404, recorder.Code)
	})
	t.Run("internalError", func(t *testing.T) {
		input := gin.H{
			"username": "test",
			"password": "admin",
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		mockErr := fmt.Errorf("wrap : %w", errors.New("some internal error"))
		repo.EXPECT().GetDoctorByUsername("test").Return(model.Doctor{}, mockErr).Once()

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)
		
		router.POST("/", webH.Login)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("invalidCredential", func(t *testing.T) {
		doctor := model.Doctor{
			ID: 1,
			Username: "test",
			Password: "admin",
			Role: "root",
		}
		input := gin.H{
			"username": "test",
			"password": "random",
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().GetDoctorByUsername("test").Return(doctor, nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/", webH.Login)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 401, recorder.Code)
	})
}