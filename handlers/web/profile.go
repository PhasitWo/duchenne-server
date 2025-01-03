package web

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/PhasitWo/duchenne-server/model"
	"github.com/gin-gonic/gin"
)

func (w *WebHandler) GetProfile(c *gin.Context) {
	id, exists := c.Get("doctorId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no 'doctorId' from auth middleware"})
		return
	}
	doctorId := id.(int)
	// fetch patient from database
	d, err := w.Repo.GetDoctorById(doctorId)
	if err != nil {
		if errors.Unwrap(err) == sql.ErrNoRows { // no rows found
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, d)
}

type profile struct {
	FirstName  string  `json:"firstName" binding:"required"`
	MiddleName *string `json:"middleName"`
	LastName   string  `json:"lastName" binding:"required"`
	Username   string  `json:"username" binding:"required"`
	Password   string  `json:"password" binding:"required"`
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
	var input profile
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := w.Repo.UpdateDoctor(
		model.Doctor{
			Id:         id,
			FirstName:  input.FirstName,
			MiddleName: input.MiddleName,
			LastName:   input.LastName,
			Username:   input.LastName,
			Password:   input.Password,
			Role:       role,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
