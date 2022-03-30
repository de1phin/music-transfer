package mockmusicservice

import (
	"log"

	"github.com/de1phin/music-transfer/internal/transfer"
)

type mockMusicService struct {
	callbackURL string
}

func NewMockMusicService(callbackURL string) *mockMusicService {
	return &mockMusicService{callbackURL}
}

func (service *mockMusicService) Name() string {
	return "Mock"
}

func (service *mockMusicService) URLName() string {
	return "mock"
}

func (service *mockMusicService) GetFavourites(interface{}) transfer.Playlist {
	return transfer.Playlist{Name: "abobus", Songs: []transfer.Song{{Name: "bibik"}}}
}

func (service *mockMusicService) AddFavourites(credentials interface{}, playlist transfer.Playlist) {
	log.Println("[mock] Asked to add", playlist)
}

func (service *mockMusicService) GetPlaylists(interface{}) []transfer.Playlist {
	return nil
}

func (service *mockMusicService) AddPlaylists(_ interface{}, playlists []transfer.Playlist) {
	log.Println("[mock] Asked to add", playlists)
}
