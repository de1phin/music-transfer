package spotify

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type credentials struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (spotify *spotifyService) GetAuthURL(id int64) string {
	return fmt.Sprintf("https://accounts.spotify.com/authorize?client_id=%s&response_type=code&redirect_uri=%s/%s&scope=%s&state=%d", spotify.clientID, spotify.callbackURL, spotify.URLName(), spotify.scope, id)
}

func (spotify *spotifyService) Authorize(callback *http.Request) (int64, interface{}) {
	m, _ := url.ParseQuery(callback.URL.RawQuery)
	userID, _ := strconv.ParseInt(m["state"][0], 10, 64)

	authorizationCode := m["code"][0]

	body := bytes.NewReader([]byte(fmt.Sprintf("grant_type=authorization_code&code=%s&redirect_uri=%s/spotify&state=authorized", authorizationCode, spotify.callbackURL)))
	request, _ := http.NewRequest("POST", "https://accounts.spotify.com/api/token", body)

	request.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(spotify.clientID+":"+spotify.clientSecret)))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	response, _ := client.Do(request)
	log.Println("Response is:", response.Request.URL)

	jsonEncoded, _ := io.ReadAll(response.Body)
	credentials := credentials{}
	json.Unmarshal(jsonEncoded, &credentials)

	return userID, credentials
}

func (spotify *spotifyService) ValidAuthCallback(callback *http.Request) bool {
	if callback == nil {
		return false
	}
	m, err := url.ParseQuery(callback.URL.RawQuery)
	if err != nil {
		return false
	}
	if len(m["state"]) != 1 {
		return false
	}
	if m["state"][0] == "Authorized" {
		return false
	}
	return true
}
