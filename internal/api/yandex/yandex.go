package yandex

import (
	"net/http"

	"github.com/de1phin/music-transfer/internal/log"
)

type YandexAPI struct {
	logger           log.Logger
	httpClient       *http.Client
	onGetCredentials OnGetCredentials
}

func NewYandexAPI(logger log.Logger) *YandexAPI {
	return &YandexAPI{
		logger:     logger,
		httpClient: &http.Client{},
	}
}
