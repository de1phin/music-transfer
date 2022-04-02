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

type OnGetTokens func(int64, Credentials)

func (api *SpotifyAPI) BindHandler(router *http.ServeMux, onGetTokens OnGetTokens) {
	router.HandleFunc("/spotify", api.callbackHandler(onGetTokens))
}

func (api *SpotifyAPI) callbackHandler(onGetTokens OnGetTokens) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		m, _ := url.ParseQuery(r.URL.RawQuery)
		userID, _ := strconv.ParseInt(m["state"][0], 10, 64)
		authorizationCode := m["code"][0]

		rawBody := fmt.Sprintf("grant_type=authorization_code&code=%s&redirect_uri=%s&state=done", authorizationCode, api.redirectURI)
		body := bytes.NewReader([]byte(rawBody))
		request, _ := http.NewRequest("POST", "https://accounts.spotify.com/api/token", body)

		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		basicToken := base64.StdEncoding.EncodeToString([]byte(api.client.ID + ":" + api.client.Secret))
		request.Header.Add("Authorization", "Basic "+basicToken)

		client := &http.Client{}
		response, _ := client.Do(request)

		respBuf, _ := io.ReadAll(response.Body)
		credentials := Credentials{}
		json.Unmarshal(respBuf, &credentials)

		onGetTokens(userID, credentials)
	}
}
