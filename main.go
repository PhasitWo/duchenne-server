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

	"github.com/gin-gonic/gin"

	// "github.com/PhasitWo/duchenne-server/repository"

	"github.com/go-sql-driver/mysql"
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
	r.Run() // listen and serve on 0.0.0.0:8080
}

// var apmtQuery = `
// select appointment.id ,date, device.id , expo_token, appointment.patient_id from appointment 
// inner join device on appointment.patient_id = device.patient_id  
// order by appointment.id asc
// `

// func scheduleNotifications(db *sql.DB) {
// 	rows, err := db.Query(apmtQuery)
// 	if err != nil {
// 		fmt.Println("scheduleNotifications : Can't query database")
// 		return
// 	}
// 	defer rows.Close()
// 	res := []model.AppointmentDevice{}
// 	for rows.Next() {
// 		var ad model.AppointmentDevice
// 		if err := rows.Scan(
// 			&ad.AppointmentId,
// 			&ad.Date,
// 			&ad.DeviceId,
// 			&ad.ExpoToken,
// 			&ad.PatientId,
// 		); err != nil {
// 			fmt.Printf("scheduleNotifications : %v", err.Error())
// 			return
// 		}
// 		res = append(res, ad)
// 	}
// 	if err := rows.Err(); err != nil {
// 		fmt.Printf("scheduleNotifications : %v", err.Error())
// 		return
// 	}
// 	// schedule
// 	for _, element := range res {
// 		fmt.Printf("scheduled appointment_id:%v , device_id:%v --> push token: %v\n", element.AppointmentId, element.DeviceId, element.ExpoToken)
// 	}
// }
