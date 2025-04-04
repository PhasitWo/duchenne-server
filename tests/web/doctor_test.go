package web_test

import (
	"encoding/json"
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
}

func TestGetAllDoctor(t *testing.T) {

}

func TestCreateDoctor(t *testing.T) {

}
func TestUpdateDoctor(t *testing.T) {

}
func TestDeleteDoctor(t *testing.T) {

}
