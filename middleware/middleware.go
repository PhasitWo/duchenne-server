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

	"github.com/PhasitWo/duchenne-server/auth"
	"github.com/PhasitWo/duchenne-server/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
)

func MobileAuthMiddleware(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization header provided"})
		c.Abort() // stop the chain of middlewares
		return
	}
	// parse token
	claims := &auth.PatientClaims{PatientId: -1, DeviceId: -1}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWT_KEY), nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	if !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	if claims.PatientId == -1 || claims.DeviceId == -1 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}
	c.Set("patientId", claims.PatientId)
	c.Set("deviceId", claims.DeviceId)
	c.Next()
}

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

func (rdc *RedisClient) RedisPersonalizedCacheMiddleware(c *gin.Context) {
	// FIXME:Setting and Getting cache will hit the performance if sever can't connect to redis
	// get patientId from auth header
	i, exists := c.Get("patientId")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no 'patientId' from auth middleware"})
		return
	}
	patientId := i.(int)
	key := "$" + strconv.Itoa(patientId) + c.Request.URL.Path
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
		rdc.Client.Set(context.Background(), key, w.body.Bytes(), time.Duration(time.Minute*5)).Err()
	}
}
