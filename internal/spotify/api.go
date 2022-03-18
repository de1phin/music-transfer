package spotify

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

func (spotify *spotifyService) Authorize(callback *http.Request) (int64, interface{}) {
	m, _ := url.ParseQuery(callback.URL.RawQuery)
	userID, _ := strconv.ParseInt(m["state"][0], 10, 64)

	authorizationCode := m["code"][0]

	body := bytes.NewReader([]byte(fmt.Sprintf("grant_type=authorization_code&code=%s&redirect_uri=%s/%s", authorizationCode, spotify.callbackURL, spotify.URLName())))
	request, _ := http.NewRequest("POST", "https://accounts.spotify.com/api/token", body)

	request.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(spotify.clientID+":"+spotify.clientSecret)))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	response, _ := client.Do(request)

	jsonEncoded, _ := io.ReadAll(response.Body)
	credentials := credentials{}
	json.Unmarshal(jsonEncoded, &credentials)

	return userID, credentials
}
