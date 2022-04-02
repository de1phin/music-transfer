package main

import (
	spotifyAPI "github.com/de1phin/music-transfer/internal/api/spotify"
	"github.com/de1phin/music-transfer/internal/config"
	"github.com/de1phin/music-transfer/internal/interactor"
	consoleInteractor "github.com/de1phin/music-transfer/internal/interactor/interactors/console"
	consoleValidator "github.com/de1phin/music-transfer/internal/interactor/validator/console"
	"github.com/de1phin/music-transfer/internal/mux"
	"github.com/de1phin/music-transfer/internal/server/callback"
	"github.com/de1phin/music-transfer/internal/service/spotify"
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

	services := []mux.Service{
		spotify,
	}

	consoleInteractor := consoleInteractor.NewConsoleInteractor(17)
	consoleValidator := consoleValidator.Validator{}
	console := interactor.NewInteractorSpec(consoleInteractor, consoleValidator)

	stateStorage := cache.NewCacheStorage[mux.UserState]()
	mux := mux.NewMux(services, console, stateStorage)

	go server.Run()
	mux.Run()

}
