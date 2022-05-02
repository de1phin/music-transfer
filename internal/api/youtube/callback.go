package youtube

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type OnGetTokens func(int64, Credentials) error

func (api *YoutubeAPI) BindHandler(router *http.ServeMux, onGetTokens OnGetTokens) {
	router.HandleFunc("/youtube", api.callbackHandler(onGetTokens))
}

func (api *YoutubeAPI) callbackHandler(onGetTokens OnGetTokens) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		m, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			api.logger.Error(fmt.Errorf("Youtube: Unable to handle callback: Unable to parse URL query: %w", err))
			return
		}
		if len(m["state"]) == 0 {
			api.logger.Error(fmt.Errorf("Youtube: Unable to handle callback: No state provided"))
			return
		}
		userID, err := strconv.ParseInt(m["state"][0], 10, 64)
		if err != nil {
			api.logger.Error(fmt.Errorf("Youtube: Unable to handle callback: Unable to parse int: %w", err))
			return
		}

		if len(m["code"]) == 0 {
			api.logger.Error(fmt.Errorf("Youtube: Unable to handle callback: No code provided"))
			return
		}
		code := m["code"][0]

		body := bytes.NewReader([]byte(fmt.Sprintf("code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=authorization_code", code, api.ClientID, api.ClientSecret, api.RedirectURI)))

		request, err := http.NewRequest("POST", "https://oauth2.googleapis.com/token", body)
		if err != nil {
			api.logger.Error(fmt.Errorf("Youtube: Unable to handle callback: Unable to create request: %w", err))
			return
		}
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		response, err := api.httpClient.Do(request)
		if err != nil {
			api.logger.Error(fmt.Errorf("Youtube: Unable to handle callback: Unable to do request: %w", err))
			return
		}

		jsonEncoded, err := ioutil.ReadAll(response.Body)
		if err != nil {
			api.logger.Error(fmt.Errorf("Youtube: Unable to handle callback: Unable to read body: %w", err))
			return
		}

		tokens := Credentials{}
		err = json.Unmarshal(jsonEncoded, &tokens)
		if err != nil {
			api.logger.Error(fmt.Errorf("Youtube: Unable to handle callback: Unable to unmarshal: %w", err))
			return
		}

		err = onGetTokens(userID, tokens)
		if err != nil {
			api.logger.Error(fmt.Errorf("Youtube: Unable to handle callback: OnGetTokens: %w", err))
			return
		}
	}
}
