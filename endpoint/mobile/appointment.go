package mobile

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/PhasitWo/duchenne-server/repository"
	"github.com/gin-gonic/gin"
)

// "fmt"
// "net/http"

// "github.com/PhasitWo/duchenne-server/auth"
// "github.com/PhasitWo/duchenne-server/model"
// "github.com/PhasitWo/duchenne-server/repository"

func (m *mobileHandler) GetAllPatientAppointment(c *gin.Context) {
	i, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no 'user_id' from auth middleware"})
		return
	}
	id := i.(int)
	aps, err := m.repo.GetAllAppointment(id, repository.PATIENTID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, aps)
}

func (m *mobileHandler) GetPatientAppointment(c *gin.Context) {
	param := c.Param("id")
	i, err := strconv.ParseInt(param, 10, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't parse 'id' param to int"})
		return
	}
	id := int(i)
	ap, err := m.repo.GetAppointment(id)
	if err != nil {
		if errors.Unwrap(err) == sql.ErrNoRows { // no rows found
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ap)
}

type appointmentInput struct {
	Id        int `json:"id" binding:"required"`
	CreateAt  int `json:"createAt" binding:"required"`
	Date      int `json:"date" binding:"required"`
	PatientId int `json:"patient" binding:"required"`
	DoctorId  int `json:"doctor" binding:"required"`
}

func (m *mobileHandler) CreateAppointment(c *gin.Context) {
	var input appointmentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//TODO validate input.date
	//TODO schedule notification
	err := m.repo.CreateAppointment(input.CreateAt, input.Date, input.PatientId, input.DoctorId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.Status(http.StatusCreated)
}

func (m *mobileHandler) Test(c *gin.Context) {
	res, _ := m.repo.GetAllAppointment(1, repository.NONE)
	c.JSON(http.StatusOK, res)
}
