package yandex

import (
	"encoding/json"
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
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("YandexAPI.GetMe: Status: " + resp.Status)
	}
	if resp.Body == nil {
		return nil, errors.New("YandexAPI.GetMe: Empty Response Body")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	url := "https://music.yandex.ru/handlers/library.jsx?owner=" + user.Login + "&filter=playlists&playlistsWithoutContent=true"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Cookie", credentials.Cookies)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Referer", "https://music.yandex.ru/users/"+user.Login+"/playlists")

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("YandexAPI.GetLibrary: Status: " + resp.Status)
	}
	if resp.Body == nil {
		return nil, errors.New("YandexAPI.GetLibrary: Empty Body returned")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	library := Library{}
	err = json.Unmarshal(body, &library)

	return &library, err
}

func (api *YandexAPI) GetPlaylist(ID int64, credentials *Credentials) (*Playlist, error) {
	user, err := api.GetMe(credentials)
	if err != nil {
		return nil, err
	}
	url := "https://music.yandex.ru/handlers/playlist.jsx?owner=" + user.Login + "&kinds=" + strconv.FormatInt(ID, 10)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Cookie", credentials.Cookies)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Referer", "https://music.yandex.ru/users/"+user.Login+"/playlists")

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("YandexAPI.GetPlaylists: Status: " + resp.Status)
	}
	if resp.Body == nil {
		return nil, errors.New("YandexAPI.GetPlaylists: Empty Body returned")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	playlistResponse := PlaylistResponse{}
	err = json.Unmarshal(body, &playlistResponse)
	return &playlistResponse.Playlist, nil
}

func (api *YandexAPI) SearchTrack(title string, artists string) (*Track, error) {
	url := "https://music.yandex.ru/handlers/music-search.jsx?text=" + url.QueryEscape(title+" "+artists) + "&type=all"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	api.logger.Log(search)

	if len(search.Tracks.Items) == 0 {
		return nil, nil
	} else {
		track := &Track{
			ID:      strconv.FormatInt(search.Tracks.Items[0].ID, 10),
			Title:   search.Tracks.Items[0].Title,
			Artists: search.Tracks.Items[0].Artist,
			Type:    search.Tracks.Items[0].Type,
		}
		return track, nil
	}
}

func (api *YandexAPI) GetAuthCSRF(credentials *Credentials) (string, error) {
	url := "https://music.yandex.ru/api/v2.1/handlers/auth"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Cookie", credentials.Cookies)
	req.Header.Add("X-Retpath-Y", "https%3A%2F%2Fmusic.yandex.ru%2Fusers%2F"+credentials.Login+"%2Fplaylists%2F1015")

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("YandexAPI.GetAuthCSRF: Status: " + resp.Status)
	}
	if resp.Body == nil {
		return "", errors.New("YandexAPI.GetAuthCSRF: Empty body returned")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	authHandlerResponse := authHandlerResponse{}
	err = json.Unmarshal(body, &authHandlerResponse)

	return authHandlerResponse.Csrf, err
}

func (api *YandexAPI) LikeTrack(track *Track, credentials *Credentials, csrf string) error {
	data := "sign=" + url.QueryEscape(csrf)
	url := "https://music.yandex.ru/api/v2.1/handlers/track/" + track.ID + "/web-own_playlists-playlist-track-main/like/add"
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Add("Cookie", credentials.Cookies)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Retpath-Y", "https%3A%2F%2Fmusic.yandex.ru%2Fusers%2F"+credentials.Login+"%2Fplaylists")

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("YandexAPI.LikeTrack: Status: " + resp.Status)
	}
	if resp.Body == nil {
		return errors.New("YandexAPI.LikeTrack: Empty body returned")
	}

	return nil
}
