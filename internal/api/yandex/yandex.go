package yandex

import (
	"net/http"

	"github.com/de1phin/music-transfer/internal/log"
)

type YandexAPI struct {
	Config
	logger           log.Logger
	httpClient       *http.Client
	onGetCredentials OnGetCredentials
}

func NewYandexAPI(logger log.Logger, config Config) *YandexAPI {
	return &YandexAPI{
		Config:     config,
		logger:     logger,
		httpClient: &http.Client{},
	}
}
