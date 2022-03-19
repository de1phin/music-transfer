package spotify

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/de1phin/music-transfer/internal/transfer"
)

func (spotify *spotifyService) AddFavourites(credentials interface{}, playlist transfer.Playlist) {
	log.Println("[Spotify] Asked to add", playlist)
}

func (spotify *spotifyService) GetFavourites(tokens interface{}) (playlist transfer.Playlist) {
	limit := 50
	offset := 0

	for {
		url := fmt.Sprintf("https://api.spotify.com/v1/me/tracks?limit=%d&offset=%d", limit, offset)
		request, _ := http.NewRequest("GET", url, nil)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+(tokens.(credentials)).AccessToken)

		client := http.Client{}
		response, _ := client.Do(request)

		spotifyPlaylist := tracksResponse{}
		body, _ := io.ReadAll(response.Body)
		log.Println(string(body))
		json.Unmarshal(body, &spotifyPlaylist)

		log.Println("Got", spotifyPlaylist)

		for _, song := range spotifyPlaylist.Items {
			playlist.Songs = append(playlist.Songs, transfer.Song{
				Name: song.Track.Name,
			})
		}

		offset += len(spotifyPlaylist.Items)

		if len(spotifyPlaylist.Items) < limit {
			break
		}
	}

	return playlist
}
