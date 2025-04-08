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

func TestGetDoctor(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("success", func(t *testing.T) {
		doctor := model.Doctor{
			ID: 1,
		}
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().GetDoctorById("1").Return(doctor, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/1", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		expectRespBody, err := json.Marshal(&doctor)
		assert.NoError(t, err)

		router.GET("/:id", webH.GetDoctor)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 200, recorder.Code)
		assert.Equal(t, expectRespBody, recorder.Body.Bytes())
	})
	t.Run("notFound", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		mockErr := fmt.Errorf("wrap : %w", gorm.ErrRecordNotFound)
		repo.EXPECT().GetDoctorById("1").Return(model.Doctor{}, mockErr).Once()

		req := httptest.NewRequest(http.MethodGet, "/1", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/:id", webH.GetDoctor)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 404, recorder.Code)
	})
	t.Run("internalError", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		mockErr := fmt.Errorf("wrap : %w", errors.New("some internal error"))
		repo.EXPECT().GetDoctorById("1").Return(model.Doctor{}, mockErr).Once()

		req := httptest.NewRequest(http.MethodGet, "/1", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/:id", webH.GetDoctor)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
}

func TestGetAllDoctor(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().GetAllDoctor().Return([]model.TrimDoctor{}, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/", webH.GetAllDoctor)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 200, recorder.Code)
	})
	t.Run("internalError", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().GetAllDoctor().Return([]model.TrimDoctor{}, errors.New("some internal error")).Once()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/", webH.GetAllDoctor)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
}

func TestCreateDoctor(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		input := model.CreateDoctorRequest{
			FirstName:  "fn",
			MiddleName: nil,
			LastName:   "ln",
			Username:   "testusername",
			Password:   "admin",
			Role:       model.ADMIN,
			Specialist: nil,
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().CreateDoctor(model.Doctor{
			FirstName:  input.FirstName,
			MiddleName: input.MiddleName,
			LastName:   input.LastName,
			Username:   input.Username,
			Password:   input.Password,
			Role:       input.Role,
			Specialist: input.Specialist,
		}).Return(1, nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/", webH.CreateDoctor)
		router.ServeHTTP(recorder, req)

		expectRespBody, err := json.Marshal(&gin.H{"id": 1})
		assert.NoError(t, err)

		assert.Equal(t, 201, recorder.Code)
		assert.Equal(t, expectRespBody, recorder.Body.Bytes())
	})
	t.Run("bindingError", func(t *testing.T) {
		input := model.CreateDoctorRequest{
			FirstName:  "", // error require
			MiddleName: nil,
			LastName:   "ln",
			Username:   "testusername",
			Password:   "admin",
			Role:       model.ADMIN,
			Specialist: nil,
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		webH := web.WebHandler{}

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/", webH.CreateDoctor)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 400, recorder.Code)
	})
	t.Run("invalidRole", func(t *testing.T) {
		input := model.CreateDoctorRequest{
			FirstName:  "fn",
			MiddleName: nil,
			LastName:   "ln",
			Username:   "testusername",
			Password:   "admin",
			Role:       "brabra",
			Specialist: nil,
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		webH := web.WebHandler{}

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/", webH.CreateDoctor)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 400, recorder.Code)
	})
	t.Run("duplicateUsername", func(t *testing.T) {
		input := model.CreateDoctorRequest{
			FirstName:  "fn",
			MiddleName: nil,
			LastName:   "ln",
			Username:   "testusername",
			Password:   "admin",
			Role:       model.ADMIN,
			Specialist: nil,
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		mockErr := fmt.Errorf("wrap : %w", repository.ErrDuplicateEntry)
		repo.EXPECT().CreateDoctor(model.Doctor{
			FirstName:  input.FirstName,
			MiddleName: input.MiddleName,
			LastName:   input.LastName,
			Username:   input.Username,
			Password:   input.Password,
			Role:       input.Role,
			Specialist: input.Specialist,
		}).Return(1, mockErr).Once()

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/", webH.CreateDoctor)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 409, recorder.Code)
	})
	t.Run("internalError", func(t *testing.T) {
		input := model.CreateDoctorRequest{
			FirstName:  "fn",
			MiddleName: nil,
			LastName:   "ln",
			Username:   "testusername",
			Password:   "admin",
			Role:       model.ADMIN,
			Specialist: nil,
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		mockErr := fmt.Errorf("wrap : %w", errors.New("some internal error"))
		repo.EXPECT().CreateDoctor(model.Doctor{
			FirstName:  input.FirstName,
			MiddleName: input.MiddleName,
			LastName:   input.LastName,
			Username:   input.Username,
			Password:   input.Password,
			Role:       input.Role,
			Specialist: input.Specialist,
		}).Return(1, mockErr).Once()

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.POST("/", webH.CreateDoctor)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
}

func TestUpdateDoctor(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		input := model.CreateDoctorRequest{
			FirstName:  "fn",
			MiddleName: nil,
			LastName:   "ln",
			Username:   "testusername",
			Password:   "admin",
			Role:       model.ADMIN,
			Specialist: nil,
		}
		doctor := model.Doctor{
			ID: 1,
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().GetDoctorById(1).Return(doctor, nil).Once()

		repo.EXPECT().UpdateDoctor(model.Doctor{
			ID:         doctor.ID,
			FirstName:  input.FirstName,
			MiddleName: input.MiddleName,
			LastName:   input.LastName,
			Username:   input.Username,
			Password:   input.Password,
			Role:       input.Role,
			Specialist: input.Specialist,
		}).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPut, "/1", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT("/:id", webH.UpdateDoctor)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 200, recorder.Code)
	})
	t.Run("bindingError", func(t *testing.T) {
		input := model.CreateDoctorRequest{
			FirstName:  "", // error require
			MiddleName: nil,
			LastName:   "ln",
			Username:   "testusername",
			Password:   "admin",
			Role:       model.ADMIN,
			Specialist: nil,
		}

		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		webH := web.WebHandler{}

		req := httptest.NewRequest(http.MethodPut, "/1", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT("/:id", webH.UpdateDoctor)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 400, recorder.Code)
	})
	t.Run("invalidRole", func(t *testing.T) {
		input := model.CreateDoctorRequest{
			FirstName:  "fn",
			MiddleName: nil,
			LastName:   "ln",
			Username:   "testusername",
			Password:   "admin",
			Role:       "brabra",
			Specialist: nil,
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		webH := web.WebHandler{}

		req := httptest.NewRequest(http.MethodPut, "/1", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT("/:id", webH.UpdateDoctor)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 400, recorder.Code)
	})
	t.Run("atoiError", func(t *testing.T) {
		input := model.CreateDoctorRequest{
			FirstName:  "fn",
			MiddleName: nil,
			LastName:   "ln",
			Username:   "testusername",
			Password:   "admin",
			Role:       model.ADMIN,
			Specialist: nil,
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		webH := web.WebHandler{}

		req := httptest.NewRequest(http.MethodPut, "/asd", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT("/:id", webH.UpdateDoctor)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 400, recorder.Code)
	})
	t.Run("notFound", func(t *testing.T) {
		input := model.CreateDoctorRequest{
			FirstName:  "fn",
			MiddleName: nil,
			LastName:   "ln",
			Username:   "testusername",
			Password:   "admin",
			Role:       model.ADMIN,
			Specialist: nil,
		}
		doctor := model.Doctor{
			ID: 1,
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		mockErr := fmt.Errorf("wrap : %w", gorm.ErrRecordNotFound)
		repo.EXPECT().GetDoctorById(1).Return(doctor, mockErr).Once()

		req := httptest.NewRequest(http.MethodPut, "/1", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT("/:id", webH.UpdateDoctor)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 404, recorder.Code)
	})
	t.Run("getDoctorInternalError", func(t *testing.T) {
		input := model.CreateDoctorRequest{
			FirstName:  "fn",
			MiddleName: nil,
			LastName:   "ln",
			Username:   "testusername",
			Password:   "admin",
			Role:       model.ADMIN,
			Specialist: nil,
		}
		doctor := model.Doctor{
			ID: 1,
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		mockErr := fmt.Errorf("wrap : %w", errors.New("some internal error"))
		repo.EXPECT().GetDoctorById(1).Return(doctor, mockErr).Once()

		req := httptest.NewRequest(http.MethodPut, "/1", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT("/:id", webH.UpdateDoctor)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("duplicateUsername", func(t *testing.T) {
		input := model.CreateDoctorRequest{
			FirstName:  "fn",
			MiddleName: nil,
			LastName:   "ln",
			Username:   "testusername",
			Password:   "admin",
			Role:       model.ADMIN,
			Specialist: nil,
		}
		doctor := model.Doctor{
			ID: 1,
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().GetDoctorById(1).Return(doctor, nil).Once()
		mockErr := fmt.Errorf("wrap : %w", repository.ErrDuplicateEntry)
		repo.EXPECT().UpdateDoctor(model.Doctor{
			ID:         doctor.ID,
			FirstName:  input.FirstName,
			MiddleName: input.MiddleName,
			LastName:   input.LastName,
			Username:   input.Username,
			Password:   input.Password,
			Role:       input.Role,
			Specialist: input.Specialist,
		}).Return(mockErr).Once()

		req := httptest.NewRequest(http.MethodPut, "/1", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT("/:id", webH.UpdateDoctor)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 409, recorder.Code)
	})
	t.Run("internalError", func(t *testing.T) {
		input := model.CreateDoctorRequest{
			FirstName:  "fn",
			MiddleName: nil,
			LastName:   "ln",
			Username:   "testusername",
			Password:   "admin",
			Role:       model.ADMIN,
			Specialist: nil,
		}
		doctor := model.Doctor{
			ID: 1,
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().GetDoctorById(1).Return(doctor, nil).Once()
		mockErr := fmt.Errorf("wrap : %w", errors.New("some internal error"))
		repo.EXPECT().UpdateDoctor(model.Doctor{
			ID:         doctor.ID,
			FirstName:  input.FirstName,
			MiddleName: input.MiddleName,
			LastName:   input.LastName,
			Username:   input.Username,
			Password:   input.Password,
			Role:       input.Role,
			Specialist: input.Specialist,
		}).Return(mockErr).Once()

		req := httptest.NewRequest(http.MethodPut, "/1", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT("/:id", webH.UpdateDoctor)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
}

func TestDeleteDoctor(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().DeleteDoctorById(1).Return(nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/1", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.DELETE("/:id", webH.DeleteDoctor)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 204, recorder.Code)
	})
	t.Run("atoiError", func(t *testing.T) {
		// setup mock
		webH := web.WebHandler{}
		req := httptest.NewRequest(http.MethodDelete, "/asd", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.DELETE("/:id", webH.DeleteDoctor)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 400, recorder.Code)
	})
	t.Run("internalError", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().DeleteDoctorById(1).Return(errors.New("some internal error")).Once()

		req := httptest.NewRequest(http.MethodDelete, "/1", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.DELETE("/:id", webH.DeleteDoctor)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
}
