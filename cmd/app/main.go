package main

import (
	cache "github.com/de1phin/music-transfer/internal/cache_storage"
	"github.com/de1phin/music-transfer/internal/config"
	"github.com/de1phin/music-transfer/internal/console"
	mockmusicservice "github.com/de1phin/music-transfer/internal/mock_music_service"
	"github.com/de1phin/music-transfer/internal/spotify"
	"github.com/de1phin/music-transfer/internal/transfer"
	"github.com/de1phin/music-transfer/internal/youtube"
)

func main() {
	config := config.NewConfig()
	storage := cache.NewCacheStorage()
	var services []transfer.MusicService
	mockMusicService := mockmusicservice.NewMockMusicService(config.GetCallbackURL())
	storage.AddService(mockMusicService.Name())
	services = append(services, mockMusicService)
	spotify := spotify.NewSpotifyService(config.GetSpotifyClientID(), config.GetSpotifyClientSecret(), config.GetSpotifyScope(), config.GetCallbackURL())
	storage.AddService(spotify.Name())
	services = append(services, spotify)
	youtube := youtube.NewYouTubeService(config.GetCallbackURL(), config.GetYouTubeScope(), config.GetYouTubeApiKEY(), config.GetYouTubeClientID(), config.GetYouTubeClientSecret())
	storage.AddService(youtube.Name())
	services = append(services, youtube)
	interactor := console.NewConsoleInteractor(17)

	transfer := transfer.Transfer{
		Interactor: interactor,
		Storage:    storage,
		Services:   services,
		Config:     config,
	}

	transfer.Run()

}
