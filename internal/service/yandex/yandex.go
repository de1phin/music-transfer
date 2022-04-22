package yandex

import (
	"github.com/de1phin/music-transfer/internal/api/yandex"
	"github.com/de1phin/music-transfer/internal/log"
	"github.com/de1phin/music-transfer/internal/storage"
)

type Yandex struct {
	api         *yandex.YandexAPI
	userStorage storage.Storage[int64, yandex.Credentials]
	logger      log.Logger
}

func NewYandexService(api *yandex.YandexAPI, userStorage storage.Storage[int64, yandex.Credentials], logger log.Logger) *Yandex {
	return &Yandex{
		api:         api,
		userStorage: userStorage,
		logger:      logger,
	}
}

func (*Yandex) Name() string {
	return "yandex"
}
