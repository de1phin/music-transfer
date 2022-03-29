package youtube

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/de1phin/music-transfer/internal/transfer"
)

func (youtube *YouTubeService) GetFavourites(tokens interface{}) transfer.Playlist {
	request, _ := http.NewRequest("GET", "https://www.googleapis.com/youtube/v3/videos?myRating=like&part=snippet,contentDetails", nil)
	request.Header.Add("Authorization", "Bearer "+tokens.(credentials).AccessToken)

	client := &http.Client{}

	response, _ := client.Do(request)
	if response.Body == nil {
		log.Println("PIZDEC")
		return transfer.Playlist{}
	}

	body, _ := ioutil.ReadAll(response.Body)
	log.Println("Response:", string(body))

	return transfer.Playlist{}
}

func (youtube *YouTubeService) AddFavourites(tokens interface{}, playlist transfer.Playlist) {
	log.Println("YouTube asked to add:\n", playlist)
}
