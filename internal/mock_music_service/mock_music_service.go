package mockmusicservice

import (
	"log"

	"github.com/de1phin/music-transfer/internal/transfer"
)

type mockMusicService struct {
}

func NewMockMusicService() *mockMusicService {
	return &mockMusicService{}
}

func (service *mockMusicService) Name() string {
	return "Mock"
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

func (service *mockMusicService) InitCallbackServer(string) (string, bool) {
	return "", false
}
