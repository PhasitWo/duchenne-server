package web

import (
	"net/http"

	"github.com/PhasitWo/duchenne-server/model"
	"github.com/gin-gonic/gin"
)

func (w *WebHandler) UpsertConsent(c *gin.Context) {
	// binding request body
	var input model.UpsertConsentRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	insertedSlug, err := w.Repo.UpsertConsent(model.Consent{
		Slug:      input.Slug,
		Body:      input.Body,
		DeletedAt: 0,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"slug": insertedSlug})
}

func (w *WebHandler) DeleteConsentById(c *gin.Context) {
	id := c.Param("id")
	err := w.Repo.DeleteConsentById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (w *WebHandler) DeleteConsentBySlug(c *gin.Context) {
	slug := c.Param("slug")
	err := w.Repo.DeleteConsentBySlug(slug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
