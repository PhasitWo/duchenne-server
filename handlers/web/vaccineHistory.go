package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type vaccineHistory struct {
	Id              string `json:"id" binding:"required"`
	VaccineName     string `json:"vaccineName" binding:"required"`
	VaccineLocation string `json:"vaccineLocation" binding:"required"`
	VaccineAt       int    `json:"vaccineAt" binding:"required"`
	Description     string `json:"description"`
}

type vaccineHistoryInput struct {
	Data []vaccineHistory `json:"data" binding:"required,dive"`
}

func (w *WebHandler) Test(c *gin.Context) {
	// input
	var input vaccineHistoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("%v", input)
	c.Status(http.StatusOK)
}
