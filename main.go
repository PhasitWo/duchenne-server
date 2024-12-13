package main

import (
	// "net/http"
	"context"
	"crypto/tls"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/PhasitWo/duchenne-server/config"
	"github.com/PhasitWo/duchenne-server/endpoint/mobile"
	"github.com/PhasitWo/duchenne-server/middleware"
	"github.com/PhasitWo/duchenne-server/notification"
	"github.com/redis/go-redis/v9"
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
	// Load app config
	config.LoadConfig()
	// Setup database connection
	db := setupDB()
	defer db.Close()
	// Setup redis
	rdc := setupRedisClient()
	// Setup router and handler
	r := setupRouter()
	m := mobile.Init(db)
	attachHandler(r, m, rdc)
	// CRON
	c := InitCronScheduler(db)
	defer c.Stop()
	r.Run() // listen and serve on 0.0.0.0:8080
}

func attachHandler(r *gin.Engine, m *mobile.MobileHandler, rdc *middleware.RedisClient) {
	mobile := r.Group("/mobile")
	mobile.POST("/testnoti", func(c *gin.Context) {
		notification.MockupScheduleNotifications(m.DBConn, notification.SendRequest)
	})
	{
		mobileAuth := mobile.Group("/auth")
		{
			mobileAuth.POST("/login", m.Login)
			mobileAuth.POST("/signup", m.Signup)
			mobileAuth.POST("/logout", middleware.MobileAuthMiddleware, m.Logout)
		}
		mobileProtected := mobile.Group("/api")
		mobileProtected.Use(middleware.MobileAuthMiddleware)
		{
			mobileProtected.GET("/test", m.Test)
			mobileProtected.GET("/profile", m.GetProfile)
			mobileProtected.GET("/appointment", rdc.RedisPersonalizedCacheMiddleware, m.GetAllPatientAppointment)
			mobileProtected.GET("/appointment/:id", m.GetAppointment)
			mobileProtected.POST("/appointment", rdc.RedisPersonalizedCacheMiddleware, m.CreateAppointment)
			mobileProtected.DELETE("/appointment/:id", rdc.RedisPersonalizedCacheMiddleware, m.DeleteAppointment)
			mobileProtected.GET("/question", rdc.RedisPersonalizedCacheMiddleware, m.GetAllPatientQuestion)
			mobileProtected.GET("/question/:id", m.GetQuestion)
			mobileProtected.POST("/question", rdc.RedisPersonalizedCacheMiddleware, m.CreateQuestion)
			mobileProtected.DELETE("/question/:id", rdc.RedisPersonalizedCacheMiddleware, m.DeleteQuestion)
			mobileProtected.GET("/doctor", m.GetAllDoctor)
		}
	}
}

func setupDB() *sql.DB {
	dsn := config.AppConfig.DATABASE_DSN
	message := "Connected to remote database"
	if config.AppConfig.MODE != "test" {
		dsn = config.AppConfig.DATABASE_DSN_LOCAL
		message = "Connected to local database"
	}
	// open db connection
	mysql.RegisterTLSConfig("tidb", &tls.Config{
		MinVersion: tls.VersionTLS12,
		ServerName: "gateway01.ap-southeast-1.prod.aws.tidbcloud.com",
	})

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		mainLogger.Panicf("Can't open connection to database : %v", err.Error())
	}
	mainLogger.Println(message)
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
	// everyday on 10.00 (GMT +7) -> spec : "00 00 03 * * *"
	c.AddFunc("00 00 03 * * *", func() {
		mainLogger.Println("Executing Push Notifications..")
		notification.MockupScheduleNotifications(db, notification.SendRequest)
	})
	c.Start()
	mainLogger.Println("Cron scheduler initialized")
	return c
}

func setupRedisClient() *middleware.RedisClient {
	var middlewareclient *middleware.RedisClient
	var serverMode string
	if config.AppConfig.MODE != "test" {
		serverMode = "local redis server"
		client := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // No password set
			DB:       0,  // Use default DB
		})
		middlewareclient = &middleware.RedisClient{Client: client}
	} else {
		serverMode = "remote redis server"
		url := config.AppConfig.REDIS_URL
		opts, err := redis.ParseURL(url)
		if err != nil {
			mainLogger.Panic(err)
		}
		middlewareclient = &middleware.RedisClient{Client: redis.NewClient(opts)}
	}
	// check connection
	if err := middlewareclient.Client.Ping(context.Background()).Err(); err != nil {
		mainLogger.Panicln("Can't connect to ", serverMode)
	}
	// delete all keys
	middlewareclient.Client.FlushDB(context.Background())
	mainLogger.Println("Connected to", serverMode)
	return middlewareclient
}
