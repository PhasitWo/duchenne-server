package main

import (
	// "net/http"
	"context"
	"crypto/tls"
	"database/sql"

	// "crypto/tls"
	// "database/sql"
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
	"github.com/PhasitWo/duchenne-server/notification"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron"
	"google.golang.org/api/option"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	// "github.com/PhasitWo/duchenne-server/repository"

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
	// Setup redis
	rdc := setupRedisClient()
	// Setup google cloud storage client
	gcsClient := setupCloudStorageClient()
	// Setup router and handler
	r := setupRouter()
	m := mobile.Init(db)
	w := web.Init(db)
	c := common.Init(db, gcsClient)
	attachHandler(r, m, w, c, rdc)
	// CRON
	cron := InitCronScheduler(db)
	defer cron.Stop()
	mainLogger.Println("Server is live! 🎉")
	r.Run() // listen and serve on 0.0.0.0:8080
}

func attachHandler(r *gin.Engine, m *mobile.MobileHandler, w *web.WebHandler, c *common.CommonHandler, rdc *middleware.RedisClient) {
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
		{
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
		c.JSON(http.StatusOK, "Duchenne Server API")
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
			g, ok := m.DBConn.(*gorm.DB)
			if !ok {
				panic("can't cast to *gorm.DB")
			}
			webProtected.POST("/sendDailyNotifications", func(c *gin.Context) {
				notification.SendDailyNotifications(g, notification.SendRequest)
				c.Status(200)
			})
			webProtected.GET("/content", c.GetAllContent)
			webProtected.GET("/content/:id", c.GetOneContent)
			webProtected.POST("/content", w.CreateContent)
			webProtected.PUT("/content/:id", w.UpdateContent)
			webProtected.DELETE("/content/:id", w.DeleteContent)
			webProtected.POST("/image/upload", c.UploadImage)
		}
	}
	// common := r.Group("/common/api").Use(middleware.CommonAuthMiddleware)
	// {
	// 	common.GET("/content", c.GetAllContent)
	// 	common.GET("/content/:id", c.GetOneContent)
	// }
}

func setupDB() *gorm.DB {
	dsn := config.AppConfig.DATABASE_DSN
	message := "Connected to remote database"
	if config.AppConfig.MODE == "dev" {
		dsn = config.AppConfig.DATABASE_DSN_LOCAL
		message = "Connected to local database"
	}
	// open db connection
	gomysql.RegisterTLSConfig("tidb", &tls.Config{
		MinVersion: tls.VersionTLS12,
		ServerName: "gateway01.ap-southeast-1.prod.aws.tidbcloud.com",
	})
	customDB, err := sql.Open("mysql", dsn)
	if err != nil {
		mainLogger.Panicf("Can't open connection to database : %v", err.Error())
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: customDB,
	}), &gorm.Config{SkipDefaultTransaction: true})

	if err != nil {
		mainLogger.Panicf("GORM can't connection to database : %v", err.Error())
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
		&model.Content{},
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
	corsConfig.AllowOrigins = config.AppConfig.CORS_ALLOW
	r.Use(cors.New(corsConfig))
	return r
}

func InitCronScheduler(db *gorm.DB) *cron.Cron {
	c := cron.New()
	// everyday on 10.00 (GMT +7) -> spec : "00 00 03 * * *"
	c.AddFunc("00 00 03 * * *", func() {
		mainLogger.Println("Executing Push Notifications..")
		notification.SendDailyNotifications(db, notification.SendRequest)
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

func setupCloudStorageClient() *storage.Client {
	gcsClient, err := storage.NewClient(context.Background(), option.WithCredentialsFile("./storage-cred.json"))
	if err != nil {
		mainLogger.Panicf("Can't create google cloud storage client : %v", err.Error())
	}
	mainLogger.Println("Connected to google cloud storage")
	return gcsClient
}
