package youtube

import (
	"fmt"
	"time"

	"github.com/de1phin/music-transfer/internal/api/youtube"
	"github.com/de1phin/music-transfer/internal/log"
	"github.com/de1phin/music-transfer/internal/mux"
	"github.com/de1phin/music-transfer/internal/storage"
)

type youtubeService struct {
	tokenStorage       storage.Storage[int64, youtube.Credentials]
	api                *youtube.YoutubeAPI
	logger             log.Logger
	OnAuthorizedNotify mux.OnAuthorized
}

func NewYouTubeService(api *youtube.YoutubeAPI, tokenStorage storage.Storage[int64, youtube.Credentials], logger log.Logger) *youtubeService {
	return &youtubeService{
		tokenStorage: tokenStorage,
		api:          api,
		logger:       logger,
	}
}

func (*youtubeService) Name() string {
	return "youtube"
}

func (yt *youtubeService) GetLiked(userID int64) (liked mux.Playlist, err error) {
	tokens, err := yt.tokenStorage.Get(userID)
	if err != nil {
		return liked, fmt.Errorf("Unable to get tokens: %w", err)
	}
	videos, err := yt.api.GetLiked(tokens)
	if err != nil {
		return liked, fmt.Errorf("Unable to get liked: %w", err)
	}
	for _, video := range videos {
		liked.Songs = append(liked.Songs, mux.Song{
			Title:   video.Snippet.Title,
			Artists: video.Snippet.VideoOwnerChannelTitle,
		})
	}
	return liked, nil
}

func (yt *youtubeService) AddLiked(userID int64, liked mux.Playlist) error {
	tokens, err := yt.tokenStorage.Get(userID)
	if err != nil {
		return err
	}
	for _, song := range liked.Songs {
		videoID, err := yt.api.SearchVideo(song.Title, song.Artists)
		if err != nil {
			yt.logger.Error(fmt.Errorf("Youtube: Unable to add liked: Unable to search video: %w", err))
			continue
		}
		err = yt.api.LikeVideo(tokens, videoID)
		if err != nil {
			yt.logger.Error(fmt.Errorf("Youtube: Unable to add liked: Unable to like video: %w", err))
			continue
		}
	}
	return nil
}

func (yt *youtubeService) GetPlaylists(userID int64) (playlists []mux.Playlist, err error) {
	tokens, err := yt.tokenStorage.Get(userID)
	if err != nil {
		return playlists, fmt.Errorf("Unable to get tokens: %w", err)
	}

	ytplaylists, err := yt.api.GetUserPlaylists(tokens)
	if err != nil {
		return playlists, fmt.Errorf("Unable to get user playlists: %w", err)
	}

	for _, playlist := range ytplaylists {
		videos, err := yt.api.GetPlaylistContent(tokens, playlist.ID)
		if err != nil {
			return nil, fmt.Errorf("Unable to get playlist content: %w", err)
		}
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

	return playlists, nil
}

func (yt *youtubeService) AddPlaylists(userID int64, playlists []mux.Playlist) error {
	tokens, err := yt.tokenStorage.Get(userID)
	if err != nil {
		return fmt.Errorf("Unable to get tokens: %w", err)
	}
	userPlaylists, err := yt.api.GetUserPlaylists(tokens)
	if err != nil {
		return fmt.Errorf("Unable to get user playlists: %w", err)
	}
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
			playlist, err := yt.api.CreatePlaylist(tokens, playlist.Title)
			if err != nil {
				return fmt.Errorf("Unable to create playlist: %w", err)
			}
			playlistID = playlist.ID
		}

		for _, v := range playlist.Songs {
			time.Sleep(time.Millisecond * 50)
			videoID, err := yt.api.SearchVideo(v.Title, v.Artists)
			if err != nil {
				yt.logger.Error(fmt.Errorf("Youtube: Unable to add playlists: %w", err))
				continue
			}
			err = yt.api.AddToPlaylist(tokens, playlistID, videoID)
			if err != nil {
				yt.logger.Error(fmt.Errorf("Youtube: Unable to add playlists: %w", err))
				continue
			}
		}
	}
	return nil
}
