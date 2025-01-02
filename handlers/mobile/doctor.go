package mobile

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (m *MobileHandler) GetAllDoctor(c *gin.Context) {
	doctors, err := m.Repo.GetAllDoctor()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, doctors)
}
