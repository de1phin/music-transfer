package mock

import (
	"fmt"
	"log"

	"github.com/de1phin/music-transfer/internal/mux"
)

type Mock struct{}

func NewMockService() *Mock {
	return &Mock{}
}

func (*Mock) Name() string {
	return "mock"
}

func (*Mock) GetAuthURL(userID int64) string {
	return fmt.Sprintf("mock/user_id=%d", userID)
}

func (*Mock) GetLiked(int64) mux.Playlist {
	return mux.Playlist{Title: "Liked", Songs: []mux.Song{
		{
			Title:   "Дед Максим",
			Artists: "RADIO TAPOK",
		},
	}}
}

func (*Mock) AddLiked(userID int64, liked mux.Playlist) {
	log.Println("[mock] Asked to like:", liked)
}

func (*Mock) GetPlaylists(int64) []mux.Playlist {
	return []mux.Playlist{
		{
			Title: "Ded",
			Songs: []mux.Song{
				{
					Title:   "Дед Максим",
					Artists: "RADIO TAPOK",
				},
			},
		},
	}
}

func (*Mock) AddPlaylists(userID int64, playlists []mux.Playlist) {
	log.Println("[mock] Asked to add:", playlists)
}
