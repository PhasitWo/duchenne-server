package mobile

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (m *MobileHandler) GetProfile(c *gin.Context) {
	id, exists := c.Get("patientId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no 'patientId' from auth middleware"})
		return
	}
	// fetch patient from database
	p, err := m.Repo.GetPatient(id)
	if err != nil {
		if errors.Unwrap(err) == sql.ErrNoRows { // no rows found
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}
