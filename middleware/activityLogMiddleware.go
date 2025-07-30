package middleware

import (
	"bytes"
	"encoding/json"
	"log"
	"os"

	"github.com/PhasitWo/duchenne-server/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var logger = log.New(os.Stdout, "[ACTIVITY_LOG] ", log.LstdFlags)

type ActivityLogMiddleware struct {
	DB *gorm.DB
}

func InitActivityLogMiddleware(db *gorm.DB) *ActivityLogMiddleware {
	return &ActivityLogMiddleware{db}
}

func (a *ActivityLogMiddleware) ActivityLog(c *gin.Context) {
	method := c.Request.Method
	if method == "POST" || method == "PUT" || method == "DELETE" {
		// setup writer wrapper
		w := responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = &w
		c.Next()
		// after handler
		data := w.body.Bytes()
		status := w.Status()
		clm, exists := c.Get("claims")
		if !exists {
			logger.Println("no claims from auth middlewares")
			c.Next()
			return
		}
		claims, err := json.Marshal(clm)
		if err != nil {
			logger.Println("error marshaling claims")
			c.Next()
			return
		}
		l := &model.ActivityLog{
			Claims:     claims,
			Method:     method,
			RequestURL: c.Request.URL.Path,
			Data:       data,
			Status:     status,
		}
		go func(log *model.ActivityLog) {
			err := a.DB.Create(l).Error
			if err != nil {
				logger.Println("error inserting data to the database")
				return
			}
		}(l)
	}
	c.Next()
}
