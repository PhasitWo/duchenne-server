package mobile

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/PhasitWo/duchenne-server/repository"
	"github.com/gin-gonic/gin"
)

// "fmt"
// "net/http"

// "github.com/PhasitWo/duchenne-server/auth"
// "github.com/PhasitWo/duchenne-server/model"
// "github.com/PhasitWo/duchenne-server/repository"

func (m *MobileHandler) GetAllPatientAppointment(c *gin.Context) {
	i, exists := c.Get("patientId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no 'patientId' from auth middleware"})
		return
	}
	id := i.(int)
	aps, err := m.Repo.GetAllAppointment(id, repository.PATIENTID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, aps)
}

func (m *MobileHandler) GetAppointment(c *gin.Context) {
	i, exists := c.Get("patientId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no 'patientId' from auth middleware"})
		return
	}
	patientId := i.(int)
	id := c.Param("id")
	ap, err := m.Repo.GetAppointment(id)
	if err != nil {
		if errors.Unwrap(err) == sql.ErrNoRows { // no rows found
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// check if this appointment belongs to the patient
	if patientId != ap.Patient.Id {
		c.Status(http.StatusUnauthorized)
		return
	}
	c.JSON(http.StatusOK, ap)
}

type appointmentInput struct {
	Date     int `json:"date" binding:"required"`
	DoctorId int `json:"doctorId" binding:"required"`
}

func (m *MobileHandler) CreateAppointment(c *gin.Context) {
	// get patientId from auth header
	i, exists := c.Get("patientId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no 'patientId' from auth middleware"})
		return
	}
	patientId := i.(int)
	// binding request body
	var input appointmentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// validate input.date
	now := int(time.Now().Add(3 * time.Minute).Unix())
	if input.Date < now {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "'Date' is before current time"})
		return
	}
	// create new appointment
	insertedId, err := m.Repo.CreateAppointment(int(time.Now().Unix()), input.Date, patientId, input.DoctorId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": insertedId})
}

func (m *MobileHandler) DeleteAppointment(c *gin.Context) {
	// prepare param from url and auth middleware
	i, exists := c.Get("patientId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no 'patientId' from auth middleware"})
		return
	}
	patientId := i.(int)
	id := c.Param("id")
	// check if this appointment belongs to the patient
	ap, err := m.Repo.GetAppointment(id)
	if err != nil {
		if errors.Unwrap(err) == sql.ErrNoRows { // no rows found
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if patientId != ap.Patient.Id {
		c.Status(http.StatusUnauthorized)
		return
	}
	// delete appointment
	err = m.Repo.DeleteAppointment(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (m *MobileHandler) Test(c *gin.Context) {
	// _, tx, err := m.Repo.CreateDevice(model.Device{Id: -1, LoginAt: 666666, DeviceName: "cxd phone", ExpoToken: "hhaha", PatientId: 1})
	// defer tx.Rollback()
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }
	// tx.Commit()
	c.Status(http.StatusOK)
}
