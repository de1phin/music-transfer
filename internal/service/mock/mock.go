package mock

import (
	"fmt"

	"github.com/de1phin/music-transfer/internal/mux"
)

type Mock struct{}

func NewMockService() *Mock {
	return &Mock{}
}

func (*Mock) Name() string {
	return "mock"
}

func (*Mock) GetAuthURL(userID int64) (string, error) {
	return "", fmt.Errorf("Mock: unreachable")
}

func (*Mock) GetLiked(int64) (mux.Playlist, error) {
	return mux.Playlist{Title: "Liked", Songs: []mux.Song{
		{
			Title:   "Дед Максим",
			Artists: "RADIO TAPOK",
		},
	}}, nil
}

func (*Mock) AddLiked(userID int64, liked mux.Playlist) error {
	fmt.Println("[mock] Asked to like:", liked)
	return nil
}

func (*Mock) GetPlaylists(int64) ([]mux.Playlist, error) {
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
	}, nil
}

func (*Mock) AddPlaylists(userID int64, playlists []mux.Playlist) error {
	fmt.Println("[mock] Asked to add:", playlists)
	return nil
}

func (*Mock) Authorized(int64) (bool, error) {
	return true, nil
}
