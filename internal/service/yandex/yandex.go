package yandex

import (
	"errors"
	"time"

	"github.com/de1phin/music-transfer/internal/api/yandex"
	"github.com/de1phin/music-transfer/internal/log"
	"github.com/de1phin/music-transfer/internal/mux"
	"github.com/de1phin/music-transfer/internal/storage"
)

const (
	likedPlaylistID = 3
	requestDelayMs  = 1700
)

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

func (ya *Yandex) GetLiked(userID int64) (playlist mux.Playlist, err error) {
	credentials, err := ya.storage.Get(userID)
	if err != nil {
		return playlist, err
	}
	liked, err := ya.api.GetPlaylist(likedPlaylistID, credentials)
	if err != nil {
		return playlist, err
	}
	playlist.Title = liked.Title
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
		playlist.Songs = append(playlist.Songs, song)
	}

	return playlist, nil
}

func (ya *Yandex) GetPlaylists(userID int64) (playlist []mux.Playlist, err error) {
	credentials, err := ya.storage.Get(userID)
	if err != nil {
		return playlist, err
	}
	library, err := ya.api.GetLibrary(credentials)
	if err != nil {
		return playlist, err
	}

	for _, id := range library.PlaylistIDs {
		if id == likedPlaylistID {
			continue
		}

		yaPlaylist, err := ya.api.GetPlaylist(id, credentials)
		if err != nil {
			return nil, err
		}
		p := mux.Playlist{
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
			p.Songs = append(p.Songs, mux.Song{
				Title:   track.Title,
				Artists: artists,
			})
		}
		playlist = append(playlist, p)
	}

	return playlist, nil
}

func (ya *Yandex) AddLiked(userID int64, liked mux.Playlist) error {
	credentials, err := ya.storage.Get(userID)
	if err != nil {
		return err
	}
	authTokens, err := ya.api.GetAuthTokens(credentials)
	if err != nil {
		return err
	}

	for _, song := range liked.Songs {
		time.Sleep(time.Millisecond * requestDelayMs)
		track, err := ya.api.SearchTrack(song.Title, song.Artists)
		if err != nil {
			return err
		}
		time.Sleep(time.Millisecond * requestDelayMs)
		err = ya.api.LikeTrack(track, credentials, authTokens)
		if err != nil {
			ya.logger.Log("Yandex.AddLiked:", err)
		}
	}

	return nil
}

func (ya *Yandex) addPlaylist(playlist mux.Playlist, credentials yandex.Credentials, authTokens yandex.AuthTokens) error {
	playlistSnippet, err := ya.api.AddPlaylist(playlist.Title, credentials, authTokens)
	if err != nil {
		return err
	}

	tracks := []yandex.TrackSnippet{}
	for _, s := range playlist.Songs {
		time.Sleep(time.Millisecond * requestDelayMs)
		t, err := ya.api.SearchTrack(s.Title, s.Artists)
		if err != nil {
			ya.logger.Log(errors.New("Yandex.AddPlaylists:" + err.Error()))
			continue
		}
		if len(t.Albums) == 0 {
			ya.logger.Log(errors.New("Yandex.AddPlaylists: Bad track returned:"), t)
			continue
		}
		tracks = append(tracks, yandex.TrackSnippet{
			ID:      t.ID,
			AlbumID: t.Albums[0].ID,
		})
	}
	time.Sleep(time.Millisecond * requestDelayMs)
	return ya.api.AddToPlaylist(tracks, playlistSnippet, credentials, authTokens)
}

func (ya *Yandex) AddPlaylists(userID int64, playlist []mux.Playlist) error {
	credentials, err := ya.storage.Get(userID)
	if err != nil {
		return err
	}
	authTokens, err := ya.api.GetAuthTokens(credentials)
	if err != nil {
		return err
	}

	for _, p := range playlist {
		err = ya.addPlaylist(p, credentials, authTokens)
		if err != nil {
			ya.logger.Log("Yandex.AddPlaylists:", err)
			continue
		}
	}

	return nil
}
