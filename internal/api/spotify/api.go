package spotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/de1phin/music-transfer/internal/log"
)

type SpotifyAPI struct {
	httpClient  *http.Client
	client      Client
	logger      log.Logger
	redirectURI string
}

func NewSpotifyAPI(config Config, logger log.Logger) *SpotifyAPI {
	return &SpotifyAPI{
		httpClient: &http.Client{},
		client: Client{
			ID:     config.ClientID,
			Secret: config.ClientSecret,
		},
		logger:      logger,
		redirectURI: config.RedirectURI,
	}
}

func (api *SpotifyAPI) GetLiked(tokens Credentials) (playlist Playlist, err error) {
	limit := 50
	offset := 0

	for {
		url := fmt.Sprintf("https://api.spotify.com/v1/me/tracks?limit=%d&offset=%d", limit, offset)
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return playlist, fmt.Errorf("Unable to create request: %w", err)
		}
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+(tokens.AccessToken))

		response, err := api.httpClient.Do(request)
		if err != nil {
			return playlist, fmt.Errorf("Unable to do request: %w", err)
		}
		if response.StatusCode != http.StatusOK {
			return playlist, fmt.Errorf("Bad response status: %s", response.Status)
		}

		spotifyPlaylist := tracksResponse{}
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return playlist, fmt.Errorf("Unable to read body: %w", err)
		}
		err = json.Unmarshal(body, &spotifyPlaylist)
		if err != nil {
			return playlist, fmt.Errorf("Unable to unmarshal: %w", err)
		}

		for _, item := range spotifyPlaylist.Items {
			playlist.Tracks.Items = append(playlist.Tracks.Items, TrackItem{
				Track: Track{
					Name:    item.Track.Name,
					Artists: item.Track.Artists,
				},
			})
		}

		offset += len(spotifyPlaylist.Items)

		if len(spotifyPlaylist.Items) < limit {
			break
		}
	}

	return playlist, nil
}

func (api *SpotifyAPI) SearchTrack(tokens Credentials, title string, artists string) (results []Track, err error) {
	query := url.QueryEscape(title + " " + artists)
	url := "https://api.spotify.com/v1/search?type=track&q=" + query

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return results, fmt.Errorf("Unable to create request: %w", err)
	}
	request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	request.Header.Add("Content-Type", "application/json")

	response, err := api.httpClient.Do(request)
	if err != nil {
		return results, fmt.Errorf("Unable to do request: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return results, fmt.Errorf("Bad response status: %s", response.Status)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return results, fmt.Errorf("Unable to read body: %w", err)
	}
	searchResponse := searchResponse{}

	err = json.Unmarshal(body, &searchResponse)
	if err != nil {
		return results, fmt.Errorf("Unable to unmarshal: %w", err)
	}

	results = searchResponse.Tracks.Items
	return searchResponse.Tracks.Items, nil
}

func (api *SpotifyAPI) LikeTracks(tokens Credentials, tracks []Track) error {
	limit := 50

	for i := 0; i < len(tracks); i += limit {
		if len(tracks)-i < limit {
			limit = len(tracks) - i
		}
		ids := make([]string, limit)
		for j := i; j < i+limit; j++ {
			ids[j-i] = tracks[j].ID
		}
		url := fmt.Sprintf("https://api.spotify.com/v1/me/tracks?ids=%s", strings.Join(ids, ","))

		request, err := http.NewRequest("PUT", url, nil)
		if err != nil {
			return fmt.Errorf("Unable to create request: %w", err)
		}
		request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
		request.Header.Add("Content-Type", "application/json")

		response, err := api.httpClient.Do(request)
		if err != nil {
			return fmt.Errorf("Unable to do request: %w", err)
		}
		if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNoContent {
			return fmt.Errorf("Bad response status: %s", response.Status)
		}
		if i+limit >= len(tracks) {
			break
		}
	}
	return nil
}

func (api *SpotifyAPI) GetUserPlaylists(tokens Credentials) (playlists []Playlist, err error) {
	limit := 50
	offset := 0

	for {
		url := fmt.Sprintf("https://api.spotify.com/v1/me/playlists?limit=%d&offset=%d", limit, offset)
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return playlists, fmt.Errorf("Unable to create request: %w", err)
		}
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+(tokens.AccessToken))

		response, err := api.httpClient.Do(request)
		if err != nil {
			return playlists, fmt.Errorf("Unable to do request: %w", err)
		}
		if response.StatusCode != http.StatusOK {
			return playlists, fmt.Errorf("Bad response status: %s", response.Status)
		}

		spotifyPlaylists := playlistResponse{}
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return playlists, fmt.Errorf("Unable to read body: %w", err)
		}
		err = json.Unmarshal(body, &spotifyPlaylists)
		if err != nil {
			return playlists, fmt.Errorf("Unable to unmarshal: %w", err)
		}

		playlists = append(playlists, spotifyPlaylists.Items...)

		offset += limit
		if len(spotifyPlaylists.Items) < limit {
			break
		}
	}

	return playlists, nil
}

func (api *SpotifyAPI) GetPlaylistTracks(tokens Credentials, playlistID string) (tracks []TrackItem, err error) {
	limit := 50
	offset := 0

	for {
		url := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks?limit=%d&offset=%d&fields=items(track(name,artists(name)))", playlistID, limit, offset)
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return tracks, fmt.Errorf("Unable to create request: %w", err)
		}
		request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
		request.Header.Add("Content-Type", "application/json")

		response, err := api.httpClient.Do(request)
		if err != nil {
			return tracks, fmt.Errorf("Unable to do request: %w", err)
		}
		if response.StatusCode != http.StatusOK {
			return tracks, fmt.Errorf("Bad response status: %s", response.Status)
		}

		tracksResponse := tracksResponse{}
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return tracks, fmt.Errorf("Unable to read body: %w", err)
		}
		err = json.Unmarshal(body, &tracksResponse)
		if err != nil {
			return tracks, fmt.Errorf("Unable to unmarshal: %w", err)
		}
		tracks = append(tracks, tracksResponse.Items...)

		offset += limit
		if len(tracksResponse.Items) < limit {
			break
		}
	}

	return tracks, nil
}

func (api *SpotifyAPI) GetUser(tokens Credentials) (user User, err error) {
	request, err := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	if err != nil {
		return user, fmt.Errorf("Unable to create request: %w", err)
	}
	request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	request.Header.Add("Content-Type", "application/json")

	response, err := api.httpClient.Do(request)
	if err != nil {
		return user, fmt.Errorf("Unable to do request: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return user, fmt.Errorf("Bad response status: %s", response.Status)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return user, fmt.Errorf("Unable to read body: %w", err)
	}

	err = json.Unmarshal(body, &user)
	if err != nil {
		return user, fmt.Errorf("Unable to unmarshal: %w", err)
	}

	return user, nil
}

func (api *SpotifyAPI) CreatePlaylist(tokens Credentials, name string) (playlist Playlist, err error) {
	body := bytes.NewReader([]byte("{\"name\": \"" + name + "\"}"))

	user, err := api.GetUser(tokens)
	if err != nil {
		return playlist, fmt.Errorf("Unable to get user: %w", err)
	}

	url := fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists", user.ID)
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return playlist, fmt.Errorf("Unable to create request: %w", err)
	}
	request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	request.Header.Add("Content-Type", "application/json")

	response, err := api.httpClient.Do(request)
	if err != nil {
		return playlist, fmt.Errorf("Unable to do request: %w", err)
	}
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		return playlist, fmt.Errorf("Bad response status: %s", response.Status)
	}

	reqBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return playlist, err
	}

	err = json.Unmarshal(reqBody, &playlist)
	if err != nil {
		return playlist, err
	}

	return playlist, nil
}

func (api *SpotifyAPI) AddToPlaylist(tokens Credentials, playlistID string, tracks []Track) error {
	limit := 100

	for i := 0; i < len(tracks); i += limit {
		if i+limit > len(tracks) {
			limit = len(tracks) - i
		}

		uris := make([]string, limit)
		for j := i; j < i+limit; j++ {
			uris[j-i] = `"` + tracks[j].URI + `"`
		}
		data := "{\"uris\": [ " + strings.Join(uris, ",") + " ],\"position\":0}"
		dataReader := bytes.NewReader([]byte(data))

		request, err := http.NewRequest("POST", "https://api.spotify.com/v1/playlists/"+playlistID+"/tracks", dataReader)
		if err != nil {
			return fmt.Errorf("Unable to create request: %w", err)
		}
		request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
		request.Header.Add("Content-Type", "application/json")

		response, err := api.httpClient.Do(request)
		if err != nil {
			return fmt.Errorf("Unable to do request: %w", err)
		}

		if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
			return fmt.Errorf("Bad response status: %s", response.Status)
		}

		if i+limit >= len(tracks) {
			break
		}
	}
	return nil
}

func (api *SpotifyAPI) Authorized(tokens Credentials) (bool, error) {
	_, err := api.GetUser(tokens)
	return err == nil, nil
}
