package yandex

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

func (api *YandexAPI) GetMe(credentials Credentials) (User, error) {
	req, err := http.NewRequest("GET", "https://api.passport.yandex.ru/all_accounts", nil)
	if err != nil {
		return User{}, err
	}
	req.Header.Add("X-Yandex-Music-Client", "YandexMusicAPI")
	req.Header.Add("Accept", "application/json")

	for _, c := range credentials.cookies {
		req.AddCookie(c)
	}

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return User{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return User{}, errors.New("YandexAPI.GetMe: Status: " + resp.Status)
	}
	if resp.Body == nil {
		return User{}, errors.New("YandexAPI.GetMe: Empty Response Body")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return User{}, err
	}
	acc := Accounts{}
	json.Unmarshal(body, &acc)
	for _, u := range acc.Users {
		if u.ID == credentials.UID {
			return u, nil
		}
	}
	return User{}, errors.New("YandexAPI.GetMe: No valid user returned")
}
