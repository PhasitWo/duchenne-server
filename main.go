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
ok POST /appointment -> create new appointment
ok DELETE /appointment/:id

ok specialize claim type -> PatientClaim, DoctorClaim -> different auth middleware

QUESTION
ok GET /question -> return all patient's question
ok GET /question/:id -> return patient's question
ok POST /question -> create new question
ok DELETE /question/:id

-> new login device will replace oldest login device -> only 'MAX_DEVICE' number of devices will get notifications
TODO change login logic to accept expo_token, device_name, add /logout endpoint

NOTIFICATION package

web app
POST /question/:id/answer
*/

func main() {
	// read config
	config.LoadConfig()
	// open db connection
	db, err := sql.Open("mysql", config.AppConfig.DATABASE_DSN)
	if err != nil {
		panic(fmt.Sprintf("Can't open connection to database : %v", err.Error()))
	}
	defer db.Close()
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	// ping db
	if err = db.Ping(); err != nil {
		panic("Can't connect to database")
	}
	// setup router
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	m := mobile.Init(db)
	mobile := r.Group("/mobile")
	{
		mobileAuth := mobile.Group("/auth")
		{
			mobileAuth.POST("/login", m.Login)
			mobileAuth.POST("/signup", m.Signup)
			mobileAuth.Use(middleware.MobileAuthMiddleware).POST("/logout", m.Logout)
		}
		mobileProtected := mobile.Group("/api")
		mobileProtected.Use(middleware.MobileAuthMiddleware)
		{
			mobileProtected.GET("/test", m.Test)
			mobileProtected.GET("/profile", m.GetProfile)
			mobileProtected.GET("/appointment", m.GetAllPatientAppointment)
			mobileProtected.GET("/appointment/:id", m.GetAppointment)
			mobileProtected.POST("/appointment", m.CreateAppointment)
			mobileProtected.DELETE("/appointment/:id", m.DeleteAppointment)
			mobileProtected.GET("/question", m.GetAllPatientQuestion)
			mobileProtected.GET("/question/:id", m.GetQuestion)
			mobileProtected.POST("/question", m.CreateQuestion)
			mobileProtected.DELETE("/question/:id", m.DeleteQuestion)
		}
	}
	r.Run() // listen and serve on 0.0.0.0:8080
}
