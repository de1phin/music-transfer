package yandex

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type OnGetCredentials func(int64, Credentials)

func (api *YandexAPI) BindOnGetCredentials(onGetCredentials OnGetCredentials) {
	api.onGetCredentials = onGetCredentials
}

func (api *YandexAPI) checkStatus(userID int64, formTokens loginFormTokens, submitResponse submitResponse) {
	timeLimit := time.Now().Add(time.Second * 100)
	submitResponse.CsrfToken = url.QueryEscape(submitResponse.CsrfToken)
	url := "https://passport.yandex.ru/auth/new/magic/status/"
	for time.Now().Before(timeLimit) {
		time.Sleep(time.Second)
		data := "track_id=" + submitResponse.TrackID + "&csrf_token=" + submitResponse.CsrfToken
		req, err := http.NewRequest("POST", url, strings.NewReader(data))
		if err != nil {
			api.logger.Log("YandexAPI.checkStatus:", err)
			continue
		}
		for _, c := range formTokens.cookies {
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
			authResponse := authResponse{}
			json.Unmarshal(body, &authResponse)
			credentials := Credentials{
				UID: strconv.FormatInt(authResponse.UID, 10),
			}
			for _, c := range resp.Cookies() {
				credentials.Cookies += c.Name + "=" + c.Value + "; "
			}
			for _, c := range formTokens.cookies {
				credentials.Cookies += c.Name + "=" + c.Value + "; "
			}
			api.onGetCredentials(userID, credentials)
			break
		}
	}
}

func (api *YandexAPI) GetAuthURL(userID int64) (url string, err error) {
	formTokens, err := api.getYandexLoginFormTokens()
	if err != nil {
		api.logger.Log("YandexAPI.GetAuthURL.getYandexLoginFormTokens:", err)
		return url, err
	}

	submitResponse, err := api.getYandexSubmitResponse(formTokens)
	if err != nil {
		api.logger.Log("YandexAPI.GetAuthURL.getYandexSubmitResponse:", err)
		return url, err
	}

	if api.fixedAuthMagicToken == "" {
		timer := time.Now()
		svgQR, err := api.getQRCodeSVG(submitResponse.TrackID)
		if err != nil {
			api.logger.Log("YandexAPI.GetAuthURL.getQRCode:", err)
			return url, err
		}

		url, err = decodeQR(svgQR)
		if err != nil {
			api.logger.Log("YandexAPI.GetAuthURL.decodeQR:", err)
			return url, err
		}
		api.logger.Log("YandexAPI: QR fetched and decoded in", time.Since(timer))
	} else {
		url = "https://passport.yandex.ru/am/push/qrsecure?track_id=" + submitResponse.TrackID +
			"&magic=" + api.fixedAuthMagicToken
	}

	go api.checkStatus(userID, formTokens, submitResponse)

	return url, err
}

func (api *YandexAPI) getYandexLoginFormTokens() (tokens loginFormTokens, err error) {
	req, err := http.NewRequest("GET", "https://passport.yandex.ru/auth", nil)
	if err != nil {
		return tokens, err
	}
	resp, err := api.httpClient.Do(req)
	if err != nil {
		return tokens, err
	}
	if resp.StatusCode != http.StatusOK {
		return tokens, errors.New("YandexAPI.getYandexLoginFormTokens: Status - " + resp.Status)
	}
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return tokens, err
	}

	tokens.cookies = resp.Cookies()

	csrfTokenBegin := bytes.Index(html, []byte("csrf_token"))
	if csrfTokenBegin == -1 {
		return tokens, errors.New("YandexAPI.getYandexLoginFormTokens: No csrf token provided")
	}
	csrfTokenBegin = csrfTokenBegin + bytes.Index(html[csrfTokenBegin:], []byte("value=\"")) + 7
	csrfTokenEnd := csrfTokenBegin + bytes.Index(html[csrfTokenBegin:], []byte("\""))
	tokens.csrf = string(html[csrfTokenBegin:csrfTokenEnd])

	processUUIDBegin := bytes.Index(html, []byte("process_uuid=")) + 13
	if processUUIDBegin == -1 {
		return tokens, errors.New("YandexAPI.getYandexLoginFormTokens: No processUUID provided")
	}
	processUUIDEnd := processUUIDBegin + bytes.Index(html[processUUIDBegin:], []byte("\""))
	tokens.processUUID = string(html[processUUIDBegin:processUUIDEnd])

	return tokens, nil
}

func (api *YandexAPI) getYandexSubmitResponse(formTokens loginFormTokens) (submit submitResponse, err error) {
	data := "csrf_token=" + url.QueryEscape(formTokens.csrf) + "&process_uuid=" +
		url.QueryEscape(formTokens.processUUID) + "&with_code=1"
	req, err := http.NewRequest("POST", "https://passport.yandex.ru/registration-validations/auth/password/submit", strings.NewReader(data))
	if err != nil {
		return submit, err
	}
	for _, c := range formTokens.cookies {
		req.AddCookie(c)
	}

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return submit, err
	}
	if resp.StatusCode != http.StatusOK {
		return submit, errors.New("YandexAPI.getYandexSubmitResponse: Status:" + resp.Status)
	}
	if resp.Body == nil {
		return submit, errors.New("YandexAPI.getYandexSubmitResponse: Empty body returned")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return submit, err
	}

	err = json.Unmarshal(body, &submit)
	if err != nil {
		return submit, err
	}

	return submit, nil
}

func (api *YandexAPI) getQRCodeSVG(trackID string) (svg []byte, err error) {
	url := "https://passport.yandex.ru/auth/magic/code/?track_id=" + trackID
	data := "track_id=" + trackID
	req, err := http.NewRequest("GET", url, strings.NewReader(data))
	if err != nil {
		return svg, err
	}

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return svg, err
	}
	if resp.StatusCode != http.StatusOK {
		return svg, errors.New("YandexAPI.getQRCode: Status: " + resp.Status)
	}
	if resp.Body == nil {
		return svg, errors.New("YandexAPI.getQRCode: Empty body returned")
	}

	svg, err = ioutil.ReadAll(resp.Body)
	return svg, err
}
