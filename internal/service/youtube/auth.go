package youtube

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

func (youtube *YouTubeService) GetAuthURL(id int64) string {
	return fmt.Sprintf("https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&response_type=code&scope=%s&redirect_uri=%s%s&state=%d", youtube.clientID, youtube.scope, youtube.redirectURL, youtube.endpoint, id)
}

func (youtube *YouTubeService) Authorize(callback *http.Request) (int64, interface{}) {
	m, _ := url.ParseQuery(callback.URL.RawQuery)
	userID, _ := strconv.ParseInt(m["state"][0], 10, 64)

	code := m["code"][0]

	body := bytes.NewReader([]byte(fmt.Sprintf("code=%s&client_id=%s&client_secret=%s&redirect_uri=%s/%s&grant_type=authorization_code", code, youtube.clientID, youtube.clientSecret, youtube.redirectURL, youtube.endpoint)))

	request, _ := http.NewRequest("POST", "https://oauth2.googleapis.com/token", body)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	response, _ := client.Do(request)

	jsonEncoded, _ := ioutil.ReadAll(response.Body)

	log.Println(string(jsonEncoded))
	credentials := credentials{}
	json.Unmarshal(jsonEncoded, &credentials)

	return userID, credentials
}

func (youtube *YouTubeService) ValidAuthCallback(callback *http.Request) bool {
	if callback.URL.Path != youtube.endpoint {
		return false
	}
	return true
}
