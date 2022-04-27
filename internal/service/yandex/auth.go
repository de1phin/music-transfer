package yandex

import (
	"errors"

	"github.com/de1phin/music-transfer/internal/api/yandex"
	"github.com/de1phin/music-transfer/internal/mux"
)

func (ya *Yandex) BindOnAuthorized(OnAuthorized mux.OnAuthorized) {
	ya.OnAuthorized = OnAuthorized
}

func (ya *Yandex) OnGetCredentials(userID int64, credentials yandex.Credentials) {
	user, err := ya.api.GetMe(credentials)
	if err != nil {
		ya.logger.Log(err)
		return
	}
	credentials.Login = user.Login
	err = ya.storage.Put(userID, credentials)
	if err != nil {
		ya.logger.Log(errors.New("Yandex.OnGetCredentials: " + err.Error()))
	}

	ya.OnAuthorized(ya, userID)
}

func (ya *Yandex) GetAuthURL(userID int64) (string, error) {
	return ya.api.GetAuthURL(userID)
}

func (ya *Yandex) Authorized(userID int64) (bool, error) {
	ok, err := ya.storage.Exist(userID)
	if err != nil {
		return ok, err
	}
	if !ok {
		return false, nil
	}
	credentials, err := ya.storage.Get(userID)
	if err != nil {
		return false, nil
	}
	_, err = ya.api.GetMe(credentials)
	if err != nil {
		return false, err
	}
	return true, nil
}
