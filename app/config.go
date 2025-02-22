package app

import (
	"github.com/azbagas/url-shortening-backend/helper"
	"github.com/spf13/viper"
)

func NewConfig() *viper.Viper {
	config := viper.New()
	config.SetConfigFile(".env")
	config.AddConfigPath(".")
	config.AutomaticEnv()

	err := config.ReadInConfig()
	helper.PanicIfError(err)

	return config
}
