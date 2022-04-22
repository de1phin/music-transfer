package config

import (
	"github.com/spf13/viper"
)

type config struct{}

func NewConfig(configPath string, configName string, configType string) *config {
	viper.AddConfigPath(configPath)
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)

	viper.ReadInConfig()

	return &config{}
}
