package yandex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type OnGetCredentials func(int64, Credentials) error

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
			api.logger.Error(fmt.Errorf("Yandex: Unable to check status: Unable to create request: %w", err))
			continue
		}
		for _, c := range formTokens.cookies {
			req.AddCookie(c)
		}
		resp, err := api.httpClient.Do(req)
		if err != nil {
			api.logger.Error(fmt.Errorf("Yandex: Unable to check status: Unable to do request: %w", err))
			continue
		}
		if resp.StatusCode != http.StatusOK {
			api.logger.Error(fmt.Errorf("Yandex: Unable to check status: Bad response status: %s", resp.Status))
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
				api.logger.Error(fmt.Errorf("Yandex: Unable to check status: Empty body returned"))
				continue
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				api.logger.Error(fmt.Errorf("Yandex: Unable to check status: Unable to read body: %w", err))
				continue
			}
			authResponse := authResponse{}
			err = json.Unmarshal(body, &authResponse)
			if err != nil {
				api.logger.Error(fmt.Errorf("Yandex: Unable to check status: Unable to unmarshal: %w", err))
				continue
			}
			credentials := Credentials{
				UID: strconv.FormatInt(authResponse.UID, 10),
			}
			for _, c := range resp.Cookies() {
				credentials.Cookies += c.Name + "=" + c.Value + "; "
			}
			for _, c := range formTokens.cookies {
				credentials.Cookies += c.Name + "=" + c.Value + "; "
			}
			err = api.onGetCredentials(userID, credentials)
			if err != nil {
				api.logger.Error(fmt.Errorf("Yandex: Unable to check status: OnGetCredentials error: %w", err))
			}
			break
		}
	}
}

func (api *YandexAPI) GetAuthURL(userID int64) (url string, err error) {
	formTokens, err := api.getYandexLoginFormTokens()
	if err != nil {
		return url, fmt.Errorf("Unable to get yandex login form tokens: %w", err)
	}

	submitResponse, err := api.getYandexSubmitResponse(formTokens)
	if err != nil {
		return url, fmt.Errorf("Unable to get yandex submit response: %w", err)
	}

	if !api.UseFixedURL {
		timer := time.Now()
		svgQR, err := api.getQRCodeSVG(submitResponse.TrackID)
		if err != nil {
			return url, fmt.Errorf("Unable to get QR code: %w", err)
		}

		url, err = decodeQR(svgQR)
		if err != nil {
			return url, fmt.Errorf("Unable to decode QR: %w", err)
		}
		api.logger.Info("YandexAPI: QR fetched and decoded in", time.Since(timer))
	} else {
		url = "https://passport.yandex.ru/am/push/qrsecure?track_id=" + submitResponse.TrackID +
			"&magic=" + api.Magic
	}

	go api.checkStatus(userID, formTokens, submitResponse)

	return url, nil
}

func (api *YandexAPI) getYandexLoginFormTokens() (tokens loginFormTokens, err error) {
	req, err := http.NewRequest("GET", "https://passport.yandex.ru/auth", nil)
	if err != nil {
		return tokens, fmt.Errorf("Unable to create request: %w", err)
	}
	resp, err := api.httpClient.Do(req)
	if err != nil {
		return tokens, fmt.Errorf("Unable to do request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return tokens, fmt.Errorf("Bad response status: %s", resp.Status)
	}
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return tokens, fmt.Errorf("Unable to read body: %w", err)
	}

	tokens.cookies = resp.Cookies()

	csrfTokenBegin := bytes.Index(html, []byte("csrf_token"))
	if csrfTokenBegin == -1 {
		return tokens, fmt.Errorf("No csrf token provided")
	}
	csrfTokenBegin = csrfTokenBegin + bytes.Index(html[csrfTokenBegin:], []byte("value=\"")) + 7
	csrfTokenEnd := csrfTokenBegin + bytes.Index(html[csrfTokenBegin:], []byte("\""))
	tokens.csrf = string(html[csrfTokenBegin:csrfTokenEnd])

	processUUIDBegin := bytes.Index(html, []byte("process_uuid=")) + 13
	if processUUIDBegin == -1 {
		return tokens, fmt.Errorf("No processUUID provided")
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
		return submit, fmt.Errorf("Unable to create request: %w", err)
	}
	for _, c := range formTokens.cookies {
		req.AddCookie(c)
	}

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return submit, fmt.Errorf("Unable to do request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return submit, fmt.Errorf("Bad response status:" + resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return submit, fmt.Errorf("Unable to read body: %w", err)
	}

	err = json.Unmarshal(body, &submit)
	if err != nil {
		return submit, fmt.Errorf("Unable to unmarshal: %w", err)
	}

	return submit, nil
}

func (api *YandexAPI) getQRCodeSVG(trackID string) (svg []byte, err error) {
	url := "https://passport.yandex.ru/auth/magic/code/?track_id=" + trackID
	data := "track_id=" + trackID
	req, err := http.NewRequest("GET", url, strings.NewReader(data))
	if err != nil {
		return svg, fmt.Errorf("Unable to create request: %w", err)
	}

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return svg, fmt.Errorf("Unable to do request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return svg, fmt.Errorf("Bad response status: %w", err)
	}

	svg, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return svg, fmt.Errorf("Unable to read body: %w", err)
	}

	return svg, nil
}
