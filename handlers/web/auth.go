package web

import (
	"errors"
	"net/http"

	"github.com/PhasitWo/duchenne-server/auth"
	"github.com/PhasitWo/duchenne-server/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	// fetch doctor from database
	storedDoctor, err := w.Repo.GetDoctorByUsername(input.Username)
	if err != nil {
		if errors.Unwrap(err) == gorm.ErrRecordNotFound { // no rows found
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
	token, err := auth.GenerateDoctorToken(storedDoctor.ID, storedDoctor.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// set cookie
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(config.Constants.WEB_ACCESS_COOKIE_NAME, token, 60*60*10, "/", config.AppConfig.SERVER_DOMAIN, true, true)
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (w *WebHandler) Logout(c *gin.Context) {
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie(config.Constants.WEB_ACCESS_COOKIE_NAME, "", 1, "/", config.AppConfig.SERVER_DOMAIN, true, true)
	c.Status(http.StatusOK)
}

func (w *WebHandler) GetUserData(c *gin.Context) {
	id, exists := c.Get("doctorId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no 'doctorId' from auth middleware"})
		c.Abort()
		return
	}
	role, exists := c.Get("doctorRole")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no 'doctorRole' from auth middleware"})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"doctorId": id, "role": role})
}
