package main

import (
	cache "github.com/de1phin/music-transfer/internal/cache_storage"
	"github.com/de1phin/music-transfer/internal/config"
	"github.com/de1phin/music-transfer/internal/console"
	mockmusicservice "github.com/de1phin/music-transfer/internal/mock_music_service"
	"github.com/de1phin/music-transfer/internal/transfer"
)

func main() {
	storage := cache.NewCacheStorage()
	var services []transfer.MusicService
	services = append(services, mockmusicservice.NewMockMusicService())
	storage.AddService(services[0].Name())
	interactor := console.NewConsoleInteractor(0)
	config := config.NewConfig()

	transfer := transfer.Transfer{
		Interactor: interactor,
		Storage:    storage,
		Services:   services,
		Config:     config,
	}

	transfer.Run()

}
