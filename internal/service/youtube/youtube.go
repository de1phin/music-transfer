package youtube

import (
	"github.com/de1phin/music-transfer/internal/api/youtube"
	"github.com/de1phin/music-transfer/internal/mux"
	"github.com/de1phin/music-transfer/internal/storage"
)

type youtubeService struct {
	tokenStorage storage.Storage[youtube.Credentials]
	config       *youtube.YoutubeConfig
	api          *youtube.YoutubeAPI
}

func NewYouTubeService(api *youtube.YoutubeAPI, tokenStorage storage.Storage[youtube.Credentials], config *youtube.YoutubeConfig) *youtubeService {
	return &youtubeService{
		tokenStorage: tokenStorage,
		config:       config,
		api:          api,
	}
}

func (*youtubeService) Name() string {
	return "youtube"
}

func (yt *youtubeService) GetLiked(userID int64) (liked mux.Playlist) {
	tokens := yt.tokenStorage.Get(userID)
	videos := yt.api.GetLiked(tokens)
	for _, video := range videos {
		liked.Songs = append(liked.Songs, mux.Song{
			Title:   video.Snippet.Title,
			Artists: video.Snippet.ChannelTitle,
		})
	}
	return liked
}

// TODO
func (yt *youtubeService) AddLiked(userID int64, liked mux.Playlist) {
}

// TODO
func (yt *youtubeService) GetPlaylists(userID int64) (playlists []mux.Playlist) {
	return playlists
}

// TODO
func (yt *youtubeService) AddPlaylists(userID int64, playlists []mux.Playlist) {

}
