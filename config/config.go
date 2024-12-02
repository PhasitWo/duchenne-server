package config

import (
	"fmt"
	"github.com/spf13/viper"
	"reflect"
)

type config struct {
	DATABASE_DSN string
	JWT_KEY      string
	MAX_DEVICE   int
}

// shared config across packages
var AppConfig config

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
}