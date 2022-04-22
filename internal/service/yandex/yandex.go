package yandex

import (
	"errors"

	"github.com/de1phin/music-transfer/internal/api/yandex"
	"github.com/de1phin/music-transfer/internal/log"
	"github.com/de1phin/music-transfer/internal/mux"
	"github.com/de1phin/music-transfer/internal/storage"
)

const likedPlaylistID = 3

type Yandex struct {
	api     *yandex.YandexAPI
	storage storage.Storage[int64, yandex.Credentials]
	logger  log.Logger
}

func NewYandexService(api *yandex.YandexAPI, storage storage.Storage[int64, yandex.Credentials], logger log.Logger) *Yandex {
	return &Yandex{
		api:     api,
		storage: storage,
		logger:  logger,
	}
}

func (*Yandex) Name() string {
	return "yandex"
}

func (ya *Yandex) GetLiked(userID int64) (*mux.Playlist, error) {
	credentials, err := ya.storage.Get(userID)
	if err != nil {
		return nil, err
	}
	liked, err := ya.api.GetPlaylist(likedPlaylistID, &credentials)
	if err != nil {
		return nil, err
	}
	result := &mux.Playlist{}
	result.Title = liked.Title
	for _, track := range liked.Tracks {
		if track.Type != "music" {
			continue
		}
		song := mux.Song{
			Title: track.Title,
		}
		artists := ""
		for i, a := range track.Artists {
			if i > 0 {
				artists += " "
			}
			artists += a.Name
		}
		song.Artists = artists
		result.Songs = append(result.Songs, song)
	}

	return result, nil
}

func (ya *Yandex) GetPlaylists(userID int64) ([]mux.Playlist, error) {
	credentials, err := ya.storage.Get(userID)
	if err != nil {
		return nil, err
	}
	library, err := ya.api.GetLibrary(&credentials)
	if err != nil {
		return nil, err
	}
	if library == nil {
		return nil, errors.New("YandexAPI.GetPlaylists: Empty library returned")
	}

	result := make([]mux.Playlist, 0)
	for _, id := range library.PlaylistIDs {
		if id == likedPlaylistID {
			continue
		}

		yaPlaylist, err := ya.api.GetPlaylist(id, &credentials)
		if err != nil {
			return nil, err
		}
		playlist := mux.Playlist{
			Title: yaPlaylist.Title,
		}
		for _, track := range yaPlaylist.Tracks {
			artists := ""
			for i, a := range track.Artists {
				if i > 0 {
					artists += " "
				}
				artists += a.Name
			}
			playlist.Songs = append(playlist.Songs, mux.Song{
				Title:   track.Title,
				Artists: artists,
			})
		}
		result = append(result, playlist)
	}

	return result, nil
}

func (ya *Yandex) AddLiked(userID int64, liked mux.Playlist) error {
	credentials, err := ya.storage.Get(userID)
	if err != nil {
		return err
	}
	csrf, err := ya.api.GetAuthCSRF(&credentials)
	if err != nil {
		return err
	}

	for _, song := range liked.Songs {
		track, err := ya.api.SearchTrack(song.Title, song.Artists)
		if err != nil {
			return err
		}
		if track == nil {
			continue
		}
		err = ya.api.LikeTrack(track, &credentials, csrf)
		if err != nil {
			return err
		}
	}

	return nil
}
