package common

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


func (c *CommonHandler) GetConsentById(ctx *gin.Context) {
	id := ctx.Param("id")
	consent, err := c.Repo.GetConsentById(id)
	if err != nil {
		if errors.Unwrap(err) == gorm.ErrRecordNotFound { // no rows found
			ctx.Status(http.StatusNotFound)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, consent)
}

func (c *CommonHandler) GetConsentBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")
	consent, err := c.Repo.GetConsentBySlug(slug)
	if err != nil {
		if errors.Unwrap(err) == gorm.ErrRecordNotFound { // no rows found
			ctx.Status(http.StatusNotFound)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, consent)
}

