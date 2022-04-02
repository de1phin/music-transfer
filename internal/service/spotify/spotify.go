package spotify

import (
	"github.com/de1phin/music-transfer/internal/api/spotify"
	"github.com/de1phin/music-transfer/internal/storage"
)

type SpotifyConfig struct {
	Client spotify.Client
	Scopes string
}

type spotifyService struct {
	scopes      string
	client      spotify.Client
	api         *spotify.SpotifyAPI
	redirectURI string
	storage     storage.Storage[spotify.Credentials]
}

func NewSpotifyService(config SpotifyConfig, redirectURI string, spotifyAPI *spotify.SpotifyAPI, storage storage.Storage[spotify.Credentials]) *spotifyService {
	return &spotifyService{
		scopes:      config.Scopes,
		client:      config.Client,
		api:         spotifyAPI,
		redirectURI: redirectURI,
		storage:     storage,
	}
}

func (spotify *spotifyService) Name() string {
	return "spotify"
}
