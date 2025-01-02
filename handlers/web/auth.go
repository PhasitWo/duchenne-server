package web

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/PhasitWo/duchenne-server/auth"
	"github.com/gin-gonic/gin"
)

type login struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (w *WebHandler) Login(c *gin.Context) {
	var input login
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// fetch patient from database
	storedDoctor, err := w.Repo.GetDoctorByUsername(input.Username)
	if err != nil {
		if errors.Unwrap(err) == sql.ErrNoRows { // no rows found
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// checking
	if storedDoctor.Password != input.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credential"})
		return
	}
	// generate token
	token, err := auth.GenerateDoctorToken(storedDoctor.Id, storedDoctor.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
