package mobile_test

import (
	"bytes"
	"encoding/json"
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
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("bindingError", func(t *testing.T) {
		input := gin.H{
			"hn":         "test",
			"firstName":  "fn",
			"lastName":   "ln",
			"DeviceName": "goTest",
			"expoToken":  "", // error require
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		mobileH := mobile.MobileHandler{}

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/", mobileH.Login)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 400, recorder.Code)
	})
	t.Run("notFound", func(t *testing.T) {
		input := gin.H{
			"hn":         "test",
			"firstName":  "fn",
			"lastName":   "ln",
			"DeviceName": "goTest",
			"expoToken":  "expo",
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		mobileH := mobile.MobileHandler{Repo: repo}

		mockErr := fmt.Errorf("wrap : %w", gorm.ErrRecordNotFound)
		repo.EXPECT().GetPatientByHN("test").Return(model.Patient{}, mockErr)

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/", mobileH.Login)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 404, recorder.Code)
	})
	t.Run("getPatientInternalError", func(t *testing.T) {
		input := gin.H{
			"hn":         "test",
			"firstName":  "fn",
			"lastName":   "ln",
			"DeviceName": "goTest",
			"expoToken":  "expo",
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		mobileH := mobile.MobileHandler{Repo: repo}

		mockErr := fmt.Errorf("wrap : %w", errors.New("err"))
		repo.EXPECT().GetPatientByHN("test").Return(model.Patient{}, mockErr)

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/", mobileH.Login)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("unverifiedAccount", func(t *testing.T) {
		input := gin.H{
			"hn":         "test",
			"firstName":  "fn",
			"lastName":   "ln",
			"DeviceName": "goTest",
			"expoToken":  "expo",
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		mobileH := mobile.MobileHandler{Repo: repo}

		repo.EXPECT().GetPatientByHN("test").Return(model.Patient{Verified: false}, nil)

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/", mobileH.Login)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 403, recorder.Code)
	})
	t.Run("invalidCredential", func(t *testing.T) {
		input := gin.H{
			"hn":         "test",
			"firstName":  "fn",
			"lastName":   "ln",
			"DeviceName": "goTest",
			"expoToken":  "expo",
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)

		patient := model.Patient{
			FirstName: "somefn",
			LastName:  "someln",
			Verified: true,
		}
		// setup mock
		repo := repository.NewMockRepo(t)
		mobileH := mobile.MobileHandler{Repo: repo}

		repo.EXPECT().GetPatientByHN("test").Return(patient, nil)

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/", mobileH.Login)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 401, recorder.Code)
	})
	t.Run("getAllDeviceInternalError", func(t *testing.T) {
		input := gin.H{
			"hn":         "test",
			"firstName":  "fn",
			"lastName":   "ln",
			"DeviceName": "goTest",
			"expoToken":  "expo",
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)

		patient := model.Patient{
			FirstName: "fn",
			LastName:  "ln",
			Verified: true,
		}
		// setup mock
		repo := repository.NewMockRepo(t)
		mobileH := mobile.MobileHandler{Repo: repo}

		repo.EXPECT().GetPatientByHN("test").Return(patient, nil)
		repo.EXPECT().GetAllDevice(mock.Anything).Return([]model.Device{}, errors.New("err"))

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/", mobileH.Login)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("txBeginError", func(t *testing.T) {
		input := gin.H{
			"hn":         "test",
			"firstName":  "fn",
			"lastName":   "ln",
			"DeviceName": "goTest",
			"expoToken":  "expo",
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)

		patient := model.Patient{
			FirstName: "fn",
			LastName:  "ln",
			Verified: true,
		}
		mg := &gorm.DB{Error: errors.New("err")}
		// setup mock
		repo := repository.NewMockRepo(t)
		g := repository.NewMockGorm(t)
		mobileH := mobile.MobileHandler{Repo: repo, DBConn: g}

		g.EXPECT().Begin().Return(mg)
		repo.EXPECT().GetPatientByHN("test").Return(patient, nil)
		repo.EXPECT().GetAllDevice(mock.Anything).Return([]model.Device{}, nil)

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/", mobileH.Login)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
}
