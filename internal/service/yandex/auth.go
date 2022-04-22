package yandex

import (
	"errors"

	"github.com/de1phin/music-transfer/internal/api/yandex"
)

func (y *Yandex) OnGetCredentials(userID int64, credentials yandex.Credentials) {
	err := y.userStorage.Put(userID, credentials)
	if err != nil {
		y.logger.Log(errors.New("Yandex.OnGetCredentials: " + err.Error()))
	}
}
