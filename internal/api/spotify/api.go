package spotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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

func (api *SpotifyAPI) GetLiked(tokens Credentials) (playlist Playlist) {
	limit := 50
	offset := 0

	for {
		url := fmt.Sprintf("https://api.spotify.com/v1/me/tracks?limit=%d&offset=%d", limit, offset)
		request, _ := http.NewRequest("GET", url, nil)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+(tokens.AccessToken))

		response, _ := api.httpClient.Do(request)

		spotifyPlaylist := tracksResponse{}
		body, _ := io.ReadAll(response.Body)
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

	return playlist
}

func (api *SpotifyAPI) SearchTrack(tokens Credentials, title string, artists string) []Track {
	query := url.QueryEscape(title + " " + artists)
	url := "https://api.spotify.com/v1/search?type=track&q=" + query

	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	request.Header.Add("Content-Type", "application/json")

	response, _ := api.httpClient.Do(request)
	body, _ := ioutil.ReadAll(response.Body)
	searchResponse := searchResponse{}
	json.Unmarshal(body, &searchResponse)
	return searchResponse.Tracks.Items
}

func (api *SpotifyAPI) LikeTracks(tokens Credentials, tracks []Track) {
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

		request, _ := http.NewRequest("PUT", url, nil)
		request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
		request.Header.Add("Content-Type", "application/json")

		api.httpClient.Do(request)
		if i+limit >= len(tracks) {
			break
		}
	}
}

func (api *SpotifyAPI) GetUserPlaylists(tokens Credentials) (playlists []Playlist) {
	limit := 50
	offset := 0

	for {
		url := fmt.Sprintf("https://api.spotify.com/v1/me/playlists?limit=%d&offset=%d", limit, offset)
		request, _ := http.NewRequest("GET", url, nil)
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", "Bearer "+(tokens.AccessToken))

		response, _ := api.httpClient.Do(request)

		spotifyPlaylists := playlistResponse{}
		body, _ := io.ReadAll(response.Body)
		json.Unmarshal(body, &spotifyPlaylists)

		playlists = append(playlists, spotifyPlaylists.Items...)

		offset += limit
		if len(spotifyPlaylists.Items) < limit {
			break
		}
	}

	return playlists
}

func (api *SpotifyAPI) GetPlaylistTracks(tokens Credentials, playlistID string) (items []TrackItem) {
	limit := 50
	offset := 0

	for {
		url := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks?limit=%d&offset=%d&fields=items(track(name,artists(name)))", playlistID, limit, offset)
		request, _ := http.NewRequest("GET", url, nil)
		request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
		request.Header.Add("Content-Type", "application/json")

		response, _ := api.httpClient.Do(request)

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

func (api *SpotifyAPI) GetUser(tokens Credentials) User {
	request, _ := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	request.Header.Add("Content-Type", "application/json")

	response, _ := api.httpClient.Do(request)
	body, _ := ioutil.ReadAll(response.Body)
	user := User{}
	json.Unmarshal(body, &user)
	return user
}

func (api *SpotifyAPI) CreatePlaylist(tokens Credentials, name string) Playlist {
	body := bytes.NewReader([]byte("{\"name\": \"" + name + "\"}"))

	userID := api.GetUser(tokens).ID
	url := fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists", userID)
	request, _ := http.NewRequest("POST", url, body)
	request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	request.Header.Add("Content-Type", "application/json")

	response, _ := api.httpClient.Do(request)
	reqBody, _ := ioutil.ReadAll(response.Body)
	playlist := Playlist{}
	json.Unmarshal(reqBody, &playlist)
	return playlist
}

func (api *SpotifyAPI) AddToPlaylist(tokens Credentials, playlistID string, tracks []Track) {
	if len(tracks) == 0 {
		return
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

		request, _ := http.NewRequest("POST", "https://api.spotify.com/v1/playlists/"+playlistID+"/tracks", dataReader)
		request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
		request.Header.Add("Content-Type", "application/json")

		api.httpClient.Do(request)
	}
}

func (api *SpotifyAPI) Authorized(tokens Credentials) bool {
	request, _ := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	response, _ := api.httpClient.Do(request)
	return response.StatusCode != 400
}
