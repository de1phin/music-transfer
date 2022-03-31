package spotify

type spotifyService struct {
	clientID     string
	clientSecret string
	scope        string
	callbackURL  string
	endpoint     string
}

type spotifyConfig interface {
	GetSpotifyClientID() string
	GetSpotifyClientSecret() string
	GetSpotifyScope() string
	GetSpotifyEndpoint() string
}

func NewSpotifyService(config spotifyConfig) *spotifyService {
	return &spotifyService{
		clientID:     config.GetSpotifyClientID(),
		clientSecret: config.GetSpotifyClientSecret(),
		scope:        config.GetSpotifyScope(),
		endpoint:     config.GetSpotifyEndpoint(),
	}
}

func (spotify *spotifyService) InitCallbackServer(url string) (endpoint string, doSetup bool) {
	spotify.callbackURL = url
	return spotify.endpoint, true
}

func (spotify *spotifyService) Name() string {
	return "Spotify"
}
