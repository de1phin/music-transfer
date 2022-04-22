package yandex

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

func (api *YandexAPI) GetMe(credentials Credentials) (User, error) {
	result := User{}

	req, err := http.NewRequest("GET", "https://api.passport.yandex.ru/all_accounts", nil)
	if err != nil {
		return result, err
	}
	req.Header.Add("X-Yandex-Music-Client", "YandexMusicAPI")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Cookie", credentials.Cookies)

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return result, err
	}
	if resp.StatusCode != http.StatusOK {
		return result, errors.New("YandexAPI.GetMe: Status: " + resp.Status)
	}
	if resp.Body == nil {
		return result, errors.New("YandexAPI.GetMe: Empty Response Body")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}
	acc := Accounts{}
	json.Unmarshal(body, &acc)
	for _, u := range acc.Users {
		if u.ID == credentials.UID {
			return u, nil
		}
	}
	return result, errors.New("YandexAPI.GetMe: No valid user returned")
}

func (api *YandexAPI) GetLibrary(credentials Credentials) (*Library, error) {
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

func (api *YandexAPI) GetPlaylist(ID int64, credentials Credentials) (*Playlist, error) {
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
