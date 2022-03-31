package mockmusicservice

import (
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type credentials struct {
	token string
}

func (service *mockMusicService) GetAuthURL(id int64) string {
	return "mock has no url:(("
}

func (service *mockMusicService) Authorize(callback *http.Request) (int64, interface{}) {
	m, _ := url.ParseQuery(callback.URL.RawQuery)
	log.Println(m["id"])

	id, _ := strconv.ParseInt(m["id"][0], 10, 64)

	return id, credentials{token: "bibus"}
}

func (service *mockMusicService) ValidAuthCallback(callback *http.Request) bool {
	return false
}
