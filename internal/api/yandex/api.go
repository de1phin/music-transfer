package yandex

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

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
