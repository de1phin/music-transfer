package youtube

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func (api *YoutubeAPI) GetLiked(tokens Credentials) (videos []Video, err error) {
	limit := 50

	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?myRating=like&part=id,snippet&maxResults=%d", limit)

	for {
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return videos, fmt.Errorf("Unable to create request: %w", err)
		}
		request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)

		response, err := api.httpClient.Do(request)
		if err != nil {
			return videos, fmt.Errorf("Unable to do request: %w", err)
		}
		if response.StatusCode != http.StatusOK {
			return videos, fmt.Errorf("Bad response status: %s", response.Status)
		}
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return videos, fmt.Errorf("Unable to read body: %w", err)
		}

		videoList := videoListResponse{}
		err = json.Unmarshal(body, &videoList)
		if err != nil {
			return videos, fmt.Errorf("Unable to unmarshal: %w", err)
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
	url := fmt.Sprintf("https://youtube.googleapis.com/youtube/v3/playlists?part=id,snippet&mine=true&key=%s", api.APIKey)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return playlists, fmt.Errorf("Unable to create request: %w", err)
	}
	request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	request.Header.Add("Accept", "application/json")

	response, err := api.httpClient.Do(request)
	if err != nil {
		return playlists, fmt.Errorf("Unable to do request: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return playlists, fmt.Errorf("Bad response status: %s", response.Status)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return playlists, fmt.Errorf("Unable to read body: %w", err)
	}
	playlistResponse := playlistListResponse{}
	err = json.Unmarshal(body, &playlistResponse)
	if err != nil {
		return playlists, fmt.Errorf("Unable to unmarshal: %w", err)
	}

	playlists = playlistResponse.Items
	return playlists, err
}

func (api *YoutubeAPI) GetPlaylistContent(tokens Credentials, playlistID string) (playlists []PlaylistItem, err error) {
	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/playlistItems?part=id,snippet&playlistId=%s", playlistID)

	for {
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return playlists, fmt.Errorf("Unable to create request: %w", err)
		}
		request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
		request.Header.Add("Accept", "application/json")

		response, err := api.httpClient.Do(request)
		if err != nil {
			return playlists, fmt.Errorf("Unable to do request: %w", err)
		}
		if response.StatusCode != http.StatusOK {
			return playlists, fmt.Errorf("Bad response status: %w", response.Status)
		}
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return playlists, fmt.Errorf("Unable to read body: %w", err)
		}

		itemList := playlistItemListResponse{}
		err = json.Unmarshal(body, &itemList)
		if err != nil {
			return playlists, fmt.Errorf("Unable to unmarshal: %w", err)
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
	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("Unable to create request: %w", err)
	}
	request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	response, err := api.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("Unable to do request: %w", err)
	}
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Bad response status: %s", response.Status)
	}
	return nil
}

func (api *YoutubeAPI) CreatePlaylist(tokens Credentials, title string) (playlist Playlist, err error) {
	data := `{ "snippet": { "title": "` + title + `" } }`
	url := "https://www.googleapis.com/youtube/v3/playlists?part=id,snippet"
	request, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return playlist, fmt.Errorf("Unable to create request: %w", err)
	}
	request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")

	response, err := api.httpClient.Do(request)
	if err != nil {
		return playlist, fmt.Errorf("Unable to do request: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return playlist, fmt.Errorf("Bad response status: %s", response.Status)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return playlist, fmt.Errorf("Unable to read body: %w", err)
	}
	err = json.Unmarshal(body, &playlist)

	return playlist, fmt.Errorf("Unable to unmarshal: %w", err)
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
		return fmt.Errorf("Unable to create request: %w", err)
	}
	request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")

	response, err := api.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("Unable to do request: %w", err)
	}
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Bad response status: %s", response.Status)
	}

	return nil
}

func (api *YoutubeAPI) Authorized(tokens Credentials) (bool, error) {
	url := "https://www.googleapis.com/youtube/v3/channels?part=id&mine=true&key=" + api.APIKey
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("Unable to create request: %w", err)
	}
	request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)

	response, err := api.httpClient.Do(request)
	if err != nil {
		return false, fmt.Errorf("Unable to do request: %w", err)
	}
	if response.StatusCode == http.StatusUnauthorized {
		return false, nil
	}
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNoContent {
		return false, fmt.Errorf("Bad response status: %s", response.Status)
	}
	return true, nil
}
