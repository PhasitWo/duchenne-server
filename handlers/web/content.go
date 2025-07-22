package web

import (
	"net/http"
	"strconv"

	"github.com/PhasitWo/duchenne-server/model"
	"github.com/gin-gonic/gin"
)

func (w *WebHandler) CreateContent(c *gin.Context) {
	// binding request body
	var input model.CreateContentRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	insertedID, err := w.Repo.CreateContent(model.Content{
		Title:         input.Title,
		Body:          input.Body,
		IsPublished:   input.IsPublished,
		Order:         input.Order,
		CoverImageURL: input.CoverImageURL,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": insertedID})
}

func (w *WebHandler) UpdateContent(c *gin.Context) {
	i := c.Param("id")
	id, err := strconv.Atoi(i)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// binding request body
	var input model.CreateContentRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = w.Repo.UpdateContent(model.Content{
		ID:            id,
		Title:         input.Title,
		Body:          input.Body,
		IsPublished:   input.IsPublished,
		Order:         input.Order,
		CoverImageURL: input.CoverImageURL,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (w *WebHandler) DeleteContent(c *gin.Context) {
	i := c.Param("id")
	id, err := strconv.Atoi(i)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// delete appointment
	err = w.Repo.DeleteContent(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
