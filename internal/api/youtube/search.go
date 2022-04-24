package youtube

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

type searchResponse struct {
	Contents searchContent `json:"contents"`
}

type searchContent struct {
	TwoColumnSearchResultsRenderer twoColumnSearchResultsRenderer `json:"twoColumnSearchResultsRenderer"`
}

type twoColumnSearchResultsRenderer struct {
	PrimaryContents primaryContents `json:"primaryContents"`
}

type primaryContents struct {
	SectionListRenderer sectionListRenderer `json:"sectionListRenderer"`
}

type sectionListRenderer struct {
	Contents []sectionListRendererContent `json:"contents"`
}

type sectionListRendererContent struct {
	ItemSectionRenderer itemSectionRenderer `json:"itemSectionRenderer"`
}

type itemSectionRenderer struct {
	Contents []itemSectionRendererContent `json:"contents"`
}

type itemSectionRendererContent struct {
	VideoRenderer videoRenderer `json:"videoRenderer"`
}

type videoRenderer struct {
	VideoID string `json:"videoId"`
}

type client struct {
	ClientName    string `json:"clientName"`
	ClientVersion string `json:"clientVersion"`
}

type searchContext struct {
	Client client `json:"client"`
}

type searchRequest struct {
	Context searchContext `json:"context"`
	Query   string        `json:"query"`
}

func (api *YoutubeAPI) SearchVideo(title string, artists string) (string, error) {
	url := "https://www.youtube.com/youtubei/v1/search?key=AIzaSyAO_FJ2SlqU8Q4STEHLGCilw_Y9_11qcW8&prettyPrint=false"
	query := searchRequest{
		Context: searchContext{
			Client: client{
				ClientName:    "WEB",
				ClientVersion: "2.20220422.01.00",
			},
		},
		Query: strings.ReplaceAll(title+" "+artists, `"`, `\"`),
	}
	data, err := json.Marshal(query)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	search := searchResponse{}
	err = json.Unmarshal(body, &search)
	if err != nil {
		return "", err
	}

	if len(search.Contents.TwoColumnSearchResultsRenderer.PrimaryContents.SectionListRenderer.Contents) == 0 {
		return "", errors.New("YoutubeAPI.SearchVideo: Contents empty")
	}
	if len(search.Contents.TwoColumnSearchResultsRenderer.PrimaryContents.SectionListRenderer.Contents[0].ItemSectionRenderer.Contents) == 0 {
		return "", errors.New("YoutubeAPI.SearchVideo: Contents empty")
	}

	return search.Contents.TwoColumnSearchResultsRenderer.PrimaryContents.SectionListRenderer.Contents[0].ItemSectionRenderer.Contents[0].VideoRenderer.VideoID, nil

}
