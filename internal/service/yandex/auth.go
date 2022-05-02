package yandex

import (
	"fmt"

	"github.com/de1phin/music-transfer/internal/api/yandex"
	"github.com/de1phin/music-transfer/internal/mux"
)

func (ya *Yandex) BindOnAuthorized(OnAuthorized mux.OnAuthorized) {
	ya.OnAuthorized = OnAuthorized
}

func (ya *Yandex) OnGetCredentials(userID int64, credentials yandex.Credentials) error {
	user, err := ya.api.GetMe(credentials)
	if err != nil {
		return fmt.Errorf("Unable to get me: %w", err)
	}
	credentials.Login = user.Login
	err = ya.storage.Set(userID, credentials)
	if err != nil {
		return fmt.Errorf("Unable to set credentials: %w", err)
	}

	ya.OnAuthorized(ya, userID)
	return nil
}

func (ya *Yandex) GetAuthURL(userID int64) (string, error) {
	return ya.api.GetAuthURL(userID)
}

func (ya *Yandex) Authorized(userID int64) (bool, error) {
	ok, err := ya.storage.Exist(userID)
	if err != nil {
		return ok, fmt.Errorf("Unable to check credentials: %w", err)
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
		return false, fmt.Errorf("Unable to get me: %w", err)
	}
	return true, nil
}
