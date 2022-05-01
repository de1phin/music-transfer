package config

import (
	"github.com/de1phin/music-transfer/internal/api/spotify"
	"github.com/de1phin/music-transfer/internal/api/yandex"
	"github.com/de1phin/music-transfer/internal/api/youtube"
	"github.com/de1phin/music-transfer/internal/interactor/interactors/telegram"
	"github.com/de1phin/music-transfer/internal/server"
	"github.com/de1phin/music-transfer/internal/storage/postgres"
	"github.com/de1phin/music-transfer/internal/storage/redis"
	"github.com/spf13/viper"
)

type Config struct {
	Server   server.Config   `yaml:"server"`
	Spotify  spotify.Config  `yaml:"spotify"`
	Yandex   yandex.Config   `yaml:"yandex"`
	Youtube  youtube.Config  `yaml:"youtube"`
	Telegram telegram.Config `yaml:"telegram"`
	Postgres postgres.Config `yaml:"postgres"`
	Redis    redis.Config    `yaml:"redis"`
}

func ReadConfig(cfgPath, cfgName, cfgType string) (cfg Config, err error) {
	viper.AddConfigPath(cfgPath)
	viper.SetConfigName(cfgName)
	viper.SetConfigType(cfgType)

	err = viper.ReadInConfig()
	if err != nil {
		return cfg, err
	}

	err = viper.Unmarshal(&cfg)

	return cfg, err
}
