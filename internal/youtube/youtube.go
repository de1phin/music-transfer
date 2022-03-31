package youtube

type YouTubeService struct {
	apiKey       string
	clientID     string
	clientSecret string
	scope        string
	endpoint     string
	redirectURL  string
}

type YouTubeConfig interface {
	GetYouTubeApiKEY() string
	GetYouTubeClientID() string
	GetYouTubeClientSecret() string
	GetYouTubeScope() string
	GetYouTubeEndpoint() string
}

func NewYouTubeService(config YouTubeConfig) *YouTubeService {
	return &YouTubeService{
		apiKey:       config.GetYouTubeApiKEY(),
		clientID:     config.GetYouTubeClientID(),
		clientSecret: config.GetYouTubeClientSecret(),
		scope:        config.GetYouTubeScope(),
		endpoint:     config.GetYouTubeEndpoint(),
	}
}

func (*YouTubeService) Name() string {
	return "YouTube"
}

func (youtube *YouTubeService) InitCallbackServer(url string) (string, bool) {
	youtube.redirectURL = url

	return youtube.endpoint, true
}
