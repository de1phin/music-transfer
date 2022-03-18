package transfer

import "net/http"

type MusicService interface {
	Name() string
	URLName() string
	GetAuthURL(int64) string
	Authorize(callback *http.Request) (int64, interface{})
}
