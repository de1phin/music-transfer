package main

import (
	spotifyAPI "github.com/de1phin/music-transfer/internal/api/spotify"
	youtubeAPI "github.com/de1phin/music-transfer/internal/api/youtube"
	"github.com/de1phin/music-transfer/internal/config"
	"github.com/de1phin/music-transfer/internal/interactor"
	consoleInteractor "github.com/de1phin/music-transfer/internal/interactor/interactors/console"
	consoleValidator "github.com/de1phin/music-transfer/internal/interactor/validator/console"
	"github.com/de1phin/music-transfer/internal/mux"
	"github.com/de1phin/music-transfer/internal/server/callback"
	"github.com/de1phin/music-transfer/internal/service/mock"
	"github.com/de1phin/music-transfer/internal/service/spotify"
	"github.com/de1phin/music-transfer/internal/service/youtube"
	"github.com/de1phin/music-transfer/internal/storage/cache"
)

func main() {

	config := config.NewConfig()

	spotifyConfig := spotify.SpotifyConfig{
		Scopes: config.GetSpotifyScope(),
		Client: spotifyAPI.Client{
			ID:     config.GetSpotifyClientID(),
			Secret: config.GetSpotifyClientSecret(),
		},
	}
	server := callback.NewCallbackServer(config.GetServerURL())

	spotifyStorage := cache.NewCacheStorage[spotifyAPI.Credentials]()
	spotifyAPI := spotifyAPI.NewSpotifyAPI(spotifyConfig.Client, config.GetCallbackURL())
	spotify := spotify.NewSpotifyService(spotifyConfig, config.GetCallbackURL(), spotifyAPI, spotifyStorage)
	spotifyAPI.BindHandler(server.ServeMux, spotify.OnGetTokens)

	youtubeStorage := cache.NewCacheStorage[youtubeAPI.Credentials]()
	youtubeConfig := youtubeAPI.YoutubeConfig{
		APIKey:       config.GetYouTubeApiKEY(),
		ClientID:     config.GetYouTubeClientID(),
		ClientSecret: config.GetYouTubeClientSecret(),
		Scopes:       config.GetYouTubeScope(),
		RedirectURI:  config.GetCallbackURL() + "/youtube",
	}
	youtubeAPI := youtubeAPI.NewYoutubeAPI(&youtubeConfig)
	youtube := youtube.NewYouTubeService(youtubeAPI, youtubeStorage, &youtubeConfig)
	youtubeAPI.BindHandler(server.ServeMux, youtube.OnGetTokens)

	services := []mux.Service{
		spotify,
		youtube,
		mock.NewMockService(),
	}

	consoleInteractor := consoleInteractor.NewConsoleInteractor(17)
	consoleValidator := consoleValidator.Validator{}
	console := interactor.NewInteractorSpec(consoleInteractor, consoleValidator)

	stateStorage := cache.NewCacheStorage[mux.UserState]()
	transferStorage := cache.NewCacheStorage[mux.Transfer]()
	mux := mux.NewMux(services, console, stateStorage, transferStorage)

	go server.Run()
	mux.Run()

}
