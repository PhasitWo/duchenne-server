package mobile

import (
	"net/http"


	"github.com/gin-gonic/gin"
)

func (m *mobileHandler) GetAllDoctor(c *gin.Context) {
	doctors, err := m.repo.GetAllDoctor()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, doctors)
}