package spotify

import (
	"github.com/de1phin/music-transfer/internal/api/spotify"
	spotifyAPI "github.com/de1phin/music-transfer/internal/api/spotify"
	"github.com/de1phin/music-transfer/internal/mux"
	"github.com/de1phin/music-transfer/internal/storage"
)

type SpotifyConfig struct {
	Client spotify.Client
	Scopes string
}

type spotifyService struct {
	scopes       string
	client       spotify.Client
	api          *spotify.SpotifyAPI
	redirectURI  string
	tokenStorage storage.Storage[int64, spotify.Credentials]
}

func NewSpotifyService(config SpotifyConfig, redirectURI string, spotifyAPI *spotify.SpotifyAPI, tokenStorage storage.Storage[int64, spotify.Credentials]) *spotifyService {
	return &spotifyService{
		scopes:       config.Scopes,
		client:       config.Client,
		api:          spotifyAPI,
		redirectURI:  redirectURI,
		tokenStorage: tokenStorage,
	}
}

func (spotify *spotifyService) Name() string {
	return "spotify"
}

func (spotify *spotifyService) GetLiked(userID int64) (mux.Playlist, error) {
	tokens, err := spotify.tokenStorage.Get(userID)
	if err != nil {
		return mux.Playlist{}, err
	}
	liked := mux.Playlist{}
	playlist, err := spotify.api.GetLiked(tokens)
	if err != nil {
		return liked, err
	}
	liked.Title = playlist.Name
	for _, track := range playlist.Tracks.Items {
		artists := ""
		for _, artist := range track.Track.Artists {
			artists += artist.Name + " "
		}
		artists = artists[:len(artists)-1]
		liked.Songs = append(liked.Songs, mux.Song{
			Title:   track.Track.Name,
			Artists: artists,
		})
	}

	return liked, nil
}

func (spotify *spotifyService) AddLiked(userID int64, liked mux.Playlist) error {
	tokens, err := spotify.tokenStorage.Get(userID)
	if err != nil {
		return err
	}
	tracks := make([]spotifyAPI.Track, 0)
	for _, track := range liked.Songs {
		search, err := spotify.api.SearchTrack(tokens, track.Title, track.Artists)
		if err != nil {
			return err
		}
		if len(search) == 0 {
			continue
		}
		tracks = append(tracks, search[0])
	}
	return spotify.api.LikeTracks(tokens, tracks)
}

func (spotify *spotifyService) GetPlaylists(userID int64) ([]mux.Playlist, error) {
	tokens, err := spotify.tokenStorage.Get(userID)
	if err != nil {
		return nil, err
	}
	spotifyPlaylists, err := spotify.api.GetUserPlaylists(tokens)
	if err != nil {
		return nil, err
	}
	playlists := []mux.Playlist{}
	for _, playlist := range spotifyPlaylists {
		tracks, err := spotify.api.GetPlaylistTracks(tokens, playlist.ID)
		if err != nil {
			return playlists, err
		}
		muxSongs := make([]mux.Song, len(tracks))
		for i := range tracks {
			muxSongs[i].Title = tracks[i].Track.Name
			artists := ""
			for _, artist := range tracks[i].Track.Artists {
				artists += artist.Name + " "
			}
			artists = artists[:len(artists)-1]
			muxSongs[i].Artists = artists
		}
		playlists = append(playlists, mux.Playlist{
			Title: playlist.Name,
			Songs: muxSongs,
		})
	}
	return playlists, nil
}

func (spotify *spotifyService) AddPlaylists(userID int64, playlists []mux.Playlist) error {
	tokens, err := spotify.tokenStorage.Get(userID)
	if err != nil {
		return err
	}
	userPlaylists, err := spotify.api.GetUserPlaylists(tokens)
	if err != nil {
		return err
	}

	for _, playlist := range playlists {
		playlistID := ""
		for _, userPlaylist := range userPlaylists {
			if userPlaylist.Name == playlist.Title {
				playlistID = userPlaylist.ID
				break
			}
		}
		if playlistID == "" {
			playlist, err := spotify.api.CreatePlaylist(tokens, playlist.Title)
			if err != nil {
				return err
			}
			playlistID = playlist.ID
		}

		tracks := make([]spotifyAPI.Track, 0)
		for _, song := range playlist.Songs {
			search, err := spotify.api.SearchTrack(tokens, song.Title, song.Artists)
			if err != nil {
				return err
			}
			if len(search) == 0 {
				continue
			}
			tracks = append(tracks, search[0])
		}
		spotify.api.AddToPlaylist(tokens, playlistID, tracks)
	}

	return nil
}
