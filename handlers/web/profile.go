package web

import (
	"database/sql"
	"errors"
	"net/http"

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
