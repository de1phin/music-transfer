package spotify

type SpotifyAPI struct {
	client      Client
	redirectURI string
}

func NewSpotifyAPI(client Client, hostname string) *SpotifyAPI {
	return &SpotifyAPI{client: client, redirectURI: hostname + "/spotify"}
}
