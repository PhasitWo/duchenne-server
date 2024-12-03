package mobile

import (
	"database/sql"
	"errors"

	// "github.com/PhasitWo/duchenne-server/repository"
	"net/http"

	"github.com/PhasitWo/duchenne-server/repository"
	"github.com/gin-gonic/gin"
)

func (m *mobileHandler) GetAllPatientQuestion(c *gin.Context) {
	i, exists := c.Get("patientId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no 'patientId' from auth middleware"})
		return
	}
	id := i.(int)
	qs, err := m.repo.GetAllQuestion(id, repository.PATIENTID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, qs)
}

func (m *mobileHandler) GetQuestion(c *gin.Context) {
	i, exists := c.Get("patientId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no 'patientId' from auth middleware"})
		return
	}
	patientId := i.(int)
	id := c.Param("id")
	q, err := m.repo.GetQuestion(id)
	if err != nil {
		if errors.Unwrap(err) == sql.ErrNoRows { // no rows found
			c.Status(http.StatusNotFound)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// check if this question belongs to the patient
	if patientId != q.Patient.Id {
		c.Status(http.StatusUnauthorized)
		return
	}
	c.JSON(http.StatusOK, q)
}
