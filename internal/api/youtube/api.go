package youtube

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/de1phin/music-transfer/internal/log"
)

type YoutubeAPI struct {
	httpClient *http.Client
	config     *YoutubeConfig
	logger     log.Logger
}

func NewYoutubeAPI(config *YoutubeConfig, logger log.Logger) *YoutubeAPI {
	return &YoutubeAPI{
		config:     config,
		logger:     logger,
		httpClient: &http.Client{},
	}
}

func (api *YoutubeAPI) GetLiked(tokens Credentials) (videos []Video, err error) {
	limit := 50

	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?myRating=like&part=id,snippet&maxResults=%d", limit)

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
		if response.StatusCode != http.StatusOK {
			return videos, errors.New("YoutubeAPI.GetLiked Error: Response Status: " + response.Status)
		}
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return videos, err
		}

		videoList := videoListResponse{}
		err = json.Unmarshal(body, &videoList)
		if err != nil {
			return videos, err
		}

		videos = append(videos, videoList.Items...)

		if videoList.NextPageToken == "" {
			break
		}

		url = fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?myRating=like&part=id,snippet&maxResults=%d&pageToken=%s", limit, videoList.NextPageToken)
	}

	return videos, nil
}

func (api *YoutubeAPI) GetUserPlaylists(tokens Credentials) (playlists []Playlist, err error) {
	url := fmt.Sprintf("https://youtube.googleapis.com/youtube/v3/playlists?part=id,snippet&mine=true&key=%s", api.config.APIKey)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return playlists, err
	}
	request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	request.Header.Add("Accept", "application/json")

	response, err := api.httpClient.Do(request)
	if err != nil {
		return playlists, err
	}
	if response.StatusCode != http.StatusOK {
		return playlists, errors.New("YoutubeAPI.GetUserPlaylists: Response Status: " + response.Status)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return playlists, err
	}
	playlistResponse := playlistListResponse{}
	err = json.Unmarshal(body, &playlistResponse)
	if err != nil {
		return playlists, err
	}

	playlists = playlistResponse.Items
	return playlists, err
}

func (api *YoutubeAPI) GetPlaylistContent(tokens Credentials, playlistID string) (playlists []PlaylistItem, err error) {
	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/playlistItems?part=id,snippet&playlistId=%s", playlistID)

	for {
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return playlists, err
		}
		request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
		request.Header.Add("Accept", "application/json")

		response, err := api.httpClient.Do(request)
		if err != nil {
			return playlists, err
		}
		if response.StatusCode != http.StatusOK {
			return playlists, errors.New("YoutubeAPI.GetPlaylistContent: Response Status: " + response.Status)
		}
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return playlists, err
		}

		itemList := playlistItemListResponse{}
		err = json.Unmarshal(body, &itemList)
		if err != nil {
			return playlists, err
		}

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

func (api *YoutubeAPI) CreatePlaylist(tokens Credentials, title string) (playlist Playlist, err error) {
	data := `{ "snippet": { "title": "` + title + `" } }`
	url := "https://www.googleapis.com/youtube/v3/playlists?part=id,snippet"
	request, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return playlist, err
	}
	request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")

	response, err := api.httpClient.Do(request)
	if err != nil {
		return playlist, err
	}
	if response.StatusCode != http.StatusOK {
		return playlist, errors.New("YoutubeAPI.CreatePlaylist: Response Status: " + response.Status)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return playlist, err
	}
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
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNoContent {
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
	if response.StatusCode == http.StatusUnauthorized {
		return false, nil
	}
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNoContent {
		return false, errors.New("YoutubeAPI.Authorized: Response Status: " + response.Status)
	}
	return true, nil
}
