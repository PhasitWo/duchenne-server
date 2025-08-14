package main

import (
	"context"
	"crypto/tls"
	"database/sql"

	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/PhasitWo/duchenne-server/config"
	"github.com/PhasitWo/duchenne-server/handlers/common"
	"github.com/PhasitWo/duchenne-server/handlers/mobile"
	"github.com/PhasitWo/duchenne-server/handlers/web"
	"github.com/PhasitWo/duchenne-server/middleware"
	"github.com/PhasitWo/duchenne-server/model"
	"github.com/PhasitWo/duchenne-server/services/notification"
	"github.com/robfig/cron"
	"google.golang.org/api/option"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	gomysql "github.com/go-sql-driver/mysql"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var mainLogger = log.New(os.Stdout, "[MAIN] ", log.LstdFlags)

func main() {
	// Load app config
	config.LoadConfig()
	// Setup database connection
	db := setupDB()
	// Setup google cloud storage client
	gcsClient := setupCloudStorageClient()
	// Setup router and handler
	r := setupRouter()
	m := mobile.Init(db)
	w := web.Init(db)
	c := common.Init(db, gcsClient)
	a := middleware.InitActivityLogMiddleware(db)
	attachHandler(r, m, w, c, a)
	// CRON
	if config.AppConfig.ENABLE_CRON {
		cron := InitCronScheduler(w.NotiService)
		defer cron.Stop()
	}
	mainLogger.Println("Server is Live! 🎉")
	r.Run() // listen and serve on 0.0.0.0:8080
}

func attachHandler(r *gin.Engine, m *mobile.MobileHandler, w *web.WebHandler, c *common.CommonHandler, a *middleware.ActivityLogMiddleware) {
	mobile := r.Group("/mobile")
	{
		mobileAuth := mobile.Group("/auth")
		{
			mobileAuth.POST("/login", m.Login)
			mobileAuth.POST("/signup", m.Signup)
			mobileAuth.POST("/logout", middleware.MobileAuthMiddleware, m.Logout)
		}
		mobileProtected := mobile.Group("/api")
		mobileProtected.Use(middleware.MobileAuthMiddleware)
		mobileProtected.Use(a.ActivityLog)
		{
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
			mobileProtected.GET("/device", m.GetAllDevice)
			mobileProtected.POST("/device", m.CreateDevice)
			mobileProtected.GET("/content", c.GetAllContent)
			mobileProtected.GET("/content/:id", c.GetOneContent)
		}
	}
	web := r.Group("/web")
	// r.Static("/static", "./assets")
	// r.GET("/", func(c *gin.Context) {
	// 	c.Redirect(http.StatusMovedPermanently, "/static")
	// })
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "DMD We Care API")
	})
	{
		web.POST("/sendDailyNotifications", w.SendDailyNotifications)
		webAuth := web.Group("/auth")
		{
			webAuth.POST("/login", w.Login)
			webAuth.POST("/logout", w.Logout)
		}
		webProtected := web.Group("/api")
		webProtected.Use(middleware.WebAuthMiddleware)
		webProtected.Use(a.ActivityLog)
		{
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
			webProtected.PUT("/patient/:id/vaccineHistory", middleware.WebRBACMiddleware(middleware.UpdatePatientPermission), w.UpdatePatientVaccineHistory)
			webProtected.PUT("/patient/:id/medicine", middleware.WebRBACMiddleware(middleware.UpdatePatientPermission), w.UpdatePatientMedicine)
			webProtected.DELETE("/patient/:id", middleware.WebRBACMiddleware(middleware.DeletePatientPermission), w.DeletePatient)
			webProtected.GET("/appointment", w.GetAllAppointment)
			webProtected.GET("/appointment/:id", w.GetAppointment)
			webProtected.POST("/appointment", w.CreateAppointment)
			webProtected.PUT("/appointment/:id", w.UpdateAppointment)
			webProtected.DELETE("/appointment/:id", w.DeleteAppointment)
			webProtected.GET("/question", w.GetAllQuestion)
			webProtected.GET("/question/:id", w.GetQuestion)
			webProtected.PUT("/question/:id/answer", w.AnswerQuestion)
			webProtected.GET("/content", c.GetAllContent)
			webProtected.GET("/content/:id", c.GetOneContent)
			webProtected.POST("/content", w.CreateContent)
			webProtected.PUT("/content/:id", w.UpdateContent)
			webProtected.DELETE("/content/:id", w.DeleteContent)
			webProtected.POST("/image/upload", c.UploadImage)
		}
	}
}

func setupDB() *gorm.DB {
	dsn := config.AppConfig.DATABASE_DSN
	// open db connection
	gomysql.RegisterTLSConfig("tidb", &tls.Config{
		MinVersion: tls.VersionTLS12,
		ServerName: "gateway01.ap-southeast-1.prod.aws.tidbcloud.com",
	})
	customDB, err := sql.Open("mysql", dsn)
	if err != nil {
		mainLogger.Panicf("can't open connection to database : %v", err.Error())
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: customDB,
	}), &gorm.Config{SkipDefaultTransaction: true})

	if err != nil {
		mainLogger.Panicf("GORM can't connection to database : %v", err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		mainLogger.Panicf("can't setup connection config : %v", err.Error())
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Migrate
	db.AutoMigrate(
		&model.ActivityLog{},
		&model.Appointment{},
		&model.Device{},
		&model.Doctor{},
		&model.Patient{},
		&model.Question{},
		&model.Content{},
	)

	mainLogger.Println("connected to the database")
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
	if config.AppConfig.MODE == "dev" {
		gin.SetMode(gin.TestMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowCredentials = true
	corsConfig.AllowOrigins = config.AppConfig.CORS_ALLOW
	r.Use(cors.New(corsConfig))
	return r
}

func InitCronScheduler(service notification.INotificationService) *cron.Cron {
	c := cron.New()
	// everyday on 10.00 (GMT +7) -> spec : "00 00 03 * * *"
	c.AddFunc("00 00 03 * * *", func() {
		mainLogger.Println("executing Push Notifications..")
		service.SendDailyNotifications(nil)
	})
	c.Start()
	mainLogger.Println("cron scheduler initialized")
	return c
}

func setupCloudStorageClient() *storage.Client {
	var opts []option.ClientOption
	const credFileName = "./storage-cred.json"
	_, err := os.Stat(credFileName)
	if !os.IsNotExist(err) {
		mainLogger.Println("using service account json file for google cloud storage")
		opts = append(opts, option.WithCredentialsFile(credFileName))
	}
	gcsClient, err := storage.NewClient(context.Background(), opts...)
	if err != nil {
		mainLogger.Panicf("can't create google cloud storage client : %v", err.Error())
	}
	mainLogger.Println("connected to google cloud storage")
	return gcsClient
}

// func setupRedisClient() *middleware.RedisClient {
// 	// check config
// 	if !config.AppConfig.ENABLE_REDIS {
// 		return nil
// 	}
// 	var middlewareclient *middleware.RedisClient
// 	var serverMode string
// 	if config.AppConfig.MODE != "test" {
// 		serverMode = "local redis server"
// 		client := redis.NewClient(&redis.Options{
// 			Addr:     "localhost:6379",
// 			Password: "", // No password set
// 			DB:       0,  // Use default DB
// 		})
// 		middlewareclient = &middleware.RedisClient{Client: client}
// 	} else {
// 		serverMode = "remote redis server"
// 		url := config.AppConfig.REDIS_URL
// 		opts, err := redis.ParseURL(url)
// 		if err != nil {
// 			mainLogger.Panic(err)
// 		}
// 		middlewareclient = &middleware.RedisClient{Client: redis.NewClient(opts)}
// 	}
// 	// check connection
// 	if err := middlewareclient.Client.Ping(context.Background()).Err(); err != nil {
// 		mainLogger.Panicln("Can't connect to ", serverMode)
// 	}
// 	// delete all keys
// 	middlewareclient.Client.FlushDB(context.Background())
// 	mainLogger.Println("Connected to", serverMode)
// 	return middlewareclient
// }
