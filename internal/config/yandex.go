package config

import "github.com/spf13/viper"

func (*config) GetYandexMagicToken() string {
	return viper.GetString("yandex.magic")
}
