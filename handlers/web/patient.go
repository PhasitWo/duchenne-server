package web

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/PhasitWo/duchenne-server/model"
	"github.com/PhasitWo/duchenne-server/repository"
	"github.com/gin-gonic/gin"
)

func (w *WebHandler) GetPatient(c *gin.Context) {
	id := c.Param("id")
	patient, err := w.Repo.GetPatientById(id)
	if err != nil {
		if errors.Unwrap(err) == sql.ErrNoRows { // no rows found
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, patient)
}

func (w *WebHandler) GetAllPatient(c *gin.Context) {
	patients, err := w.Repo.GetAllPatient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, patients)
}

type patientInput struct {
	Hn         string  `json:"hn" binding:"required"`
	FirstName  string  `json:"firstName" binding:"required"`
	MiddleName *string `json:"middleName"`
	LastName   string  `json:"lastName" binding:"required"`
	Email      *string `json:"email"`
	Phone      *string `json:"phone"`
	Verified   bool    `json:"verified"`
}

func (w *WebHandler) CreatePatient(c *gin.Context) {
	var input patientInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	insertedId, err := w.Repo.CreatePatient(model.Patient{
		Hn:         input.Hn,
		FirstName:  input.FirstName,
		MiddleName: input.MiddleName,
		LastName:   input.LastName,
		Email:      input.Email,
		Phone:      input.Phone,
		Verified:   input.Verified,
	})
	if err != nil {
		if errors.Unwrap(err) == repository.ErrDuplicateEntry {
			c.JSON(http.StatusConflict, gin.H{"error": "duplicate hn"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": insertedId})
}

// func (w *WebHandler) UpdateDoctor(c *gin.Context) {
// 	var input doctorInput
// 	if err := c.ShouldBindJSON(&input); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	i := c.Param("id")
// 	id, err := strconv.Atoi(i)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	err = w.Repo.UpdateDoctor(model.Doctor{
// 		Id:         id,
// 		FirstName:  input.FirstName,
// 		MiddleName: input.MiddleName,
// 		LastName:   input.LastName,
// 		Username:   input.Username,
// 		Password:   input.Password,
// 		Role:       input.Role,
// 	})
// 	if err != nil {
// 		if errors.Unwrap(err) == repository.ErrDuplicateEntry {
// 			c.Status(http.StatusConflict)
// 			return
// 		}
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.Status(http.StatusOK)
// }

// // DeleteDoctor is idempotent
// func (w *WebHandler) DeleteDoctor(c *gin.Context) {
// 	i := c.Param("id")
// 	id, err := strconv.Atoi(i)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	err = w.Repo.DeleteDoctorById(id)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.Status(http.StatusNoContent)
// }
