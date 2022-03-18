package config

import (
	"github.com/spf13/viper"
)

type config struct{}

func NewConfig() *config {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.ReadInConfig()

	return &config{}
}

func (*config) GetCallbackURL() string {
	return viper.GetString("callbackURL")
}

func (*config) GetServerURL() string {
	return viper.GetString("serverURL")
}
