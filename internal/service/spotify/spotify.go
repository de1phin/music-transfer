package spotify

import (
	"log"

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
	scopes      string
	client      spotify.Client
	api         *spotify.SpotifyAPI
	redirectURI string
	storage     storage.Storage[spotify.Credentials]
}

func NewSpotifyService(config SpotifyConfig, redirectURI string, spotifyAPI *spotify.SpotifyAPI, storage storage.Storage[spotify.Credentials]) *spotifyService {
	return &spotifyService{
		scopes:      config.Scopes,
		client:      config.Client,
		api:         spotifyAPI,
		redirectURI: redirectURI,
		storage:     storage,
	}
}

func (spotify *spotifyService) Name() string {
	return "spotify"
}

func (spotify *spotifyService) GetLiked(userID int64) (liked mux.Playlist) {
	tokens := spotify.storage.Get(userID)
	playlist := spotify.api.GetLiked(tokens)
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

	return liked
}

func (spotify *spotifyService) AddLiked(userID int64, liked mux.Playlist) {
	tokens := spotify.storage.Get(userID)
	tracks := make([]spotifyAPI.Track, 0)
	for _, track := range liked.Songs {
		search := spotify.api.SearchTrack(tokens, track.Title, track.Artists)
		if len(search) == 0 {
			continue
		}
		tracks = append(tracks, search[0])
	}
	spotify.api.LikeTracks(tokens, tracks)
}

func (spotify *spotifyService) GetPlaylists(userID int64) (playlists []mux.Playlist) {
	tokens := spotify.storage.Get(userID)
	spotifyPlaylists := spotify.api.GetUserPlaylists(tokens)
	for _, playlist := range spotifyPlaylists {
		tracks := spotify.api.GetPlaylistTracks(tokens, playlist.ID)
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
		log.Println("Total tracks:", len(muxSongs))
		playlists = append(playlists, mux.Playlist{
			Title: playlist.Name,
			Songs: muxSongs,
		})
	}
	return playlists
}

func (spotify *spotifyService) AddPlaylists(userID int64, playlists []mux.Playlist) {
	tokens := spotify.storage.Get(userID)
	userPlaylists := spotify.api.GetUserPlaylists(tokens)

	for _, playlist := range playlists {
		playlistID := ""
		for _, userPlaylist := range userPlaylists {
			if userPlaylist.Name == playlist.Title {
				playlistID = userPlaylist.ID
				break
			}
		}
		if playlistID == "" {
			playlistID = spotify.api.CreatePlaylist(tokens, playlist.Title).ID
		}

		tracks := make([]spotifyAPI.Track, 0)
		for _, song := range playlist.Songs {
			search := spotify.api.SearchTrack(tokens, song.Title, song.Artists)
			if len(search) == 0 {
				continue
			}
			tracks = append(tracks, search[0])
		}
		spotify.api.AddToPlaylist(tokens, playlistID, tracks)
	}
}
