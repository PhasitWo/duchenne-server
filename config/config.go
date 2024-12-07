package config

import (
	"fmt"
	"github.com/spf13/viper"
	"reflect"
)

type config struct {
	MODE               string
	DATABASE_DSN       string
	DATABASE_DSN_LOCAL string
	JWT_KEY            string
	MAX_DEVICE         int
	NOTIFY_IN_RANGE    int
}

// shared config across packages
var AppConfig = config{MODE: "dev"}

func LoadConfig() {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		panic("Can't read config file")
	}
	// load config by struct field's name
	f := reflect.ValueOf(&AppConfig).Elem()
	fmt.Print("Loading Config:\n")
	for i := 0; i < f.NumField(); i++ {
		field := f.Field(i)
		fieldName := f.Type().Field(i).Name
		fieldValue := field.Interface()
		if !viper.IsSet(fieldName) {
			panic(fmt.Sprintf("Can't read %s from .env", fieldName))
		}
		switch fieldValue.(type) {
		case string:
			field.SetString(viper.GetString(fieldName))
		case int:
			field.SetInt(int64(viper.GetInt(fieldName)))
		default:
			panic("invalid config type")
		}
		fmt.Printf("\t%-15s\t=>\t%-10v\n", fieldName, f.Field(i).Interface())
	}
	fmt.Printf("Config Loaded\n\n")
	fmt.Printf("Server is running in mode `%v`\n", AppConfig.MODE)
	if AppConfig.MODE == "dev" {
		AppConfig.DATABASE_DSN = AppConfig.DATABASE_DSN_LOCAL
		fmt.Printf("Replacing AppConfig.DATABASE_DSN with => %v\n\n", AppConfig.DATABASE_DSN_LOCAL)
	}
	fmt.Printf("Uses this DSN => %v\n", AppConfig.DATABASE_DSN)
}
