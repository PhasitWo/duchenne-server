package main

import (
	// "net/http"
	"context"
	// "crypto/tls"
	// "database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PhasitWo/duchenne-server/config"
	"github.com/PhasitWo/duchenne-server/handlers/mobile"
	"github.com/PhasitWo/duchenne-server/handlers/web"
	"github.com/PhasitWo/duchenne-server/middleware"
	"github.com/PhasitWo/duchenne-server/model"
	"github.com/PhasitWo/duchenne-server/notification"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	// "github.com/PhasitWo/duchenne-server/repository"

	// "github.com/go-sql-driver/mysql"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var mainLogger = log.New(os.Stdout, "[MAIN] ", log.LstdFlags)

func main() {
	// Load app config
	config.LoadConfig()
	// Setup database connection
	db := setupDB()
	// Setup redis
	rdc := setupRedisClient()
	// Setup router and handler
	r := setupRouter()
	m := mobile.Init(db)
	w := web.Init(db)
	attachHandler(r, m, w, rdc)
	// CRON
	c := InitCronScheduler(db)
	defer c.Stop()
	mainLogger.Println("Server is live! 🎉")
	r.Run() // listen and serve on 0.0.0.0:8080
}

func attachHandler(r *gin.Engine, m *mobile.MobileHandler, w *web.WebHandler, rdc *middleware.RedisClient) {
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
			mobileProtected.GET("/appointment", rdc.UseRedisMiddleware(m.GetAllPatientAppointment)...)
			mobileProtected.GET("/appointment/:id", m.GetAppointment)
			mobileProtected.POST("/appointment", rdc.UseRedisMiddleware(m.CreateAppointment)...)
			mobileProtected.DELETE("/appointment/:id", rdc.UseRedisMiddleware(m.DeleteAppointment)...)
			mobileProtected.GET("/question", rdc.UseRedisMiddleware(m.GetAllPatientQuestion)...)
			mobileProtected.GET("/question/:id", m.GetQuestion)
			mobileProtected.POST("/question", rdc.UseRedisMiddleware(m.CreateQuestion)...)
			mobileProtected.DELETE("/question/:id", rdc.UseRedisMiddleware(m.DeleteQuestion)...)
			mobileProtected.GET("/doctor", m.GetAllDoctor)
			mobileProtected.GET("/device", m.GetAllDevice)
			mobileProtected.POST("/device", m.CreateDevice)
		}
	}
	web := r.Group("/web")
	r.Static("/static", "./assets")
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/static")
	})
	{
		webAuth := web.Group("/auth")
		{
			webAuth.POST("/login", w.Login)
			webAuth.POST("/logout", w.Logout)
		}
		webProtected := web.Group("/api")
		webProtected.Use(middleware.WebAuthMiddleware)
		{
			webProtected.POST("/test", w.Test)
			webProtected.GET("/userData", w.GetUserData)
			webProtected.GET("/profile", w.GetProfile)
			webProtected.PUT("/profile", w.UpdateProfile)
			webProtected.GET("/doctor", w.GetAllDoctor)
			webProtected.POST("/doctor", middleware.WebRBACMiddleware(middleware.CreateDoctorPermission), w.CreateDoctor)
			webProtected.GET("/doctor/:id", w.GetDoctor)
			webProtected.PUT("/doctor/:id", middleware.WebRBACMiddleware(middleware.UpdateDoctorPermission), w.UpdateDoctor)
			webProtected.DELETE("/doctor/:id", middleware.WebRBACMiddleware(middleware.DeleteDoctorPermission), w.DeleteDoctor)
			webProtected.GET("/patient", w.GetAllPatient)
			webProtected.POST("/patient", middleware.WebRBACMiddleware(middleware.CreatePatientPermission), w.CreatePatient)
			webProtected.GET("/patient/:id", w.GetPatient)
			webProtected.PUT("/patient/:id", middleware.WebRBACMiddleware(middleware.UpdatePatientPermission), w.UpdatePatient)
			webProtected.DELETE("/patient/:id", middleware.WebRBACMiddleware(middleware.DeletePatientPermission), w.DeletePatient)
			webProtected.GET("/appointment", w.GetAllAppointment)
			webProtected.GET("/appointment/:id", w.GetAppointment)
			webProtected.POST("/appointment", w.CreateAppointment)
			webProtected.PUT("/appointment/:id", w.UpdateAppointment)
			webProtected.DELETE("/appointment/:id", w.DeleteAppointment)
			webProtected.GET("/question", w.GetAllQuestion)
			webProtected.GET("/question/:id", w.GetQuestion)
			webProtected.PUT("/question/:id/answer", w.AnswerQuestion)
		}
	}
}

func setupDB() *gorm.DB {
	dsn := config.AppConfig.DATABASE_DSN
	message := "Connected to remote database"
	if config.AppConfig.MODE == "dev" {
		dsn = config.AppConfig.DATABASE_DSN_LOCAL
		message = "Connected to local database"
	}
	// open db connection
	// mysql.RegisterTLSConfig("tidb", &tls.Config{
	// 	MinVersion: tls.VersionTLS12,
	// 	ServerName: "gateway01.ap-southeast-1.prod.aws.tidbcloud.com",
	// })

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	// db, err := sql.Open("mysql", dsn)
	if err != nil {
		mainLogger.Panicf("Can't open connection to database : %v", err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		mainLogger.Panicf("Can't setup connection config : %v", err.Error())
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Migrate
	db.AutoMigrate(
		&model.Appointment{},
		&model.Device{},
		&model.Doctor{},
		&model.Patient{},
		&model.Question{},
	)

	mainLogger.Println(message)
	// db.SetConnMaxLifetime(time.Minute * 3)
	// db.SetMaxOpenConns(10)
	// db.SetMaxIdleConns(10)
	// // ping db
	// if err = db.Ping(); err != nil {
	// 	mainLogger.Panic("Can't connect to database")
	// }
	return db
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowCredentials = true
	corsConfig.AllowOrigins = []string{"http://localhost:5173", "http://localhost:4173", "https://duchenne-web.onrender.com/"}
	r.Use(cors.New(corsConfig))
	return r
}

func InitCronScheduler(db *gorm.DB) *cron.Cron {
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
	// check config
	if !config.AppConfig.ENABLE_REDIS {
		return nil
	}
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

// func testRepo() {
// 	db := setupDB()
// 	repo := repository.New(db)
// 	mn := "mid na"
// 	res, err := repo.CreateDoctor(model.Doctor{
// 		Id:         -1,
// 		FirstName:  "myrepo",
// 		MiddleName: &mn,
// 		LastName:   "ln na",
// 		Username:   "myrepousername",
// 		Password:   "1234",
// 		Role:       model.USER,
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// 	log.Println(res)
// }
