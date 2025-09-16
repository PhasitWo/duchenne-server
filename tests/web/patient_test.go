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
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestGetOnPatient(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("success", func(t *testing.T) {
		patient := model.Patient{
			ID: 1,
		}
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().GetPatientById("1").Return(patient, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/1", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		expectRespBody, err := json.Marshal(&patient)
		assert.NoError(t, err)

		router.GET("/:id", webH.GetPatient)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 200, recorder.Code)
		assert.Equal(t, expectRespBody, recorder.Body.Bytes())
	})
	t.Run("notFound", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		mockErr := fmt.Errorf("wrap : %w", gorm.ErrRecordNotFound)
		repo.EXPECT().GetPatientById("1").Return(model.Patient{}, mockErr).Once()

		req := httptest.NewRequest(http.MethodGet, "/1", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/:id", webH.GetPatient)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 404, recorder.Code)
	})
	t.Run("internalError", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		mockErr := fmt.Errorf("wrap : %w", errors.New("some internal error"))
		repo.EXPECT().GetPatientById("1").Return(model.Patient{}, mockErr).Once()

		req := httptest.NewRequest(http.MethodGet, "/1", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/:id", webH.GetPatient)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
}

func TestGetAllPatient(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().GetAllPatient(mock.Anything, mock.Anything).Return([]model.Patient{}, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/", webH.GetAllPatient)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 200, recorder.Code)
	})
	t.Run("internalError", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().GetAllPatient(mock.Anything, mock.Anything).Return([]model.Patient{}, errors.New("some internal error")).Once()

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.GET("/", webH.GetAllPatient)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
}

// func TestCreatePatient(t *testing.T) {
// 	t.Run("success", func(t *testing.T) {
// 		input := model.CreatePatientRequest{
// 			Hn:        "testhn",
// 			FirstName: "fn",
// 			LastName:  "ln",
// 			Verified:  false,
// 		}
// 		rawInput, err := json.Marshal(&input)
// 		assert.NoError(t, err)
// 		// setup mock
// 		repo := repository.NewMockRepo(t)
// 		webH := web.WebHandler{Repo: repo}

// 		repo.EXPECT().CreatePatient(
// 			model.Patient{
// 				Hn:        input.Hn,
// 				FirstName: input.FirstName,
// 				LastName:  input.LastName,
// 				Verified:  input.Verified,
// 			},
// 		).Return(1, nil).Once()

// 		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(rawInput))
// 		recorder := httptest.NewRecorder()
// 		_, router := gin.CreateTestContext(recorder)

// 		router.POST("/", webH.CreatePatient)
// 		router.ServeHTTP(recorder, req)

// 		expectRespBody, err := json.Marshal(gin.H{"id": 1})
// 		assert.NoError(t, err)

// 		assert.Equal(t, 201, recorder.Code)
// 		assert.Equal(t, expectRespBody, recorder.Body.Bytes())
// 	})
// 	t.Run("bindingError", func(t *testing.T) {
// 		input := model.CreatePatientRequest{
// 			Hn:        "", // error require
// 			FirstName: "fn",
// 			LastName:  "ln",
// 			Verified:  false,
// 		}
// 		rawInput, err := json.Marshal(&input)
// 		assert.NoError(t, err)
// 		// setup mock
// 		webH := web.WebHandler{}

// 		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(rawInput))
// 		recorder := httptest.NewRecorder()
// 		_, router := gin.CreateTestContext(recorder)

// 		router.POST("/", webH.CreatePatient)
// 		router.ServeHTTP(recorder, req)

// 		assert.Equal(t, 400, recorder.Code)
// 	})
// 	t.Run("duplicateHn", func(t *testing.T) {
// 		input := model.CreatePatientRequest{
// 			Hn:        "testhn",
// 			FirstName: "fn",
// 			LastName:  "ln",
// 			Verified:  false,
// 		}
// 		rawInput, err := json.Marshal(&input)
// 		assert.NoError(t, err)
// 		// setup mock
// 		repo := repository.NewMockRepo(t)
// 		webH := web.WebHandler{Repo: repo}

// 		mockErr := fmt.Errorf("wrap : %w", repository.ErrDuplicateEntry)
// 		repo.EXPECT().CreatePatient(
// 			model.Patient{
// 				Hn:        input.Hn,
// 				FirstName: input.FirstName,
// 				LastName:  input.LastName,
// 				Verified:  input.Verified,
// 			},
// 		).Return(1, mockErr).Once()

// 		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(rawInput))
// 		recorder := httptest.NewRecorder()
// 		_, router := gin.CreateTestContext(recorder)

// 		router.POST("/", webH.CreatePatient)
// 		router.ServeHTTP(recorder, req)

// 		assert.Equal(t, 409, recorder.Code)
// 	})
// 	t.Run("internalError", func(t *testing.T) {
// 		input := model.CreatePatientRequest{
// 			Hn:        "testhn",
// 			FirstName: "fn",
// 			LastName:  "ln",
// 			Verified:  false,
// 		}
// 		rawInput, err := json.Marshal(&input)
// 		assert.NoError(t, err)
// 		// setup mock
// 		repo := repository.NewMockRepo(t)
// 		webH := web.WebHandler{Repo: repo}

// 		mockErr := fmt.Errorf("wrap : %w", errors.New("some internal error"))
// 		repo.EXPECT().CreatePatient(
// 			model.Patient{
// 				Hn:        input.Hn,
// 				FirstName: input.FirstName,
// 				LastName:  input.LastName,
// 				Verified:  input.Verified,
// 			},
// 		).Return(1, mockErr).Once()

// 		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(rawInput))
// 		recorder := httptest.NewRecorder()
// 		_, router := gin.CreateTestContext(recorder)

// 		router.POST("/", webH.CreatePatient)
// 		router.ServeHTTP(recorder, req)

// 		assert.Equal(t, 500, recorder.Code)
// 	})
// }
func TestUpdatePatient(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		input := model.UpdatePatientRequest{
			Hn:        "testhn",
			FirstName: "fn",
			LastName:  "ln",
			Verified:  false,
		}
		patient := model.Patient{
			ID:        1,
			FirstName: "oldfn",
			LastName:  "oldln",
			Verified:  false,
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().GetPatientById(1).Return(patient, nil).Once()
		repo.EXPECT().UpdatePatient(
			model.Patient{
				ID:        1,
				Hn:        input.Hn,
				FirstName: input.FirstName,
				LastName:  input.LastName,
				Verified:  input.Verified,
			},
		).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPut, "/1", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT("/:id", webH.UpdatePatient)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 200, recorder.Code)
	})
	t.Run("bindingError", func(t *testing.T) {
		input := model.UpdatePatientRequest{
			Hn:        "", // error require
			FirstName: "fn",
			LastName:  "ln",
			Verified:  false,
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		webH := web.WebHandler{}

		req := httptest.NewRequest(http.MethodPut, "/1", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT("/:id", webH.UpdatePatient)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 400, recorder.Code)
	})
	t.Run("atoiError", func(t *testing.T) {
		input := model.UpdatePatientRequest{
			Hn:        "hn",
			FirstName: "fn",
			LastName:  "ln",
			Verified:  false,
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		webH := web.WebHandler{}

		req := httptest.NewRequest(http.MethodPut, "/asd", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT("/:id", webH.UpdatePatient)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 400, recorder.Code)
	})
	t.Run("notFound", func(t *testing.T) {
		input := model.UpdatePatientRequest{
			Hn:        "testhn",
			FirstName: "fn",
			LastName:  "ln",
			Verified:  false,
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		mockErr := fmt.Errorf("wrap : %w", gorm.ErrRecordNotFound)
		repo.EXPECT().GetPatientById(1).Return(model.Patient{}, mockErr).Once()

		req := httptest.NewRequest(http.MethodPut, "/1", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT("/:id", webH.UpdatePatient)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 404, recorder.Code)
	})
	t.Run("getPatientInternalError", func(t *testing.T) {
		input := model.UpdatePatientRequest{
			Hn:        "testhn",
			FirstName: "fn",
			LastName:  "ln",
			Verified:  false,
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		mockErr := fmt.Errorf("wrap : %w", errors.New("some internal error"))
		repo.EXPECT().GetPatientById(1).Return(model.Patient{}, mockErr).Once()

		req := httptest.NewRequest(http.MethodPut, "/1", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT("/:id", webH.UpdatePatient)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("duplicateHn", func(t *testing.T) {
		input := model.UpdatePatientRequest{
			Hn:        "testhn",
			FirstName: "fn",
			LastName:  "ln",
			Verified:  false,
		}
		patient := model.Patient{
			ID:        1,
			FirstName: "oldfn",
			LastName:  "oldln",
			Verified:  false,
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().GetPatientById(1).Return(patient, nil).Once()
		mockErr := fmt.Errorf("wrap : %w", repository.ErrDuplicateEntry)
		repo.EXPECT().UpdatePatient(
			model.Patient{
				ID:        1,
				Hn:        input.Hn,
				FirstName: input.FirstName,
				LastName:  input.LastName,
				Verified:  input.Verified,
			},
		).Return(mockErr).Once()

		req := httptest.NewRequest(http.MethodPut, "/1", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT("/:id", webH.UpdatePatient)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 409, recorder.Code)
	})
	t.Run("updateInternalError", func(t *testing.T) {
		input := model.UpdatePatientRequest{
			Hn:        "testhn",
			FirstName: "fn",
			LastName:  "ln",
			Verified:  false,
		}
		patient := model.Patient{
			ID:        1,
			FirstName: "oldfn",
			LastName:  "oldln",
			Verified:  false,
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().GetPatientById(1).Return(patient, nil).Once()
		mockErr := fmt.Errorf("wrap : %w", errors.New("some internal error"))
		repo.EXPECT().UpdatePatient(
			model.Patient{
				ID:        1,
				Hn:        input.Hn,
				FirstName: input.FirstName,
				LastName:  input.LastName,
				Verified:  input.Verified,
			},
		).Return(mockErr).Once()

		req := httptest.NewRequest(http.MethodPut, "/1", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT("/:id", webH.UpdatePatient)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
}
func TestDeletePatient(t *testing.T) {
	t.Run("atoiError", func(t *testing.T) {
		// setup mock
		webH := web.WebHandler{}

		req := httptest.NewRequest(http.MethodDelete, "/asd", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.DELETE("/:id", webH.DeletePatient)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 400, recorder.Code)
	})
	t.Run("internalError", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().DeletePatientById(1).Return(errors.New("some internal error"))

		req := httptest.NewRequest(http.MethodDelete, "/1", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.DELETE("/:id", webH.DeletePatient)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("success", func(t *testing.T) {
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().DeletePatientById(1).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/1", nil)
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.DELETE("/:id", webH.DeletePatient)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 204, recorder.Code)
	})
}
func TestUpdatePatientVaccineHistory(t *testing.T) {
	t.Run("bindingError", func(t *testing.T) {
		input := model.UpdateVaccineHistoryRequest{
			Data: []model.VaccineHistory{{Id: "test", VaccineName: "hello", VaccineAt: 0}}, // error require
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		webH := web.WebHandler{}

		req := httptest.NewRequest(http.MethodPut, "/1", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT("/:id", webH.UpdatePatientVaccineHistory)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 400, recorder.Code)
	})
	t.Run("atoiError", func(t *testing.T) {
		input := model.UpdateVaccineHistoryRequest{
			Data: []model.VaccineHistory{{Id: "test", VaccineName: "hello", VaccineAt: 123}},
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		webH := web.WebHandler{}

		req := httptest.NewRequest(http.MethodPut, "/asd", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT("/:id", webH.UpdatePatientVaccineHistory)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 400, recorder.Code)
	})
	t.Run("atoiError", func(t *testing.T) {
		input := model.UpdateVaccineHistoryRequest{
			Data: []model.VaccineHistory{{Id: "test", VaccineName: "hello", VaccineAt: 123}},
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		mockErr := fmt.Errorf("wrap : %w", gorm.ErrRecordNotFound)
		repo.EXPECT().GetPatientById(1).Return(model.Patient{}, mockErr).Once()

		req := httptest.NewRequest(http.MethodPut, "/1", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT("/:id", webH.UpdatePatientVaccineHistory)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 404, recorder.Code)
	})
	t.Run("getPatientInternalError", func(t *testing.T) {
		input := model.UpdateVaccineHistoryRequest{
			Data: []model.VaccineHistory{{Id: "test", VaccineName: "hello", VaccineAt: 123}},
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		mockErr := fmt.Errorf("wrap : %w", errors.New("some internal error"))
		repo.EXPECT().GetPatientById(1).Return(model.Patient{}, mockErr).Once()

		req := httptest.NewRequest(http.MethodPut, "/1", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT("/:id", webH.UpdatePatientVaccineHistory)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("updateInternalError", func(t *testing.T) {
		input := model.UpdateVaccineHistoryRequest{
			Data: []model.VaccineHistory{{Id: "test", VaccineName: "hello", VaccineAt: 123}},
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().GetPatientById(1).Return(model.Patient{}, nil).Once()
		repo.EXPECT().UpdatePatientVaccineHistory(1, input.Data).Return(errors.New("some internal error")).Once()

		req := httptest.NewRequest(http.MethodPut, "/1", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT("/:id", webH.UpdatePatientVaccineHistory)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("success", func(t *testing.T) {
		input := model.UpdateVaccineHistoryRequest{
			Data: []model.VaccineHistory{{Id: "test", VaccineName: "hello", VaccineAt: 123}},
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().GetPatientById(1).Return(model.Patient{}, nil).Once()
		repo.EXPECT().UpdatePatientVaccineHistory(1, input.Data).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPut, "/1", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT("/:id", webH.UpdatePatientVaccineHistory)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 200, recorder.Code)
	})
}
func TestUpdatePatientMedicine(t *testing.T) {
	t.Run("bindingError", func(t *testing.T) {
		input := model.UpdateMedicineRequest{
			Data: []model.Medicine{{Id: "", MedicineName: "hello"}}, // error require
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		webH := web.WebHandler{}

		req := httptest.NewRequest(http.MethodPut, "/1", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT("/:id", webH.UpdatePatientMedicine)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 400, recorder.Code)
	})
	t.Run("atoiError", func(t *testing.T) {
		input := model.UpdateMedicineRequest{
			Data: []model.Medicine{{Id: "test", MedicineName: "hello"}}, // error require
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		webH := web.WebHandler{}

		req := httptest.NewRequest(http.MethodPut, "/asd", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT("/:id", webH.UpdatePatientMedicine)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 400, recorder.Code)
	})
	t.Run("notFound", func(t *testing.T) {
		input := model.UpdateMedicineRequest{
			Data: []model.Medicine{{Id: "test", MedicineName: "hello"}}, // error require
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		mockErr := fmt.Errorf("wrap : %w", gorm.ErrRecordNotFound)
		repo.EXPECT().GetPatientById(1).Return(model.Patient{}, mockErr).Once()

		req := httptest.NewRequest(http.MethodPut, "/1", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT("/:id", webH.UpdatePatientMedicine)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 404, recorder.Code)
	})
	t.Run("getPatientInternalError", func(t *testing.T) {
		input := model.UpdateMedicineRequest{
			Data: []model.Medicine{{Id: "test", MedicineName: "hello"}}, // error require
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		mockErr := fmt.Errorf("wrap : %w", errors.New("some internal error"))
		repo.EXPECT().GetPatientById(1).Return(model.Patient{}, mockErr).Once()

		req := httptest.NewRequest(http.MethodPut, "/1", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT("/:id", webH.UpdatePatientMedicine)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("updateInternalError", func(t *testing.T) {
		input := model.UpdateMedicineRequest{
			Data: []model.Medicine{{Id: "test", MedicineName: "hello"}}, // error require
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().GetPatientById(1).Return(model.Patient{}, nil).Once()
		repo.EXPECT().UpdatePatientMedicine(1, input.Data).Return(errors.New("some internal error")).Once()

		req := httptest.NewRequest(http.MethodPut, "/1", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT("/:id", webH.UpdatePatientMedicine)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 500, recorder.Code)
	})
	t.Run("success", func(t *testing.T) {
		input := model.UpdateMedicineRequest{
			Data: []model.Medicine{{Id: "test", MedicineName: "hello"}}, // error require
		}
		rawInput, err := json.Marshal(&input)
		assert.NoError(t, err)
		// setup mock
		repo := repository.NewMockRepo(t)
		webH := web.WebHandler{Repo: repo}

		repo.EXPECT().GetPatientById(1).Return(model.Patient{}, nil).Once()
		repo.EXPECT().UpdatePatientMedicine(1, input.Data).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPut, "/1", bytes.NewReader(rawInput))
		recorder := httptest.NewRecorder()
		_, router := gin.CreateTestContext(recorder)

		router.PUT("/:id", webH.UpdatePatientMedicine)
		router.ServeHTTP(recorder, req)

		assert.Equal(t, 200, recorder.Code)
	})
}
