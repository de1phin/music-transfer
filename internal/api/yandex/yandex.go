package yandex

import (
	"net/http"

	"github.com/de1phin/music-transfer/internal/log"
)

type YandexAPI struct {
	fixedAuthMagicToken string
	logger              log.Logger
	httpClient          *http.Client
	onGetCredentials    OnGetCredentials
}

func NewYandexAPI(logger log.Logger, fixedAuthMagicToken string) *YandexAPI {
	return &YandexAPI{
		fixedAuthMagicToken: fixedAuthMagicToken,
		logger:              logger,
		httpClient:          &http.Client{},
	}
}
