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

type OnGetTokens func(int64, Credentials) error

func (api *SpotifyAPI) BindHandler(router *http.ServeMux, onGetTokens OnGetTokens) {
	router.HandleFunc("/spotify", api.callbackHandler(onGetTokens))
}

func (api *SpotifyAPI) callbackHandler(onGetTokens OnGetTokens) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		m, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			api.logger.Log(fmt.Errorf("SpotifyAPI.callbackHandler: url.ParseQuery error: %w", err))
			return
		}
		if len(m["state"]) == 0 {
			api.logger.Log(fmt.Errorf("SpotifyAPI.callbackHandler: No state provided"))
			return
		}
		userID, err := strconv.ParseInt(m["state"][0], 10, 64)
		if err != nil {
			api.logger.Log(fmt.Errorf("SpotifyAPI.callbackHandler: strconv.ParseInt error: %w", err))
			return
		}
		if len(m["state"]) == 0 {
			api.logger.Log(fmt.Errorf("SpotifyAPI.callbackHandler: No code provided"))
			return
		}
		authorizationCode := m["code"][0]

		rawBody := fmt.Sprintf("grant_type=authorization_code&code=%s&redirect_uri=%s&state=done", authorizationCode, api.redirectURI)
		body := bytes.NewReader([]byte(rawBody))
		request, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", body)
		if err != nil {
			api.logger.Log(fmt.Errorf("SpotifyAPI.callbackHandler: http.NewRequest error: %w", err))
			return
		}

		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		basicToken := base64.StdEncoding.EncodeToString([]byte(api.client.ID + ":" + api.client.Secret))
		request.Header.Add("Authorization", "Basic "+basicToken)

		response, err := api.httpClient.Do(request)
		if err != nil {
			api.logger.Log(fmt.Errorf("SpotifyAPI.callbackHandler: httpClient.Do error: %w", err))
			return
		}

		respBuf, err := io.ReadAll(response.Body)
		if err != nil {
			api.logger.Log(fmt.Errorf("SpotifyAPI.callbackHandler: Read response body error: %w", err))
			return
		}

		credentials := Credentials{}
		err = json.Unmarshal(respBuf, &credentials)
		if err != nil {
			api.logger.Log(fmt.Errorf("SpotifyAPI.callbackHandler: Unmarshal response body error: %w", err))
			return
		}

		if err := onGetTokens(userID, credentials); err != nil {
			api.logger.Log(fmt.Errorf("SpotifyAPI.callbackHandler: OnGetTokens error: %w", err))
		}
	}
}
