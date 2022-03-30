package spotify

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/de1phin/music-transfer/internal/transfer"
)

func (spotify *spotifyService) AddFavourites(tokens interface{}, playlist transfer.Playlist) {
	log.Println("[Spotify] Asked to add", playlist)
}

func (spotify *spotifyService) GetFavourites(tokens interface{}) (playlist transfer.Playlist) {
	limit := 50
	offset := 0

	client := &http.Client{}
	for {
		url := fmt.Sprintf("https://api.spotify.com/v1/me/tracks?limit=%d&offset=%d", limit, offset)
		request, _ := http.NewRequest("GET", url, nil)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+(tokens.(credentials)).AccessToken)

		response, _ := client.Do(request)

		spotifyPlaylist := tracksResponse{}
		body, _ := io.ReadAll(response.Body)
		json.Unmarshal(body, &spotifyPlaylist)

		for _, song := range spotifyPlaylist.Items {
			artists := ""
			for _, artist := range song.Track.Artists {
				artists += " " + artist.Name
			}
			playlist.Songs = append(playlist.Songs, transfer.Song{
				Name:    song.Track.Name,
				Artists: artists[1:],
			})
		}

		offset += len(spotifyPlaylist.Items)

		if len(spotifyPlaylist.Items) < limit {
			break
		}
	}

	return playlist
}

func (spotify *spotifyService) getSongsFromPlaylist(tokens interface{}, playlistID string) []trackItem {
	client := &http.Client{}
	limit := 50
	offset := 0

	items := make([]trackItem, 0)
	for {
		url := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks?limit=%d&offset=%d&fields=items(track(name,artists(name)))", playlistID, limit, offset)
		request, _ := http.NewRequest("GET", url, nil)
		request.Header.Add("Authorization", "Bearer "+tokens.(credentials).AccessToken)
		request.Header.Add("Content-Type", "application/json")

		response, _ := client.Do(request)

		tracks := tracksResponse{}
		body, _ := ioutil.ReadAll(response.Body)
		json.Unmarshal(body, &tracks)
		items = append(items, tracks.Items...)

		offset += limit
		if len(tracks.Items) < limit {
			break
		}
	}
	return items
}

func (spotify *spotifyService) GetPlaylists(tokens interface{}) []transfer.Playlist {
	limit := 50
	offset := 0

	client := &http.Client{}
	playlists := make([]transfer.Playlist, 0)
	for {
		url := fmt.Sprintf("https://api.spotify.com/v1/me/playlists?limit=%d&offset=%d", limit, offset)
		request, _ := http.NewRequest("GET", url, nil)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+(tokens.(credentials)).AccessToken)

		response, _ := client.Do(request)

		spotifyPlaylists := playlistResponse{}
		body, _ := io.ReadAll(response.Body)
		json.Unmarshal(body, &spotifyPlaylists)

		for _, item := range spotifyPlaylists.Items {
			tracks := spotify.getSongsFromPlaylist(tokens, item.ID)
			songs := make([]transfer.Song, len(tracks))
			for i := range tracks {
				artists := ""
				for _, artist := range tracks[i].Track.Artists {
					artists += " " + artist.Name
				}
				songs[i] = transfer.Song{Name: tracks[i].Track.Name, Artists: artists[1:]}
			}
			playlists = append(playlists, transfer.Playlist{Name: item.Name, Songs: songs})
		}

		offset += limit
		if len(spotifyPlaylists.Items) < limit {
			break
		}
	}
	return playlists
}

func (spotify *spotifyService) AddPlaylists(tokens interface{}, playlists []transfer.Playlist) {

}
