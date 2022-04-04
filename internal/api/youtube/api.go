package youtube

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type YoutubeAPI struct {
	httpClient *http.Client
	config     *YoutubeConfig
}

func NewYoutubeAPI(config *YoutubeConfig) *YoutubeAPI {
	return &YoutubeAPI{
		config:     config,
		httpClient: &http.Client{},
	}
}

func (api *YoutubeAPI) GetLiked(tokens Credentials) (videos []Video) {
	limit := 50

	url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?myRating=like&part=id,snippet&maxResults=%d", limit)

	for {
		request, _ := http.NewRequest("GET", url, nil)
		request.Header.Add("Authorization", "Bearer "+tokens.AccessToken)

		response, _ := api.httpClient.Do(request)
		body, _ := ioutil.ReadAll(response.Body)

		videoList := videoListResponse{}
		json.Unmarshal(body, &videoList)

		videos = append(videos, videoList.Items...)

		if videoList.NextPageToken == "" {
			break
		}

		url = fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?myRating=like&part=id,snippet&maxResults=%d&pageToken=%s", limit, videoList.NextPageToken)
	}

	return videos
}
