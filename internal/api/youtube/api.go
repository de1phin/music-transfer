package youtube

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type YoutubeAPI struct {
	httpClient *http.Client
	config     *YoutubeConfig
}

func NewYoutubeAPI(config *YoutubeConfig) *YoutubeAPI {
	return &YoutubeAPI{
		config:     config,
		httpClient: &http.Client{},
	}
}

func (api *YoutubeAPI) GetLiked(tokens Credentials) (videos []Video) {
	limit := 50

	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?myRating=like&part=id,snippet&maxResults=%d", limit)

	for {
		request, _ := http.NewRequest("GET", url, nil)
		request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)

		response, _ := api.httpClient.Do(request)
		body, _ := ioutil.ReadAll(response.Body)

		videoList := videoListResponse{}
		json.Unmarshal(body, &videoList)

		videos = append(videos, videoList.Items...)

		if videoList.NextPageToken == "" {
			break
		}

		url = fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?myRating=like&part=id,snippet&maxResults=%d&pageToken=%s", limit, videoList.NextPageToken)
	}

	return videos
}

func (api *YoutubeAPI) GetUserPlaylists(tokens Credentials) []Playlist {
	url := fmt.Sprintf("https://youtube.googleapis.com/youtube/v3/playlists?part=id,snippet&mine=true&key=%s", api.config.APIKey)
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	request.Header.Add("Accept", "application/json")

	response, _ := api.httpClient.Do(request)
	body, _ := ioutil.ReadAll(response.Body)
	playlists := playlistListResponse{}
	json.Unmarshal(body, &playlists)

	return playlists.Items
}

func (api *YoutubeAPI) GetPlaylistContent(tokens Credentials, playlistID string) (playlists []PlaylistItem) {
	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/playlistItems?part=id,snippet&playlistId=%s", playlistID)

	for {
		request, _ := http.NewRequest("GET", url, nil)
		request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
		request.Header.Add("Accept", "application/json")

		response, _ := api.httpClient.Do(request)
		body, _ := ioutil.ReadAll(response.Body)
		itemList := playlistItemListResponse{}
		json.Unmarshal(body, &itemList)
		playlists = append(playlists, itemList.Items...)

		if itemList.NextPageToken == "" {
			break
		}
		url = fmt.Sprintf("https://www.googleapis.com/youtube/v3/playlistItems?part=id,snippet&playlistId=%s&pageToken=%s", playlistID, itemList.NextPageToken)
	}

	return playlists

}

func (api *YoutubeAPI) LikeVideo(tokens Credentials, videoID string) {
	url := "https://www.googleapis.com/youtube/v3/videos/rate?rating=like&id=" + videoID
	request, _ := http.NewRequest("POST", url, nil)
	request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	api.httpClient.Do(request)
}

func (api *YoutubeAPI) SearchVideo(title string, artists string) string {
	query := "https://youtube.com/results?search_query=" + strings.Map(func(r rune) rune {
		if r == ' ' {
			return '+'
		} else {
			return r
		}
	}, url.QueryEscape(title+"+"+artists))
	response, _ := http.Get(query)
	body, _ := ioutil.ReadAll(response.Body)
	idx := bytes.Index(body, []byte("/watch?v="))
	idx2 := bytes.Index(body[idx:], []byte("\""))
	return string(body[idx+9 : idx+idx2])
}

func (api *YoutubeAPI) CreatePlaylist(tokens Credentials, title string) Playlist {
	data := `{ "snippet": { "title": "` + title + `" } }`
	url := "https://www.googleapis.com/youtube/v3/playlists?part=id,snippet"
	log.Println("Create", title)
	request, _ := http.NewRequest("POST", url, strings.NewReader(data))
	log.Println("Data:", data)
	request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")

	response, _ := api.httpClient.Do(request)
	body, _ := ioutil.ReadAll(response.Body)
	log.Println("CreatePlaylistResponse:", string(body))
	playlist := Playlist{}
	json.Unmarshal(body, &playlist)

	return playlist
}

func (api *YoutubeAPI) AddToPlaylist(tokens Credentials, playlistID string, videoID string) {
	url := "https://www.googleapis.com/youtube/v3/playlistItems?part=snippet"
	snippets := addVideosRequest{
		Snippet: snippet{
			PlaylistID: playlistID,
			ResourceID: resourceID{
				Kind:    "youtube#video",
				VideoID: videoID,
			},
		},
	}
	data, _ := json.Marshal(snippets)

	request, _ := http.NewRequest("POST", url, bytes.NewReader(data))
	request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")

	api.httpClient.Do(request)

}

func (api *YoutubeAPI) Authorized(tokens Credentials) bool {
	url := "https://www.googleapis.com/youtube/v3/channels?part=id&mine=true&key=" + api.config.APIKey
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)

	response, _ := api.httpClient.Do(request)
	return response.StatusCode != 401
}
