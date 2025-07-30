package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/spf13/viper"
)

type config struct {
	MODE               string
	DATABASE_DSN       string
	DATABASE_DSN_LOCAL string
	ENABLE_REDIS       bool
	REDIS_URL          string
	JWT_KEY            string
	MAX_DEVICE         int
	NOTIFY_IN_RANGE    int
	ENABLE_CRON        bool
	SERVER_DOMAIN      string
	CORS_ALLOW         []string
}

// shared config across packages
var AppConfig = config{MODE: "dev"}

func LoadConfig() {
	configLogger := log.New(os.Stdout, "[CONFIG] ", 0)
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		if errors.Is(err, viper.ConfigFileNotFoundError{}) || errors.Is(err, os.ErrNotExist) {
			configLogger.Println(".env file is not found, read from env variables intead")
			viper.AutomaticEnv()
		} else {
			panic("Can't read config file")
		}
	}
	// load config by struct field's name
	f := reflect.ValueOf(&AppConfig).Elem()
	configLogger.Print("Loading Config:\n")
	for i := 0; i < f.NumField(); i++ {
		field := f.Field(i)
		fieldName := f.Type().Field(i).Name
		fieldValue := field.Interface()
		if !viper.IsSet(fieldName) {
			configLogger.Panicf("Can't read %s from .env", fieldName)
		}
		switch fieldValue.(type) {
		case string:
			field.SetString(viper.GetString(fieldName))
		case int:
			field.SetInt(int64(viper.GetInt(fieldName)))
		case bool:
			field.SetBool(viper.GetBool(fieldName))
		case []string:
			field.Set(reflect.ValueOf(viper.GetStringSlice(fieldName)))
		default:
			panic("invalid config type")
		}
		fmt.Printf("\t%-15s\t=>\t%-10v\n", fieldName, f.Field(i).Interface())
	}
	if AppConfig.MODE == "dev" {
		AppConfig.SERVER_DOMAIN = "127.0.0.1"
	}

	configLogger.Printf("Config Loaded\n")
	configLogger.Printf("Server is running in mode `%v`\n", AppConfig.MODE)
	configLogger.Printf("Server Domain -> `%v`\n\n", AppConfig.SERVER_DOMAIN)
}

type constants struct {
	WEB_ACCESS_COOKIE_NAME  string
	WEB_REFRESH_COOKIE_NAME string
}

var Constants = constants{
	WEB_ACCESS_COOKIE_NAME:  "web_access_token",
	WEB_REFRESH_COOKIE_NAME: "web_refresh_token",
}
