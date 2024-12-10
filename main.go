package main

import (
	// "net/http"
	"crypto/tls"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/PhasitWo/duchenne-server/config"
	"github.com/PhasitWo/duchenne-server/endpoint/mobile"
	"github.com/PhasitWo/duchenne-server/middleware"
	"github.com/PhasitWo/duchenne-server/notification"
	"github.com/robfig/cron"

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
var mainLogger = log.New(os.Stdout, "[MAIN] ", log.LstdFlags)
func main() {
	db := setupDB()
	defer db.Close()
	r := setupRouter()
	m := mobile.Init(db)
	attachHandler(r, m)
	c := InitCronScheduler(db)
	defer c.Stop()
	r.Run() // listen and serve on 0.0.0.0:8080
}

func attachHandler(r *gin.Engine, m *mobile.MobileHandler) {
	mobile := r.Group("/mobile")
	mobile.POST("/testnoti", func(c *gin.Context) {
		notification.MockupScheduleNotifications(m.DBConn, notification.SendRequest)
	})
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
}

func setupDB() *sql.DB {
	// read config
	config.LoadConfig()
	// open db connection
	mysql.RegisterTLSConfig("tidb", &tls.Config{
		MinVersion: tls.VersionTLS12,
		ServerName: "gateway01.ap-southeast-1.prod.aws.tidbcloud.com",
	})

	db, err := sql.Open("mysql", config.AppConfig.DATABASE_DSN)
	if err != nil {
		mainLogger.Panicf("Can't open connection to database : %v", err.Error())
	}
	mainLogger.Println("Connected to database")
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	// ping db
	if err = db.Ping(); err != nil {
		mainLogger.Panic("Can't connect to database")
	}
	return db
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	return r
}

func InitCronScheduler(db *sql.DB) *cron.Cron {
	c := cron.New()
	// everyday on 10.00 -> spec : "00 10 * * *"
	c.AddFunc("00 10 * * *", func() {
		mainLogger.Println("Executing Push Notifications..")
		notification.MockupScheduleNotifications(db, notification.MockSendRequest)
	})
	c.Start()
	mainLogger.Println("Cron scheduler initialized")
	return c
}
