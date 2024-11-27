package config

import (
	"fmt"
	"github.com/spf13/viper"
	"reflect"
)

type config struct {
	DATABASE_DSN string
	JWT_KEY      []byte
}

// shared config across packages
var AppConfig config

func LoadConfig() {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		panic("Can't read config file")
	}
	AppConfig.DATABASE_DSN = viper.GetString("DATABASE_DSN")
	AppConfig.JWT_KEY = []byte(viper.GetString("JWT_KEY"))
	checkConfig()
}

func checkConfig() {
	f := reflect.ValueOf(AppConfig)
	fmt.Print("Loading Config:\n")
	for i := 0; i < f.NumField(); i++ {
		field := f.Field(i)
		fieldName := f.Type().Field(i).Name
		fieldValue := field.Interface()
		if fieldValue == "" {
			panic(fmt.Sprintf("Can't read %s from .env", fieldName))
		}
		fmt.Printf("\t%-15s\t=>\t%-10s\n", fieldName, fieldValue)
	}
}
