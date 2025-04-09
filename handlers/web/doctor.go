package web

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/PhasitWo/duchenne-server/model"
	"github.com/PhasitWo/duchenne-server/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (w *WebHandler) GetDoctor(c *gin.Context) {
	id := c.Param("id")
	doctor, err := w.Repo.GetDoctorById(id)
	if err != nil {
		if errors.Unwrap(err) == gorm.ErrRecordNotFound { // no rows found
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, doctor)
}

func (w *WebHandler) GetAllDoctor(c *gin.Context) {
	doctors, err := w.Repo.GetAllDoctor()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, doctors)
}


func (w *WebHandler) CreateDoctor(c *gin.Context) {
	var input model.CreateDoctorRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.Role != model.ADMIN && input.Role != model.ROOT && input.Role != model.USER {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role value"})
		return
	}
	insertedId, err := w.Repo.CreateDoctor(model.Doctor{
		FirstName:  input.FirstName,
		MiddleName: input.MiddleName,
		LastName:   input.LastName,
		Username:   input.Username,
		Password:   input.Password,
		Role:       input.Role,
		Specialist: input.Specialist,
	})
	if err != nil {
		if errors.Unwrap(err) == repository.ErrDuplicateEntry {
			c.JSON(http.StatusConflict, gin.H{"error": "duplicate username"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": insertedId})
}

func (w *WebHandler) UpdateDoctor(c *gin.Context) {
	var input model.CreateDoctorRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.Role != model.ADMIN && input.Role != model.ROOT && input.Role != model.USER {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role value"})
		return
	}
	i := c.Param("id")
	id, err := strconv.Atoi(i)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err = w.Repo.GetDoctorById(id) // check if this id exist
	if err != nil {
		if errors.Unwrap(err) == gorm.ErrRecordNotFound { // no rows found
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = w.Repo.UpdateDoctor(model.Doctor{
		ID:         id,
		FirstName:  input.FirstName,
		MiddleName: input.MiddleName,
		LastName:   input.LastName,
		Username:   input.Username,
		Password:   input.Password,
		Role:       input.Role,
		Specialist: input.Specialist,
	})
	if err != nil {
		if errors.Unwrap(err) == repository.ErrDuplicateEntry {
			c.JSON(http.StatusConflict, gin.H{"error": "duplicate username"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// DeleteDoctor is idempotent
func (w *WebHandler) DeleteDoctor(c *gin.Context) {
	i := c.Param("id")
	id, err := strconv.Atoi(i)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = w.Repo.DeleteDoctorById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
