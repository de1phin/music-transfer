package spotify

type tracksResponse struct {
	Items []TrackItem `json:"items"`
}

type TrackItem struct {
	Track Track `json:"track"`
}

type Track struct {
	Name    string   `json:"name"`
	Artists []Artist `json:"artists"`
}

type Artist struct {
	Name string `json:"name"`
}

type Playlist struct {
	Name   string         `json:"name"`
	ID     string         `json:"ID"`
	Tracks tracksResponse `json:"tracks"`
}

type playlistResponse struct {
	Items []Playlist `json:"items"`
}

type Credentials struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Client struct {
	ID     string
	Secret string
}
