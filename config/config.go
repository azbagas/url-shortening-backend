package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	AppName            string
	AppPort            int
	DatabaseUrl        string
	AccessTokenSecret  string
	RefreshTokenSecret string
}

var AppConfig Config

func LoadConfig() {
	config := viper.New()
	config.SetConfigFile(".env")
	config.AutomaticEnv()

	err := config.ReadInConfig()
	if err != nil {
		panic(err)
	}

	AppConfig = Config{
		AppName:            config.GetString("APP_NAME"),
		AppPort:            config.GetInt("APP_PORT"),
		DatabaseUrl:        config.GetString("DATABASE_URL"),
		AccessTokenSecret:  config.GetString("ACCESS_TOKEN_SECRET"),
		RefreshTokenSecret: config.GetString("REFRESH_TOKEN_SECRET"),
	}
}
