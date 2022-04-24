package yandex

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func (api *YandexAPI) GetMe(credentials *Credentials) (*User, error) {
	req, err := http.NewRequest("GET", "https://api.passport.yandex.ru/all_accounts", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-Yandex-Music-Client", "YandexMusicAPI")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Cookie", credentials.Cookies)

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, errors.New("YandexAPI.GetMe: " + err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("YandexAPI.GetMe: Status: " + resp.Status)
	}
	if resp.Body == nil {
		return nil, errors.New("YandexAPI.GetMe: Empty Response Body")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("YandexAPI.GetMe: " + err.Error())
	}
	acc := Accounts{}
	json.Unmarshal(body, &acc)
	for _, u := range acc.Users {
		if u.ID == credentials.UID {
			return &u, nil
		}
	}
	return nil, errors.New("YandexAPI.GetMe: No valid user returned")
}

func (api *YandexAPI) GetLibrary(credentials *Credentials) (*Library, error) {
	user, err := api.GetMe(credentials)
	if err != nil {
		return nil, errors.New("YandexAPI.GetLibrary: " + err.Error())
	}
	url := "https://music.yandex.ru/handlers/library.jsx?owner=" + user.Login + "&filter=playlists"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.New("YandexAPI.GetLibrary: " + err.Error())
	}
	req.Header.Add("Cookie", credentials.Cookies)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Referer", "https://music.yandex.ru/users/"+user.Login+"/playlists")

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, errors.New("YandexAPI.GetLibrary: " + err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("YandexAPI.GetLibrary: Status: " + resp.Status)
	}
	if resp.Body == nil {
		return nil, errors.New("YandexAPI.GetLibrary: Empty Body returned")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("YandexAPI.GetLibrary: " + err.Error())
	}
	library := Library{}
	err = json.Unmarshal(body, &library)
	if err != nil {
		return nil, errors.New("YandexAPI.GetLibrary: " + err.Error())
	}

	return &library, nil
}

func (api *YandexAPI) GetPlaylist(ID int64, credentials *Credentials) (*Playlist, error) {
	user, err := api.GetMe(credentials)
	if err != nil {
		return nil, errors.New("YandexAPI.GetPlaylist: " + err.Error())
	}
	url := "https://music.yandex.ru/handlers/playlist.jsx?owner=" + user.Login + "&kinds=" + strconv.FormatInt(ID, 10)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.New("YandexAPI.GetPlaylist: " + err.Error())
	}
	req.Header.Add("Cookie", credentials.Cookies)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Referer", "https://music.yandex.ru/users/"+user.Login+"/playlists")

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, errors.New("YandexAPI.GetPlaylist: " + err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("YandexAPI.GetPlaylist: Status: " + resp.Status)
	}
	if resp.Body == nil {
		return nil, errors.New("YandexAPI.GetPlaylist: Empty Body returned")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("YandexAPI.GetPlaylist: " + err.Error())
	}
	playlistResponse := PlaylistResponse{}
	err = json.Unmarshal(body, &playlistResponse)
	if err != nil {
		return nil, errors.New("YandexAPI.GetPlaylist: " + err.Error())
	}

	return &playlistResponse.Playlist, nil
}

func (api *YandexAPI) SearchTrack(title string, artists string) (*Track, error) {
	url := "https://music.yandex.ru/handlers/music-search.jsx?text=" + url.QueryEscape(title+" "+artists) + "&type=all"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.New("YandexAPI.SearchTrack: " + err.Error())
	}
	req.Header.Add("Accept", "application/json")

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, errors.New("YandexAPI.SearchTrack: " + err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("YandexAPI.SearchTrack: Status: " + resp.Status)
	}
	if resp.Body == nil {
		return nil, errors.New("YandexAPI.SearchTrack: Empty body returned")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	search := SearchResponse{}
	err = json.Unmarshal(body, &search)
	if err != nil {
		return nil, errors.New("YandexAPI.SearchTrack: " + err.Error())
	}

	if len(search.Tracks.Items) == 0 {
		return nil, nil
	} else {
		track := &Track{
			ID:      strconv.FormatInt(search.Tracks.Items[0].ID, 10),
			Title:   search.Tracks.Items[0].Title,
			Artists: search.Tracks.Items[0].Artist,
			Albums:  search.Tracks.Items[0].Albums,
			Type:    search.Tracks.Items[0].Type,
		}
		return track, nil
	}
}

func (api *YandexAPI) GetAuthTokens(credentials *Credentials) (*AuthTokens, error) {
	requrl := "https://music.yandex.ru/handlers/auth.jsx"
	req, err := http.NewRequest("GET", requrl, nil)
	if err != nil {
		return nil, errors.New("YandexAPI.GetAuthTokens: " + err.Error())
	}
	req.Header.Add("Cookie", credentials.Cookies)
	req.Header.Add("Referer", url.QueryEscape("https://music.yandex.ru/users/"+credentials.Login+"/playlists"))
	req.Header.Add("X-Retpath-Y", url.QueryEscape("https://music.yandex.ru/users/"+credentials.Login+"/playlists"))

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, errors.New("YandexAPI.GetAuthTokens: " + err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("YandexAPI.GetAuthTokens: Status: " + resp.Status)
	}
	if resp.Body == nil {
		return nil, errors.New("YandexAPI.GetAuthTokens: Empty body returned")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("YandexAPI.GetAuthTokens: " + err.Error())
	}
	authTokens := authTokensResponse{}
	err = json.Unmarshal(body, &authTokens)
	if err != nil {
		return nil, errors.New("YandexAPI.GetAuthTokens: " + err.Error())
	}

	return &authTokens.User, nil
}

func (api *YandexAPI) LikeTrack(track *Track, credentials *Credentials, authTokens *AuthTokens) error {
	data := "sign=" + url.QueryEscape(authTokens.Sign)
	url := "https://music.yandex.ru/api/v2.1/handlers/track/" + track.ID + "/web-own_playlists-playlist-track-main/like/add"
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return errors.New("YandexAPI.LikeTrack: " + err.Error())
	}
	req.Header.Add("Cookie", credentials.Cookies)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Retpath-Y", "https%3A%2F%2Fmusic.yandex.ru%2Fusers%2F"+credentials.Login+"%2Fplaylists")

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return errors.New("YandexAPI.LikeTrack: " + err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("YandexAPI.LikeTrack: Status: " + resp.Status)
	}
	if resp.Body == nil {
		return errors.New("YandexAPI.LikeTrack: Empty body returned")
	}

	return nil
}

func (api *YandexAPI) AddPlaylist(title string, credentials *Credentials, authTokens *AuthTokens) (*PlaylistSnippet, error) {
	data := "action=add&title=" + url.QueryEscape(title) + "&sign=" + url.QueryEscape(authTokens.Sign) +
		"&external-domain=music.yandex.ru&overembed=false&lang=ru"
	url := "https://music.yandex.ru/handlers/change-playlist.jsx"
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return nil, errors.New("YandexAPI.AddPlaylist: " + err.Error())
	}
	req.Header.Add("Cookie", credentials.Cookies)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Referer", "https://music.yandex.ru/users/"+credentials.Login+"/playlists")
	req.Header.Add("X-Retpath-Y", "https://music.yandex.ru/users/"+credentials.Login+"/playlists")

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, errors.New("YandexAPI.AddPlaylist: " + err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("YandexAPI.AddPlaylist: Status: " + resp.Status)
	}
	if resp.Body == nil {
		return nil, errors.New("YandexAPI.AddPlaylist: Empty body returned")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("YandexAPI.AddPlaylist: " + err.Error())
	}
	snippet := playlistAddResponse{}
	err = json.Unmarshal(body, &snippet)
	if err != nil {
		return nil, errors.New("YandexAPI.AddPlaylist: " + err.Error())
	}

	return &snippet.Playlist, nil
}

func (api *YandexAPI) AddToPlaylist(tracks []TrackSnippet, playlist *PlaylistSnippet, credentials *Credentials, authTokens *AuthTokens) error {
	diff := playlistPatchDifference{
		At:        0,
		Operation: "insert",
		Tracks:    tracks,
	}
	diffstr, err := json.Marshal(diff)
	if err != nil {
		return errors.New("YandexAPI.AddToPlaylist: " + err.Error())
	}
	data := fmt.Sprintf("revision=1&owner=%s&kind=%d&diff=[%s]&sign=%s", credentials.UID, playlist.Kind, url.QueryEscape(string(diffstr)), url.QueryEscape(authTokens.Sign))
	url := "https://music.yandex.ru/handlers/playlist-patch.jsx"
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return errors.New("YandexAPI.AddToPlaylist: " + err.Error())
	}
	req.Header.Add("Cookie", credentials.Cookies)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Referer", "https://music.yandex.ru/users/"+credentials.Login+"/playlists")
	req.Header.Add("X-Current-UID", credentials.UID)
	req.Header.Add("X-Retpath-Y", "https://music.yandex.ru/users/"+credentials.Login+"/playlists")

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return errors.New("YandexAPI.AddToPlaylist: " + err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("YandexAPI.AddToPlaylist: Status: " + resp.Status)
	}
	if resp.Body == nil {
		return errors.New("YandexAPI.AddToPlaylist: Empty body returned")
	}

	return nil
}
