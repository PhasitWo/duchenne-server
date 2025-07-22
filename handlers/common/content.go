package common

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PhasitWo/duchenne-server/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (c *CommonHandler) GetAllContent(ctx *gin.Context) {
	limit := 9999
	offset := 0
	criteria := []repository.Criteria{}
	var err error
	// get url query param
	if l, exist := ctx.GetQuery("limit"); exist {
		limit, err = strconv.Atoi(l)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "cannot parse limit value"})
			return
		}
	}
	if of, exist := ctx.GetQuery("offset"); exist {
		offset, err = strconv.Atoi(of)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "cannot parse offset value"})
			return
		}
	}
	if _, exist := ctx.GetQuery("isPublished"); exist {
		criteria = append(criteria, repository.Criteria{QueryCriteria: repository.IS_PUBLISHED, Value: true})
	}
	if _, exist := ctx.GetQuery("notPublished"); exist {
		criteria = append(criteria, repository.Criteria{QueryCriteria: repository.IS_PUBLISHED, Value: false})
	}
	// query
	contents, err := c.Repo.GetAllContent(limit, offset, criteria...)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, contents)
}

func (c *CommonHandler) GetOneContent(ctx *gin.Context) {
	id := ctx.Param("id")
	content, err := c.Repo.GetContent(id)
	if err != nil {
		if errors.Unwrap(err) == gorm.ErrRecordNotFound { // no rows found
			ctx.Status(http.StatusNotFound)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, content)
}

func (c *CommonHandler) UploadImage(ctx *gin.Context) {
	file, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// check content type
	contentType := file.Header.Get("Content-Type")
	if !strings.Contains(contentType, "image") {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid content type: %v", contentType)})
		return
	}
	publicURL, err := c.CloudStorageService.UploadImage(file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("error uploading file: %v", err.Error())})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"publicURL": publicURL})
}
