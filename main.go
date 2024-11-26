package main

import (
	// "net/http"
	"database/sql"
	"fmt"

	"github.com/PhasitWo/duchenne-server/endpoint/mobile"
	"github.com/gin-gonic/gin"
	// "github.com/PhasitWo/duchenne-server/repository"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

/*
mobile endpoints -> /mobile/
AUTH
POST /login	-> authenticate user
POST /signup -> verify user info

PROFILE
GET /profle -> return patient profile data

APPOINTMENT
GET /appointment  -> return maximum 20 of patient's appointments
POST /appointment -> create new appointment

ASK
GET /ask -> return patient's question history
GET /ask/:id -> return patient's question and doctor's answer
POST /ask -> create new question
*/

func main() {
	// read config
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		panic("Can't read config file")
	}
	databaseDSN := viper.GetString("DATABASE_DSN")
	// open db connection
	db, err := sql.Open("mysql", databaseDSN)
	if err != nil {
		panic(fmt.Sprintf("Can't connect to database : %v", err.Error()))
	}
	// setup router
	r := gin.Default()
	m := mobile.Init(db)
	r.POST("/mobile/login", m.Login)
	r.POST("/mobile/signup", m.Signup)
	r.GET("/mobile/test", m.Test)
	r.Run() // listen and serve on 0.0.0.0:8080

}
