package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
mobile endpoints -> /mobile/
AUTH
POST /login	-> authenticate user
POST /verify-user-info -> verify user info in singup process
POST /set-password
POST /verify-code

PROFILE
GET /profle -> return user profile data

APPOINTMENT
GET /appointment  -> return maximum 20 of patient's appointments
POST /appointment -> create new appointment

ASK
GET /ask -> return patient's question history
GET /ask/:id -> return patient's question and doctor's answer
POST /ask -> create new question

*/

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
