package spotify

import "fmt"

type spotifyService struct {
	clientID     string
	clientSecret string
	scope        string
	callbackURL  string
}

type credentials struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
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

func (spotify *spotifyService) GetAuthURL(id int64) string {
	return fmt.Sprintf("https://accounts.spotify.com/authorize?client_id=%s&response_type=code&redirect_uri=%s/%s&scope=%s&state=%d", spotify.clientID, spotify.callbackURL, spotify.URLName(), spotify.scope, id)
}
