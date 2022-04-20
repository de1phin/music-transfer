package yandex

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type OnGetCredentials func(int64, Credentials)

func (api *YandexAPI) BindOnGetCredentials(onGetCredentials OnGetCredentials) {
	api.onGetCredentials = onGetCredentials
}

func (api *YandexAPI) checkStatus(userID int64, yaFormTokens yandexLoginFormTokens, yaSubmit yandexSubmitResponse) {
	timeLimit := time.Now().Add(time.Second * 100)
	yaSubmit.CsrfToken = url.QueryEscape(yaSubmit.CsrfToken)
	url := "https://passport.yandex.ru/auth/new/magic/status/"
	for time.Now().Before(timeLimit) {
		time.Sleep(time.Second)
		data := "track_id=" + yaSubmit.TrackID + "&csrf_token=" + yaSubmit.CsrfToken
		req, err := http.NewRequest("POST", url, strings.NewReader(data))
		if err != nil {
			api.logger.Log("YandexAPI.checkStatus:", err)
			continue
		}
		for _, c := range yaFormTokens.cookies {
			req.AddCookie(c)
		}
		resp, err := api.httpClient.Do(req)
		if err != nil {
			api.logger.Log("YandexAPI.checkStatus:", err)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			api.logger.Log(errors.New("YandexAPI.checkStatus: Status: " + resp.Status))
			continue
		}

		hasSessionID := false
		for _, c := range resp.Cookies() {
			if c.Name == "Session_id" {
				hasSessionID = true
				break
			}
		}
		if hasSessionID {
			if resp.Body == nil {
				api.logger.Log(errors.New("YandexAPI.GetMe: Empty body returned"))
				continue
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				api.logger.Log(err)
				continue
			}
			authResponse := yandexAuthResponse{}
			json.Unmarshal(body, &authResponse)
			credentials := Credentials{
				UID:     authResponse.UID,
				cookies: resp.Cookies(),
			}
			credentials.cookies = append(credentials.cookies, yaFormTokens.cookies...)
			api.onGetCredentials(userID, credentials)
			break
		}
	}
}

func (api *YandexAPI) GetAuthURL(userID int64) (string, error) {
	yaFormTokens, err := api.getYandexLoginFormTokens()
	if err != nil {
		api.logger.Log("YandexAPI.GetAuthURL.getYandexLoginFormTokens:", err)
		return "", err
	}

	yaSubmit, err := api.getYandexSubmitResponse(yaFormTokens)
	if err != nil {
		api.logger.Log("YandexAPI.GetAuthURL.getYandexSubmitResponse:", err)
		return "", err
	}

	svgQR, err := api.getQRCodeSVG(yaSubmit.TrackID)
	if err != nil {
		api.logger.Log("YandexAPI.GetAuthURL.getQRCode:", err)
		return "", err
	}

	url, err := decodeQR(svgQR)

	go api.checkStatus(userID, yaFormTokens, yaSubmit)

	return url, err
}

func (api *YandexAPI) getYandexLoginFormTokens() (yandexLoginFormTokens, error) {
	result := yandexLoginFormTokens{}

	req, err := http.NewRequest("GET", "https://passport.yandex.ru/auth", nil)
	if err != nil {
		return result, err
	}
	resp, err := api.httpClient.Do(req)
	if err != nil {
		return result, err
	}
	if resp.StatusCode != http.StatusOK {
		return result, errors.New("YandexAPI.getYandexLoginFormTokens: Status - " + resp.Status)
	}
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	result.cookies = resp.Cookies()

	csrfTokenBegin := bytes.Index(html, []byte("csrf_token"))
	if csrfTokenBegin == -1 {
		return result, errors.New("YandexAPI.getYandexLoginFormTokens: No csrf token provided")
	}
	csrfTokenBegin = csrfTokenBegin + bytes.Index(html[csrfTokenBegin:], []byte("value=\"")) + 7
	csrfTokenEnd := csrfTokenBegin + bytes.Index(html[csrfTokenBegin:], []byte("\""))
	result.csrf = string(html[csrfTokenBegin:csrfTokenEnd])

	processUUIDBegin := bytes.Index(html, []byte("process_uuid=")) + 13
	if processUUIDBegin == -1 {
		return result, errors.New("YandexAPI.getYandexLoginFormTokens: No processUUID provided")
	}
	processUUIDEnd := processUUIDBegin + bytes.Index(html[processUUIDBegin:], []byte("\""))
	result.processUUID = string(html[processUUIDBegin:processUUIDEnd])

	return result, nil
}

func (api *YandexAPI) getYandexSubmitResponse(yaFormTokens yandexLoginFormTokens) (yandexSubmitResponse, error) {
	result := yandexSubmitResponse{}

	data := "csrf_token=" + url.QueryEscape(yaFormTokens.csrf) + "&process_uuid=" +
		url.QueryEscape(yaFormTokens.processUUID) + "&with_code=1"
	req, err := http.NewRequest("POST", "https://passport.yandex.ru/registration-validations/auth/password/submit", strings.NewReader(data))
	if err != nil {
		return result, err
	}
	for _, c := range yaFormTokens.cookies {
		req.AddCookie(c)
	}

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return result, err
	}
	if resp.StatusCode != http.StatusOK {
		return result, errors.New("YandexAPI.getYandexSubmitResponse: Status:" + resp.Status)
	}
	if resp.Body == nil {
		return result, errors.New("YandexAPI.getYandexSubmitResponse: Empty body returned")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}
	json.Unmarshal(body, &result)
	return result, nil
}

func (api *YandexAPI) getQRCodeSVG(trackID string) ([]byte, error) {
	url := "https://passport.yandex.ru/auth/magic/code/?track_id=" + trackID
	data := "track_id=" + trackID
	req, err := http.NewRequest("GET", url, strings.NewReader(data))
	if err != nil {
		return nil, err
	}

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("YandexAPI.getQRCode: Status: " + resp.Status)
	}
	if resp.Body == nil {
		return nil, errors.New("YandexAPI.getQRCode: Empty body returned")
	}
	return ioutil.ReadAll(resp.Body)
}
