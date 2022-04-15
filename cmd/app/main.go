package main

import (
	"os"

	spotifyAPI "github.com/de1phin/music-transfer/internal/api/spotify"
	youtubeAPI "github.com/de1phin/music-transfer/internal/api/youtube"
	"github.com/de1phin/music-transfer/internal/config"
	consoleAdapter "github.com/de1phin/music-transfer/internal/interactor/adapters/console"
	telegramAdapter "github.com/de1phin/music-transfer/internal/interactor/adapters/telegram"
	"github.com/de1phin/music-transfer/internal/interactor/interactors/console"
	"github.com/de1phin/music-transfer/internal/interactor/interactors/telegram"
	logger "github.com/de1phin/music-transfer/internal/log/file_logger"
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
	server := callback.NewCallbackServer(config.GetServerHostname())

	spotifyStorage := cache.NewCacheStorage[int64, spotifyAPI.Credentials]()
	spotifyAPI := spotifyAPI.NewSpotifyAPI(spotifyConfig.Client, "http://"+config.GetServerHostname())
	spotify := spotify.NewSpotifyService(spotifyConfig, "http://"+config.GetServerHostname(), spotifyAPI, spotifyStorage)
	spotifyAPI.BindHandler(server.ServeMux, spotify.OnGetTokens)

	youtubeStorage := cache.NewCacheStorage[int64, youtubeAPI.Credentials]()
	youtubeConfig := youtubeAPI.YoutubeConfig{
		APIKey:       config.GetYouTubeApiKEY(),
		ClientID:     config.GetYouTubeClientID(),
		ClientSecret: config.GetYouTubeClientSecret(),
		Scopes:       config.GetYouTubeScope(),
		RedirectURI:  "http://" + config.GetServerHostname() + "/youtube",
	}
	youtubeAPI := youtubeAPI.NewYoutubeAPI(&youtubeConfig)
	youtube := youtube.NewYouTubeService(youtubeAPI, youtubeStorage, &youtubeConfig)
	youtubeAPI.BindHandler(server.ServeMux, youtube.OnGetTokens)

	services := []mux.Service{
		spotify,
		youtube,
		mock.NewMockService(),
	}

	userStateStorage := cache.NewCacheStorage[int64, mux.UserState]()

	telegram, err := telegram.NewTelegramBot(config.GetTelegramToken())
	if err != nil {
		panic("Telegram init error: " + err.Error())
	}
	telegramAdapter := telegramAdapter.NewTelegramAdapter(telegram, userStateStorage)

	console := console.NewConsoleInteractor()
	consoleAdapter := consoleAdapter.NewConsoleAdapter(console, 17)

	transferStorage := cache.NewCacheStorage[int64, mux.Transfer]()
	idStorage := cache.NewCacheStorage[string, int64]()

	interactors := []mux.Interactor{
		telegramAdapter,
		consoleAdapter,
	}

	fileLogger, err := logger.NewFileLogger("./log/a.log")
	if err != nil {
		panic("FileLogger init error: " + err.Error())
	}
	mux := mux.NewMux(services, interactors, transferStorage, idStorage, fileLogger)

	go server.Run()

	muxQuit := make(chan struct{})
	mux.Run(muxQuit)

	c := make(chan os.Signal)
	<-c
	muxQuit <- struct{}{}

	fileLogger.Close()
}
