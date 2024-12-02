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
	i, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no 'userId' from auth middleware"})
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
	c.JSON(http.StatusOK, ap)
}

type appointmentInput struct {
	CreateAt int `json:"createAt" binding:"required"`
	Date     int `json:"date" binding:"required"`
	DoctorId int `json:"doctorId" binding:"required"`
}

func (m *mobileHandler) CreateAppointment(c *gin.Context) {
	// get patientId from auth header
	i, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no 'userId' from auth middleware"})
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
	tx, err := m.repo.CreateAppointment(input.CreateAt, input.Date, patientId, input.DoctorId)
	defer func() {
		if err := tx.Rollback(); !errors.Is(err, sql.ErrTxDone) {
			c.JSON(http.StatusInternalServerError, gin.H{"tx": "Can't rollback"})
		}
	}()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	/*
		TODO schedule notification
	*/
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"tx": "Can't commit"})
	}
	c.Status(http.StatusCreated)
}

func (m *mobileHandler) DeleteAppointment(c *gin.Context) {
	id := c.Param("id")
	tx, err := m.repo.DeleteAppointment(id)
	defer func() {
		if err := tx.Rollback(); !errors.Is(err, sql.ErrTxDone) {
			c.JSON(http.StatusInternalServerError, gin.H{"tx": "Can't rollback"})
		}
	}()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// TODO cancel the notification, then commit
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"tx": "Can't commit"})
	}
	c.Status(http.StatusNoContent)
}

func (m *mobileHandler) Test(c *gin.Context) {
	res, _ := m.repo.GetAllAppointment(1, repository.NONE)
	c.JSON(http.StatusOK, res)
}
