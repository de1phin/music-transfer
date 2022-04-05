package config

import "github.com/spf13/viper"

func (*config) GetTelegramToken() string {
	return viper.GetString("telegram.token")
}
