package web

import (
	"errors"
	"net/http"

	"github.com/PhasitWo/duchenne-server/model"
	"github.com/PhasitWo/duchenne-server/repository"
	"github.com/gin-gonic/gin"
)

func (w *WebHandler) GetAllDoctor(c *gin.Context) {
	doctors, err := w.Repo.GetAllDoctor()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, doctors)
}

type doctorInput struct {
	FirstName  string     `json:"firstName" binding:"required"`
	MiddleName *string    `json:"middleName"`
	LastName   string     `json:"lastName" binding:"required"`
	Username   string     `json:"username" binding:"required"`
	Password   string     `json:"password" binding:"required"`
	Role       model.Role `json:"role" binding:"required"`
}

func (w *WebHandler) CreateDoctor(c *gin.Context) {
	var input doctorInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	insertedId, err := w.Repo.CreateDoctor(model.Doctor{
		FirstName:  input.FirstName,
		MiddleName: input.MiddleName,
		LastName:   input.LastName,
		Username:   input.Username,
		Password:   input.Password,
		Role:       input.Role,
	})
	if err != nil {
		if errors.Unwrap(err) == repository.ErrDuplicateEntry {
			c.Status(http.StatusConflict)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id" : insertedId})
}
