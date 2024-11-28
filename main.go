package main

import (
	// "net/http"
	"database/sql"
	"fmt"
	"time"

	"github.com/PhasitWo/duchenne-server/config"
	"github.com/PhasitWo/duchenne-server/endpoint/mobile"
	"github.com/PhasitWo/duchenne-server/middleware"

	"github.com/gin-gonic/gin"

	// "github.com/PhasitWo/duchenne-server/repository"
	_ "github.com/go-sql-driver/mysql"
)

/*
mobile endpoints -> /mobile/
AUTH
ok POST /login	-> authenticate user
ok POST /signup -> verify user info

PROFILE
ok GET /profle -> return patient profile data

APPOINTMENT
ok GET /appointment  -> return maximum 20 of patient's appointments
ok GET /appointment/:id -> individual appointment
POST /appointment -> create new appointment
DELETE /appointment/:id

ASK
GET /ask -> return patient's question history
GET /ask/:id -> return patient's question and doctor's answer
POST /ask -> create new question
*/

func main() {
	// read config
	config.LoadConfig()
	// open db connection
	db, err := sql.Open("mysql", config.AppConfig.DATABASE_DSN)
	if err != nil {
		panic(fmt.Sprintf("Can't connect to database : %v", err.Error()))
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	// setup router
	r := gin.Default()
	m := mobile.Init(db)
	mobile := r.Group("/mobile")
	{
		mobileAuth := mobile.Group("/auth")
		{
			mobileAuth.POST("/login", m.Login)
			mobileAuth.POST("/signup", m.Signup)
		}
		mobileProtected := mobile.Group("/api")
		mobileProtected.Use(middleware.AuthMiddleware)
		{
			mobileProtected.GET("/test", m.Test)
			mobileProtected.GET("/profile", m.GetProfile)
			mobileProtected.GET("/appointment", m.GetAllPatientAppointment)
			mobileProtected.GET("/appointment/:id", m.GetPatientAppointment)
		}
	}
	r.Run() // listen and serve on 0.0.0.0:8080
}
