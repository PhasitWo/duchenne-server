package middleware

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)


var redisLogger = log.New(os.Stdout, "[REDIS] ", log.LstdFlags)

type RedisClient struct {
	Client *redis.Client
}

// wrapper around original writer
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
// UNUSED
func (rdc *RedisClient) RedisPersonalizedCacheMiddleware(c *gin.Context) {
	// Setting and Getting cache will hit the performance if sever can't connect to redis
	// get patientId from auth header
	i, exists := c.Get("patientId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no 'patientId' from auth middleware"})
		return
	}
	patientId := i.(int)
	key := "$" + strconv.Itoa(patientId) + c.Request.RequestURI
	// setup writer wrapper
	w := responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
	c.Writer = &w
	if c.Request.Method == "POST" {
		c.Next()
		if c.Writer.Status() == http.StatusCreated {
			rdc.Client.Del(context.Background(), key)
			redisLogger.Println("Invalidate ", key)
		}
		return
	}
	if c.Request.Method == "DELETE" {
		c.Next()
		if c.Writer.Status() == http.StatusNoContent {
			i := strings.LastIndex(key, "/")
			newKey := key[0:i]
			rdc.Client.Del(context.Background(), newKey)
			redisLogger.Println("Invalidate ", newKey)
		}
		return
	}
	var value []byte
	err := rdc.Client.Get(context.Background(), key).Scan(&value)
	if err == nil {
		c.Data(http.StatusOK, "application/json", value)
		redisLogger.Println("Cache hit ", key)
		c.Abort()
		return
	}
	redisLogger.Println("Cache miss")
	c.Next()
	if c.Writer.Status() == http.StatusOK {
		rdc.Client.Set(context.Background(), key, w.body.Bytes(), time.Minute*5).Err()
	}
}

func (rdc *RedisClient) UseRedisMiddleware(handler ...gin.HandlerFunc) []gin.HandlerFunc {
	if rdc == nil {
		return handler
	}
	res := []gin.HandlerFunc{rdc.RedisPersonalizedCacheMiddleware}
	res = append(res, handler...)
	return res
}