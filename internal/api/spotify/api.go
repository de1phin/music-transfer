package spotify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type SpotifyAPI struct {
	httpClient  *http.Client
	client      Client
	redirectURI string
}

func NewSpotifyAPI(client Client, hostname string) *SpotifyAPI {
	return &SpotifyAPI{
		httpClient:  &http.Client{},
		client:      client,
		redirectURI: hostname + "/spotify",
	}
}

func (api *SpotifyAPI) GetLiked(tokens Credentials) (Playlist, error) {
	limit := 50
	offset := 0

	playlist := Playlist{}

	for {
		url := fmt.Sprintf("https://api.spotify.com/v1/me/tracks?limit=%d&offset=%d", limit, offset)
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return playlist, err
		}
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+(tokens.AccessToken))

		log.Println("Request:\n", request)
		response, err := api.httpClient.Do(request)
		if err != nil {
			return playlist, err
		}
		bd, _ := io.ReadAll(response.Body)
		log.Println("Response:", string(bd))
		if response.StatusCode != 200 {
			return playlist, errors.New("SpotifyAPI.GetLiked: " + response.Status)
		}

		spotifyPlaylist := tracksResponse{}
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return playlist, err
		}
		json.Unmarshal(body, &spotifyPlaylist)

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

func (api *SpotifyAPI) SearchTrack(tokens Credentials, title string, artists string) ([]Track, error) {
	query := url.QueryEscape(title + " " + artists)
	url := "https://api.spotify.com/v1/search?type=track&q=" + query

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	request.Header.Add("Content-Type", "application/json")

	response, err := api.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New("SpotifyAPI.SearchTrack: " + response.Status)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	searchResponse := searchResponse{}
	json.Unmarshal(body, &searchResponse)
	return searchResponse.Tracks.Items, nil
}

func (api *SpotifyAPI) LikeTracks(tokens Credentials, tracks []Track) error {
	limit := 50
	for i := 0; i < len(tracks); i += limit {
		if len(tracks)-i < limit {
			limit = len(tracks) - i
		}
		ids := ""
		for j := i; j < i+limit; j++ {
			ids += tracks[j].ID + ","
		}
		ids = ids[:len(ids)-1]
		url := fmt.Sprintf("https://api.spotify.com/v1/me/tracks?ids=%s", ids)

		request, err := http.NewRequest("PUT", url, nil)
		if err != nil {
			return err
		}
		request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
		request.Header.Add("Content-Type", "application/json")

		response, err := api.httpClient.Do(request)
		if err != nil {
			return err
		}
		if response.StatusCode != 200 {
			return errors.New("SpotifyAPI.LikeTracks: " + response.Status)
		}
		if i+limit >= len(tracks) {
			break
		}
	}
	return nil
}

func (api *SpotifyAPI) GetUserPlaylists(tokens Credentials) ([]Playlist, error) {
	limit := 50
	offset := 0

	playlists := []Playlist{}

	for {
		url := fmt.Sprintf("https://api.spotify.com/v1/me/playlists?limit=%d&offset=%d", limit, offset)
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return playlists, err
		}
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+(tokens.AccessToken))

		response, err := api.httpClient.Do(request)
		if err != nil {
			return playlists, err
		}
		if response.StatusCode != 200 {
			return playlists, errors.New("SpotifyAPI.GetUserPlaylists: " + response.Status)
		}

		spotifyPlaylists := playlistResponse{}
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return playlists, err
		}
		json.Unmarshal(body, &spotifyPlaylists)

		playlists = append(playlists, spotifyPlaylists.Items...)

		offset += limit
		if len(spotifyPlaylists.Items) < limit {
			break
		}
	}

	return playlists, nil
}

func (api *SpotifyAPI) GetPlaylistTracks(tokens Credentials, playlistID string) ([]TrackItem, error) {
	limit := 50
	offset := 0

	items := []TrackItem{}

	for {
		url := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks?limit=%d&offset=%d&fields=items(track(name,artists(name)))", playlistID, limit, offset)
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return items, err
		}
		request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
		request.Header.Add("Content-Type", "application/json")

		response, err := api.httpClient.Do(request)
		if err != nil {
			return items, err
		}
		if response.StatusCode != 200 {
			return items, errors.New("SpotifyAPI.GetPlaylistTracks: " + response.Status)
		}

		tracks := tracksResponse{}
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return items, err
		}
		json.Unmarshal(body, &tracks)
		items = append(items, tracks.Items...)

		offset += limit
		if len(tracks.Items) < limit {
			break
		}
	}

	return items, nil
}

func (api *SpotifyAPI) GetUser(tokens Credentials) (User, error) {
	request, err := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	if err != nil {
		return User{}, err
	}
	request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	request.Header.Add("Content-Type", "application/json")

	response, err := api.httpClient.Do(request)
	if err != nil {
		return User{}, err
	}
	if response.StatusCode != 200 {
		return User{}, errors.New("SpotifyAPI.GetUser: " + response.Status)
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return User{}, err
	}
	user := User{}
	json.Unmarshal(body, &user)
	return user, nil
}

func (api *SpotifyAPI) CreatePlaylist(tokens Credentials, name string) (Playlist, error) {
	body := bytes.NewReader([]byte("{\"name\": \"" + name + "\"}"))

	user, err := api.GetUser(tokens)
	if err != nil {
		return Playlist{}, err
	}

	url := fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists", user.ID)
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return Playlist{}, err
	}
	request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	request.Header.Add("Content-Type", "application/json")

	response, err := api.httpClient.Do(request)
	if err != nil {
		return Playlist{}, err
	}
	if response.StatusCode != 200 {
		return Playlist{}, errors.New("SpotifyAPI.CreatePlaylist: " + response.Status)
	}

	reqBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return Playlist{}, err
	}

	playlist := Playlist{}
	json.Unmarshal(reqBody, &playlist)
	return playlist, nil
}

func (api *SpotifyAPI) AddToPlaylist(tokens Credentials, playlistID string, tracks []Track) error {
	if len(tracks) == 0 {
		return errors.New("No tracks to add")
	}
	limit := 100

	for i := 0; i < len(tracks); i += limit {
		uris := ""
		for j := i; j < i+limit && j < len(tracks); j++ {
			uris += tracks[j].URI + ",\n"
		}
		uris = uris[:len(uris)-2]
		data := "{\"uris\": [ \"" + uris + "\" ],\"position\":0}"
		dataReader := bytes.NewReader([]byte(data))

		request, err := http.NewRequest("POST", "https://api.spotify.com/v1/playlists/"+playlistID+"/tracks", dataReader)
		if err != nil {
			return err
		}
		request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
		request.Header.Add("Content-Type", "application/json")

		response, err := api.httpClient.Do(request)
		if err != nil {
			return err
		}
		if response.StatusCode != 200 {
			return errors.New("SpotifyAPI.AddToPlaylist: " + response.Status)
		}
	}
	return nil
}

func (api *SpotifyAPI) Authorized(tokens Credentials) (bool, error) {
	_, err := api.GetUser(tokens)
	return err == nil, nil
}
