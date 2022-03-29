package youtube

import "log"

type YouTubeService struct {
	redirectURL  string
	scope        string
	apiKey       string
	clientID     string
	clientSecret string
}

func NewYouTubeService(redirectURL, scope, apiKey, clientID, clientSecret string) *YouTubeService {
	log.Println("Scope is", scope)
	return &YouTubeService{redirectURL, scope, apiKey, clientID, clientSecret}
}

func (*YouTubeService) Name() string {
	return "YouTube"
}

func (*YouTubeService) URLName() string {
	return "youtube"
}
