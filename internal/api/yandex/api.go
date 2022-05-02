package yandex

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func (api *YandexAPI) GetMe(credentials Credentials) (user User, err error) {
	req, err := http.NewRequest("GET", "https://api.passport.yandex.ru/all_accounts", nil)
	if err != nil {
		return user, fmt.Errorf("Unable to create request: %w", err)
	}
	req.Header.Add("X-Yandex-Music-Client", "YandexMusicAPI")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Cookie", credentials.Cookies)

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return user, fmt.Errorf("Unable to do request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return user, fmt.Errorf("Bad response status: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return user, fmt.Errorf("Unable to read body: %w", err)
	}
	acc := Accounts{}
	err = json.Unmarshal(body, &acc)
	if err != nil {
		return user, fmt.Errorf("Unable to unmarshal: %w", err)
	}

	for _, u := range acc.Users {
		if u.ID == credentials.UID {
			user = u
			break
		}
	}

	return user, nil
}

func (api *YandexAPI) GetLibrary(credentials Credentials) (library Library, err error) {
	user, err := api.GetMe(credentials)
	if err != nil {
		return library, fmt.Errorf("Unable to get me: %w", err)
	}
	url := "https://music.yandex.ru/handlers/library.jsx?owner=" + user.Login + "&filter=playlists"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return library, fmt.Errorf("Unable to create request: %w", err)
	}
	req.Header.Add("Cookie", credentials.Cookies)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Referer", "https://music.yandex.ru/users/"+user.Login+"/playlists")

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return library, fmt.Errorf("Unable to do request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return library, fmt.Errorf("Bad response status: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return library, fmt.Errorf("Unable to read body: %w", err)
	}

	err = json.Unmarshal(body, &library)
	if err != nil {
		return library, fmt.Errorf("Unable to do unmarshal: %w", err)
	}

	return library, nil
}

func (api *YandexAPI) GetPlaylist(ID int64, credentials Credentials) (playlist Playlist, err error) {
	user, err := api.GetMe(credentials)
	if err != nil {
		return playlist, fmt.Errorf("Unable to get me: %w", err)
	}
	url := "https://music.yandex.ru/handlers/playlist.jsx?owner=" + user.Login + "&kinds=" + strconv.FormatInt(ID, 10)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return playlist, fmt.Errorf("Unable to create request: %w", err)
	}
	req.Header.Add("Cookie", credentials.Cookies)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Referer", "https://music.yandex.ru/users/"+user.Login+"/playlists")

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return playlist, fmt.Errorf("Unable to do request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return playlist, fmt.Errorf("Bad response status: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return playlist, fmt.Errorf("Unable to read body: %w", err)
	}
	playlistResponse := PlaylistResponse{}
	err = json.Unmarshal(body, &playlistResponse)
	if err != nil {
		return playlist, fmt.Errorf("Unable to unmarshal: %w", err)
	}

	playlist = playlistResponse.Playlist
	return playlistResponse.Playlist, nil
}

func (api *YandexAPI) SearchTrack(title string, artists string) (track Track, err error) {
	url := "https://music.yandex.ru/handlers/music-search.jsx?text=" + url.QueryEscape(title+" "+artists) + "&type=all"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return track, fmt.Errorf("Unable to create request: %w", err)
	}
	req.Header.Add("Accept", "application/json")

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return track, fmt.Errorf("Unable to do request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return track, fmt.Errorf("Bad response status: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return track, fmt.Errorf("Unable to read body: %w", err)
	}
	search := SearchResponse{}
	err = json.Unmarshal(body, &search)
	if err != nil {
		return track, fmt.Errorf("Unable to unmarshal: %w", err)
	}

	if len(search.Tracks.Items) == 0 {
		return track, fmt.Errorf("No results")
	}

	track.ID = strconv.FormatInt(search.Tracks.Items[0].ID, 10)
	track.Title = search.Tracks.Items[0].Title
	track.Albums = search.Tracks.Items[0].Albums
	track.Artists = search.Tracks.Items[0].Artist
	track.Type = search.Tracks.Items[0].Type
	return track, nil
}

func (api *YandexAPI) GetAuthTokens(credentials Credentials) (tokens AuthTokens, err error) {
	requrl := "https://music.yandex.ru/handlers/auth.jsx"
	req, err := http.NewRequest("GET", requrl, nil)
	if err != nil {
		return tokens, fmt.Errorf("Unable to create request: %w", err)
	}
	req.Header.Add("Cookie", credentials.Cookies)
	req.Header.Add("Referer", url.QueryEscape("https://music.yandex.ru/users/"+credentials.Login+"/playlists"))
	req.Header.Add("X-Retpath-Y", url.QueryEscape("https://music.yandex.ru/users/"+credentials.Login+"/playlists"))

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return tokens, fmt.Errorf("Unable to do request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return tokens, fmt.Errorf("Bad response status: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return tokens, fmt.Errorf("Unable to read body: %w", err)
	}
	authTokens := authTokensResponse{}
	err = json.Unmarshal(body, &authTokens)
	if err != nil {
		return tokens, fmt.Errorf("Unable to unmarshal: %w", err)
	}

	tokens = authTokens.User
	return authTokens.User, nil
}

func (api *YandexAPI) LikeTrack(track Track, credentials Credentials, authTokens AuthTokens) error {
	data := "sign=" + url.QueryEscape(authTokens.Sign)
	url := "https://music.yandex.ru/api/v2.1/handlers/track/" + track.ID + "/web-own_playlists-playlist-track-main/like/add"
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return fmt.Errorf("Unable to create request: %w", err)
	}
	req.Header.Add("Cookie", credentials.Cookies)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Retpath-Y", "https%3A%2F%2Fmusic.yandex.ru%2Fusers%2F"+credentials.Login+"%2Fplaylists")

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("Unable to do request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Bad response status: %s", resp.Status)
	}

	return nil
}

func (api *YandexAPI) AddPlaylist(title string, credentials Credentials, authTokens AuthTokens) (playlist PlaylistSnippet, err error) {
	data := "action=add&title=" + url.QueryEscape(title) + "&sign=" + url.QueryEscape(authTokens.Sign) +
		"&external-domain=music.yandex.ru&overembed=false&lang=ru"
	url := "https://music.yandex.ru/handlers/change-playlist.jsx"
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return playlist, fmt.Errorf("Unable to create request: %w", err)
	}
	req.Header.Add("Cookie", credentials.Cookies)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Referer", "https://music.yandex.ru/users/"+credentials.Login+"/playlists")
	req.Header.Add("X-Retpath-Y", "https://music.yandex.ru/users/"+credentials.Login+"/playlists")

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return playlist, fmt.Errorf("Unable to do request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return playlist, fmt.Errorf("Bad response status: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return playlist, fmt.Errorf("Unable to read body: %w", err)
	}
	snippet := playlistAddResponse{}
	err = json.Unmarshal(body, &snippet)
	if err != nil {
		return playlist, fmt.Errorf("Unable to unmarshal: %w", err)
	}

	playlist = snippet.Playlist
	return playlist, nil
}

func (api *YandexAPI) AddToPlaylist(tracks []TrackSnippet, playlist PlaylistSnippet, credentials Credentials, authTokens AuthTokens) error {
	diff := playlistPatchDifference{
		At:        0,
		Operation: "insert",
		Tracks:    tracks,
	}
	diffstr, err := json.Marshal(diff)
	if err != nil {
		return fmt.Errorf("Unable to do marshal: %w", err)
	}
	data := fmt.Sprintf("revision=1&owner=%s&kind=%d&diff=[%s]&sign=%s", credentials.UID, playlist.Kind, url.QueryEscape(string(diffstr)), url.QueryEscape(authTokens.Sign))
	url := "https://music.yandex.ru/handlers/playlist-patch.jsx"
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return fmt.Errorf("Unable to create request: %w", err)
	}
	req.Header.Add("Cookie", credentials.Cookies)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Referer", "https://music.yandex.ru/users/"+credentials.Login+"/playlists")
	req.Header.Add("X-Current-UID", credentials.UID)
	req.Header.Add("X-Retpath-Y", "https://music.yandex.ru/users/"+credentials.Login+"/playlists")

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("Unable to do request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Bad response status: %s", resp.Status)
	}
	if resp.Body == nil {
		return fmt.Errorf("Empty body returned")
	}

	return nil
}
