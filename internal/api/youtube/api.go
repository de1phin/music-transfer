package youtube

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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
	log.Println("VideoID =", videoID)
	url := "https://www.googleapis.com/youtube/v3/videos/rate?rating=like&id=" + videoID
	request, _ := http.NewRequest("POST", url, nil)
	request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	resp, _ := api.httpClient.Do(request)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("LikeVideoResponse:", string(body))
}

func (api *YoutubeAPI) SearchVideos(tokens Credentials, title string, channel string) []SearchResult {
	url := "https://www.googleapis.com/youtube/v3/search?type=video&part=snippet&q=" + url.QueryEscape(title+" "+channel)
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	request.Header.Add("Accept", "application/json")

	response, _ := api.httpClient.Do(request)
	body, _ := ioutil.ReadAll(response.Body)
	searchList := searchListResponse{}
	log.Println("Search:", string(body))
	json.Unmarshal(body, &searchList)

	return searchList.Items
}

func (api *YoutubeAPI) CreatePlaylist(tokens Credentials, title string) Playlist {
	url := "https://www.googleapis.com/youtube/v3/playlists?part=id,snippet"
	data := []byte(`{"snippet":{"title":` + title + "}}")
	request, _ := http.NewRequest("POST", url, bytes.NewReader([]byte(data)))
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
