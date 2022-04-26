package youtube

import (
	"net/http"

	"github.com/de1phin/music-transfer/internal/log"
)

type YoutubeAPI struct {
	Config
	httpClient *http.Client
	logger     log.Logger
}

func NewYoutubeAPI(config Config, logger log.Logger) *YoutubeAPI {
	return &YoutubeAPI{
		Config:     config,
		logger:     logger,
		httpClient: &http.Client{},
	}
}
