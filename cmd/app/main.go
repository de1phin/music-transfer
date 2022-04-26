package main

import (
	"os"

	spotifyAPI "github.com/de1phin/music-transfer/internal/api/spotify"
	yandexAPI "github.com/de1phin/music-transfer/internal/api/yandex"
	youtubeAPI "github.com/de1phin/music-transfer/internal/api/youtube"
	"github.com/de1phin/music-transfer/internal/config"
	consoleAdapter "github.com/de1phin/music-transfer/internal/interactor/adapters/console"
	telegramAdapter "github.com/de1phin/music-transfer/internal/interactor/adapters/telegram"
	"github.com/de1phin/music-transfer/internal/interactor/interactors/console"
	"github.com/de1phin/music-transfer/internal/interactor/interactors/telegram"
	logger "github.com/de1phin/music-transfer/internal/log/file_logger"
	"github.com/de1phin/music-transfer/internal/mux"
	"github.com/de1phin/music-transfer/internal/server"
	"github.com/de1phin/music-transfer/internal/service/mock"
	"github.com/de1phin/music-transfer/internal/service/spotify"
	"github.com/de1phin/music-transfer/internal/service/yandex"
	"github.com/de1phin/music-transfer/internal/service/youtube"
	"github.com/de1phin/music-transfer/internal/storage/cache"
	"github.com/de1phin/music-transfer/internal/storage/postgres"
)

func main() {

	config, err := config.ReadConfig("./config", "config", "yaml")
	if err != nil {
		panic(err)
	}

	fileLogger, err := logger.NewFileLogger("./log/a.log")
	if err != nil {
		panic("FileLogger init error: " + err.Error())
	}
	psql, err := postgres.NewPostgresDatabase(config.Postgres)
	if err != nil {
		panic(err)
	}

	server := server.NewServer(config.Server)

	spotifyStorage := postgres.NewTable[int64, spotifyAPI.Credentials](psql, "Spotify", "id")
	spotifyAPI := spotifyAPI.NewSpotifyAPI(config.Spotify, fileLogger)
	spotify := spotify.NewSpotifyService(config.Spotify, spotifyAPI, spotifyStorage)
	spotifyAPI.BindHandler(server.ServeMux, spotify.OnGetTokens)

	youtubeStorage := postgres.NewTable[int64, youtubeAPI.Credentials](psql, "Youtube", "id")
	youtubeAPI := youtubeAPI.NewYoutubeAPI(config.Youtube, fileLogger)
	youtube := youtube.NewYouTubeService(youtubeAPI, youtubeStorage, fileLogger)
	youtubeAPI.BindHandler(server.ServeMux, youtube.OnGetTokens)

	yandexStorage := postgres.NewTable[int64, yandexAPI.Credentials](psql, "Yandex", "id")
	yandexAPI := yandexAPI.NewYandexAPI(fileLogger, config.Yandex)
	yandex := yandex.NewYandexService(yandexAPI, yandexStorage, fileLogger)
	yandexAPI.BindOnGetCredentials(yandex.OnGetCredentials)

	services := []mux.Service{
		spotify,
		youtube,
		yandex,
		mock.NewMockService(),
	}

	userStateStorage := cache.NewCacheStorage[int64, mux.UserState]()

	telegram, err := telegram.NewTelegramBot(config.Telegram)
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

	mux := mux.NewMux(services, interactors, transferStorage, idStorage, fileLogger)

	go server.Run()

	muxQuit := make(chan struct{})
	go mux.Run(muxQuit)

	c := make(chan os.Signal)
	<-c
	muxQuit <- struct{}{}

	fileLogger.Close()
}
