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

type OnGetTokens func(int64, Credentials)

func (api *YoutubeAPI) BindHandler(router *http.ServeMux, onGetTokens OnGetTokens) {
	router.HandleFunc("/youtube", api.callbackHandler(onGetTokens))
}

func (api *YoutubeAPI) callbackHandler(onGetTokens OnGetTokens) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("oh, got something")
		m, _ := url.ParseQuery(r.URL.RawQuery)
		userID, _ := strconv.ParseInt(m["state"][0], 10, 64)

		code := m["code"][0]

		body := bytes.NewReader([]byte(fmt.Sprintf("code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=authorization_code", code, api.config.ClientID, api.config.ClientSecret, api.config.RedirectURI)))

		request, _ := http.NewRequest("POST", "https://oauth2.googleapis.com/token", body)
		request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		client := &http.Client{}
		response, _ := client.Do(request)

		jsonEncoded, _ := ioutil.ReadAll(response.Body)

		tokens := Credentials{}
		json.Unmarshal(jsonEncoded, &tokens)

		onGetTokens(userID, tokens)
	}
}
