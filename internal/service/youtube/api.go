package youtube

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/de1phin/music-transfer/internal/transfer"
)

func (youtube *YouTubeService) GetFavourites(tokens interface{}) transfer.Playlist {
	limit := 50

	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?myRating=like&part=id,snippet&maxResults=%d", limit)
	playlist := transfer.Playlist{}
	client := &http.Client{}

	for {
		request, _ := http.NewRequest("GET", url, nil)
		request.Header.Add("Authorization", "Bearer "+tokens.(credentials).AccessToken)

		response, _ := client.Do(request)
		body, _ := ioutil.ReadAll(response.Body)

		videos := videoListResponse{}
		log.Println(string(body))
		json.Unmarshal(body, &videos)

		for _, video := range videos.Items {
			playlist.Songs = append(playlist.Songs, transfer.Song{Name: video.Snippet.Title, Artists: video.Snippet.ChannelTitle})
		}

		if videos.NextPageToken == "" {
			break
		}

		url = fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?myRating=like&part=id,snippet&maxResults=%d&pageToken=%s", limit, videos.NextPageToken)
	}

	return playlist
}

func (youtube *YouTubeService) AddFavourites(tokens interface{}, playlist transfer.Playlist) {
	log.Println("[YouTube] asked to add:", playlist)
}

func (youtube *YouTubeService) GetPlaylists(tokens interface{}) []transfer.Playlist {
	url := fmt.Sprintf("https://youtube.googleapis.com/youtube/v3/playlists?mine=true&key=%s", youtube.apiKey)
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("Authorization", "Bearer "+tokens.(credentials).AccessToken)
	request.Header.Add("Accept", "application/json")

	client := &http.Client{}
	response, _ := client.Do(request)
	body, _ := ioutil.ReadAll(response.Body)
	playlists := playlistListResponse{}
	json.Unmarshal(body, &playlists)

	log.Println("Playlists:", playlists)

	return nil
}

func (youtube *YouTubeService) AddPlaylists(tokens interface{}, playlists []transfer.Playlist) {
	log.Println("[YouTube] asked to add:", playlists)
}
