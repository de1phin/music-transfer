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
			Artists: video.Snippet.VideoOwnerChannelTitle,
		})
	}
	return liked
}

func (yt *youtubeService) AddLiked(userID int64, liked mux.Playlist) {
	tokens := yt.tokenStorage.Get(userID)
	for _, song := range liked.Songs {
		videoID := yt.api.SearchVideo(song.Title, song.Artists)
		yt.api.LikeVideo(tokens, videoID)
	}
}

func (yt *youtubeService) GetPlaylists(userID int64) (playlists []mux.Playlist) {
	tokens := yt.tokenStorage.Get(userID)
	ytplaylists := yt.api.GetUserPlaylists(tokens)
	for _, playlist := range ytplaylists {
		videos := yt.api.GetPlaylistContent(tokens, playlist.ID)
		songs := make([]mux.Song, len(videos))
		for i := range videos {
			songs[i] = mux.Song{
				Title:   videos[i].Snippet.Title,
				Artists: videos[i].Snippet.VideoOwnerChannelTitle,
			}
		}
		playlists = append(playlists, mux.Playlist{
			Title: playlist.Snippet.Title,
			Songs: songs,
		})
	}

	return playlists
}

func (yt *youtubeService) AddPlaylists(userID int64, playlists []mux.Playlist) {
	tokens := yt.tokenStorage.Get(userID)
	userPlaylists := yt.api.GetUserPlaylists(tokens)
	for _, playlist := range playlists {
		exist := false
		playlistID := ""
		for _, up := range userPlaylists {
			if playlist.Title == up.Snippet.Title {
				exist = true
				playlistID = up.ID
				break
			}
		}

		if !exist {
			playlistID = yt.api.CreatePlaylist(tokens, playlist.Title).ID
		}

		for _, v := range playlist.Songs {
			videoID := yt.api.SearchVideo(v.Title, v.Artists)
			yt.api.AddToPlaylist(tokens, playlistID, videoID)
		}
	}
}
