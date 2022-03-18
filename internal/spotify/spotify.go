package spotify

type spotifyService struct {
	clientID     string
	clientSecret string
	scope        string
	callbackURL  string
}

func NewSpotifyService(clientID string, clientSecret string, scope string, callbackURL string) *spotifyService {
	return &spotifyService{clientID, clientSecret, scope, callbackURL}
}

func (spotify *spotifyService) Name() string {
	return "Spotify"
}

func (spotify *spotifyService) URLName() string {
	return "spotify"
}
