package youtube

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
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

func (api *YoutubeAPI) GetLiked(tokens Credentials) ([]Video, error) {
	limit := 50

	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?myRating=like&part=id,snippet&maxResults=%d", limit)

	videos := []Video{}

	for {
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return videos, err
		}
		request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)

		response, err := api.httpClient.Do(request)
		if err != nil {
			return videos, err
		}
		if response.StatusCode != 200 {
			return videos, errors.New("YoutubeAPI.GetLiked Error: Response Status: " + response.Status)
		}
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return videos, err
		}

		videoList := videoListResponse{}
		json.Unmarshal(body, &videoList)

		videos = append(videos, videoList.Items...)

		if videoList.NextPageToken == "" {
			break
		}

		url = fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?myRating=like&part=id,snippet&maxResults=%d&pageToken=%s", limit, videoList.NextPageToken)
	}

	return videos, nil
}

func (api *YoutubeAPI) GetUserPlaylists(tokens Credentials) ([]Playlist, error) {
	url := fmt.Sprintf("https://youtube.googleapis.com/youtube/v3/playlists?part=id,snippet&mine=true&key=%s", api.config.APIKey)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	request.Header.Add("Accept", "application/json")

	response, err := api.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New("YoutubeAPI.GetUserPlaylists: Response Status: " + response.Status)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	playlists := playlistListResponse{}
	json.Unmarshal(body, &playlists)

	return playlists.Items, err
}

func (api *YoutubeAPI) GetPlaylistContent(tokens Credentials, playlistID string) ([]PlaylistItem, error) {
	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/playlistItems?part=id,snippet&playlistId=%s", playlistID)

	playlists := []PlaylistItem{}

	for {
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
		request.Header.Add("Accept", "application/json")

		response, err := api.httpClient.Do(request)
		if err != nil {
			return nil, err
		}
		if response.StatusCode != 200 {
			return nil, errors.New("YoutubeAPI.GetPlaylistContent: Response Status: " + response.Status)
		}
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		itemList := playlistItemListResponse{}
		json.Unmarshal(body, &itemList)
		playlists = append(playlists, itemList.Items...)

		if itemList.NextPageToken == "" {
			break
		}
		url = fmt.Sprintf("https://www.googleapis.com/youtube/v3/playlistItems?part=id,snippet&playlistId=%s&pageToken=%s", playlistID, itemList.NextPageToken)
	}

	return playlists, nil
}

func (api *YoutubeAPI) LikeVideo(tokens Credentials, videoID string) error {
	url := "https://www.googleapis.com/youtube/v3/videos/rate?rating=like&id=" + videoID
	request, _ := http.NewRequest("POST", url, nil)
	request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	response, err := api.httpClient.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNoContent {
		return errors.New("YoutubeAPI.LikeVideo: Response Status: " + response.Status)
	}
	return nil
}

func (api *YoutubeAPI) CreatePlaylist(tokens Credentials, title string) (Playlist, error) {
	data := `{ "snippet": { "title": "` + title + `" } }`
	url := "https://www.googleapis.com/youtube/v3/playlists?part=id,snippet"
	request, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return Playlist{}, err
	}
	request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")

	response, err := api.httpClient.Do(request)
	if err != nil {
		return Playlist{}, err
	}
	if response.StatusCode != 200 {
		return Playlist{}, errors.New("YoutubeAPI.CreatePlaylist: Response Status: " + response.Status)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return Playlist{}, err
	}
	playlist := Playlist{}
	err = json.Unmarshal(body, &playlist)

	return playlist, err
}

func (api *YoutubeAPI) AddToPlaylist(tokens Credentials, playlistID string, videoID string) error {
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
	data, err := json.Marshal(snippets)
	if err != nil {
		return err
	}

	request, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")

	response, err := api.httpClient.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return errors.New("YoutubeAPI.AddToPlaylist: Response Status: " + response.Status)
	}

	return nil
}

func (api *YoutubeAPI) Authorized(tokens Credentials) (bool, error) {
	url := "https://www.googleapis.com/youtube/v3/channels?part=id&mine=true&key=" + api.config.APIKey
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)

	response, err := api.httpClient.Do(request)
	if err != nil {
		return false, err
	}
	if response.StatusCode == 401 {
		return false, nil
	}
	if response.StatusCode != 200 {
		return false, errors.New("YoutubeAPI.Authorized: Response Status: " + response.Status)
	}
	return true, nil
}
