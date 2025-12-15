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
	MODE                   string
	DATABASE_DSN           string
	JWT_KEY                string
	JWT_REFRESH_KEY        string
	MAX_DEVICE             int
	NOTIFY_IN_RANGE        int
	NOTIFY_SECRET          string
	ENABLE_CRON            bool
	CORS_ALLOW             []string
	REQUIRE_MOBILE_VERSION string
	ANDROID_STORE_LINK     string
	IOS_STORE_LINK         string
}

// shared config across packages
var AppConfig = config{}
var defaultConfig = config{
	MODE:                   "dev",
	DATABASE_DSN:           "root:superuser@tcp(127.0.0.1)/master",
	JWT_KEY:                "SAMPLE_KEY",
	JWT_REFRESH_KEY:        "REFRESH_KEY",
	MAX_DEVICE:             3,
	NOTIFY_IN_RANGE:        3,
	NOTIFY_SECRET:          "SAMPLE_SECRET",
	ENABLE_CRON:            false,
	CORS_ALLOW:             []string{"http://localhost:5173", "http://localhost:4173", "https://duchenne-web.onrender.com"},
	REQUIRE_MOBILE_VERSION: "0.0.0",
	ANDROID_STORE_LINK:     "https://play.google.com",
	IOS_STORE_LINK:         "https://apps.apple.com/",
}

func LoadConfig() {
	configLogger := log.New(os.Stdout, "[CONFIG] ", 0)
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		if errors.Is(err, viper.ConfigFileNotFoundError{}) || errors.Is(err, os.ErrNotExist) {
			configLogger.Println(".env file is not found, read from env variables instead")
			viper.AutomaticEnv()
		} else {
			panic("Can't read config file")
		}
	}
	// load config by struct field's name
	f := reflect.ValueOf(&AppConfig).Elem()
	df := reflect.ValueOf(&defaultConfig).Elem()
	configLogger.Print("loading config:\n")
	for i := 0; i < f.NumField(); i++ {
		field := f.Field(i)
		fieldName := f.Type().Field(i).Name
		fieldValue := field.Interface()
		if !viper.IsSet(fieldName) {
			// use default
			field.Set(reflect.ValueOf(df.Field(i).Interface()))
			fmt.Printf("\t%-15s\t=>\t%-10v(default)\n", fieldName, f.Field(i).Interface())
			continue
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
	configLogger.Printf("config loaded\n")
}
