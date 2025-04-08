package web

import (
	"errors"
	"net/http"

	"github.com/PhasitWo/duchenne-server/model"
	"github.com/PhasitWo/duchenne-server/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (w *WebHandler) GetProfile(c *gin.Context) {
	id, exists := c.Get("doctorId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no 'doctorId' from auth middleware"})
		return
	}
	doctorId := id.(int)
	// fetch doctor from database
	d, err := w.Repo.GetDoctorById(doctorId)
	if err != nil {
		if errors.Unwrap(err) == gorm.ErrRecordNotFound { // no rows found
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, d)
}

func (w *WebHandler) UpdateProfile(c *gin.Context) {
	i, exists := c.Get("doctorId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no 'doctorId' from auth middleware"})
		return
	}
	id := i.(int)
	r, exists := c.Get("doctorRole")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no 'doctorRole' from auth middleware"})
		return
	}
	role := r.(model.Role)
	// input
	var input model.UpdateProfileRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := w.Repo.UpdateDoctor(
		model.Doctor{
			ID:         id,
			FirstName:  input.FirstName,
			MiddleName: input.MiddleName,
			LastName:   input.LastName,
			Username:   input.Username,
			Password:   input.Password,
			Specialist: input.Specialist,
			Role:       role,
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
