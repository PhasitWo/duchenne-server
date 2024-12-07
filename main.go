package main

import (
	// "net/http"
	"crypto/tls"
	"database/sql"
	"fmt"
	"time"

	"github.com/PhasitWo/duchenne-server/config"
	"github.com/PhasitWo/duchenne-server/endpoint/mobile"
	"github.com/PhasitWo/duchenne-server/middleware"
	"github.com/PhasitWo/duchenne-server/notification"

	"github.com/gin-gonic/gin"

	// "github.com/PhasitWo/duchenne-server/repository"

	"github.com/go-sql-driver/mysql"
)

/*
NOTIFICATION
TODO schedule cronjob to run everyday at xx:xx
TODO set condition to filter appointments -> in range X days  

web app
POST /question/:id/answer
*/

func main() {
	// read config
	config.LoadConfig()
	// open db connection
	mysql.RegisterTLSConfig("tidb", &tls.Config{
		MinVersion: tls.VersionTLS12,
		ServerName: "gateway01.ap-southeast-1.prod.aws.tidbcloud.com",
	})

	db, err := sql.Open("mysql", config.AppConfig.DATABASE_DSN)
	if err != nil {
		panic(fmt.Sprintf("Can't open connection to database : %v", err.Error()))
	}
	fmt.Println("Connected to database")
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
	mobile.POST("/testnoti", notification.TestPushNotification(db))
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
			mobileProtected.GET("/doctor", m.GetAllDoctor)
		}
	}
	// scheduleNotifications(db)
	r.Run() // listen and serve on 0.0.0.0:8080
}
