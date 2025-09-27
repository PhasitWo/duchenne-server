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
	"github.com/PhasitWo/duchenne-server/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGetProfile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("noDoctorIdFromAuthMiddleware", func(t *testing.T) {
		// setup mock
		webH := web.WebHandler{}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/", webH.GetProfile)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("notFound", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		mockErr := fmt.Errorf("wrap : %w", gorm.ErrRecordNotFound)
		repo.EXPECT().GetDoctorById(1).Return(model.Doctor{}, mockErr).Once()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/", func(ctx *gin.Context) { ctx.Set("doctorId", 1) }, webH.GetProfile)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 404, recorder.Code)
	})
	t.Run("internalError", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		mockErr := fmt.Errorf("wrap : %w", errors.New("some internal error"))
		repo.EXPECT().GetDoctorById(1).Return(model.Doctor{}, mockErr).Once()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/", func(ctx *gin.Context) { ctx.Set("doctorId", 1) }, webH.GetProfile)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("success", func(t *testing.T) {
		doctor := model.Doctor{ID: 1}
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().GetDoctorById(1).Return(doctor, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/", func(ctx *gin.Context) { ctx.Set("doctorId", 1) }, webH.GetProfile)
		router.ServeHTTP(recorder, req)

		expectRespBody, err := json.Marshal(&doctor)
		assert.NoError(t, err)

		assert.Equal(t, 200, recorder.Code)
		assert.Equal(t, expectRespBody, recorder.Body.Bytes())
	})
}

func TestUpdateProfile(t *testing.T) {
	password := "1234"
	t.Run("noDoctorIdFromAuthMiddleware", func(t *testing.T) {
		// setup mock
		webH := web.WebHandler{}

		req := httptest.NewRequest(http.MethodPut, "/", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT("/", webH.UpdateProfile)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("noDoctorRoleFromAuthMiddleware", func(t *testing.T) {
		// setup mock
		webH := web.WebHandler{}

		req := httptest.NewRequest(http.MethodPut, "/", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT("/", func(ctx *gin.Context) { ctx.Set("doctorId", 1) }, webH.UpdateProfile)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("bindingError", func(t *testing.T) {
		input := model.UpdateProfileRequest{
			FirstName:  "", // error require
			MiddleName: nil,
			LastName:   "ln",
			Username:   "test",
			Password:   &password,
			Specialist: nil,
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		webH := web.WebHandler{}

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT(
			"/",
			func(ctx *gin.Context) {
				ctx.Set("doctorId", 1)
				ctx.Set("doctorRole", model.ADMIN)
			},
			webH.UpdateProfile,
		)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 400, recorder.Code)
	})
	t.Run("duplicateUsername", func(t *testing.T) {
		input := model.UpdateProfileRequest{
			FirstName:  "fn",
			MiddleName: nil,
			LastName:   "ln",
			Username:   "test",
			Password:   &password,
			Specialist: nil,
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		mockErr := fmt.Errorf("wrap : %w", repository.ErrDuplicateEntry)
		repo.EXPECT().UpdateDoctor(
			model.Doctor{
				ID:         1,
				FirstName:  input.FirstName,
				MiddleName: input.MiddleName,
				LastName:   input.LastName,
				Username:   input.Username,
				Password:   *input.Password,
				Specialist: input.Specialist,
				Role:       model.ADMIN,
			},
		).Return(mockErr).Once()

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT(
			"/",
			func(ctx *gin.Context) {
				ctx.Set("doctorId", 1)
				ctx.Set("doctorRole", model.ADMIN)
			},
			webH.UpdateProfile,
		)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 409, recorder.Code)
	})
	t.Run("internalError", func(t *testing.T) {
		input := model.UpdateProfileRequest{
			FirstName:  "fn",
			MiddleName: nil,
			LastName:   "ln",
			Username:   "test",
			Password:   &password,
			Specialist: nil,
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		mockErr := fmt.Errorf("wrap : %w", errors.New("some internal error"))
		repo.EXPECT().UpdateDoctor(
			model.Doctor{
				ID:         1,
				FirstName:  input.FirstName,
				MiddleName: input.MiddleName,
				LastName:   input.LastName,
				Username:   input.Username,
				Password:   *input.Password,
				Specialist: input.Specialist,
				Role:       model.ADMIN,
			},
		).Return(mockErr).Once()

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT(
			"/",
			func(ctx *gin.Context) {
				ctx.Set("doctorId", 1)
				ctx.Set("doctorRole", model.ADMIN)
			},
			webH.UpdateProfile,
		)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("success", func(t *testing.T) {
		input := model.UpdateProfileRequest{
			FirstName:  "fn",
			MiddleName: nil,
			LastName:   "ln",
			Username:   "test",
			Password:   &password,
			Specialist: nil,
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().UpdateDoctor(
			model.Doctor{
				ID:         1,
				FirstName:  input.FirstName,
				MiddleName: input.MiddleName,
				LastName:   input.LastName,
				Username:   input.Username,
				Password:   *input.Password,
				Specialist: input.Specialist,
				Role:       model.ADMIN,
			},
		).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT(
			"/",
			func(ctx *gin.Context) {
				ctx.Set("doctorId", 1)
				ctx.Set("doctorRole", model.ADMIN)
			},
			webH.UpdateProfile,
		)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 200, recorder.Code)
	})
}
