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

func (m *mobileHandler) GetAllPatientAppointment(c *gin.Context) {
	i, exists := c.Get("patientId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no 'patientId' from auth middleware"})
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

func (m *mobileHandler) GetAppointment(c *gin.Context) {
	i, exists := c.Get("patientId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no 'patientId' from auth middleware"})
		return
	}
	patientId := i.(int)
	id := c.Param("id")
	ap, err := m.repo.GetAppointment(id)
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

func (m *mobileHandler) CreateAppointment(c *gin.Context) {
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
	now := int(time.Now().Add(5 * time.Minute).Unix())
	if input.Date < now {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "'Date' is before current time"})
		return
	}
	// commit a transaction if scheduling notifications is successful
	tx, err := m.dbConn.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Tx can't rollback"})
		}
	}()
	repoWithTx := repository.New(tx)
	_, err = repoWithTx.CreateAppointment(int(time.Now().Unix()), input.Date, patientId, input.DoctorId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	/*
		TODO schedule notification
	*/
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tx can't commit"})
		return
	}
	c.Status(http.StatusCreated)
}

func (m *mobileHandler) DeleteAppointment(c *gin.Context) {
	// prepare param from url and auth middleware
	i, exists := c.Get("patientId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no 'patientId' from auth middleware"})
		return
	}
	patientId := i.(int)
	id := c.Param("id")
	// check if this appointment belongs to the patient
	ap, err := m.repo.GetAppointment(id)
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
	// commit a transaction if cancelling notifications is successful
	tx, err := m.dbConn.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Tx can't rollback"})
		}
	}()
	repoWithTx := repository.New(tx)
	err = repoWithTx.DeleteAppointment(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	/*

	 TODO cancel the notification, then commit

	*/
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tx can't commit"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (m *mobileHandler) Test(c *gin.Context) {
	// _, tx, err := m.repo.CreateDevice(model.Device{Id: -1, LoginAt: 666666, DeviceName: "cxd phone", ExpoToken: "hhaha", PatientId: 1})
	// defer tx.Rollback()
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }
	// tx.Commit()
	c.Status(http.StatusOK)
}
