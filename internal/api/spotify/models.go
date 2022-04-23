package spotify

type tracksResponse struct {
	Items []TrackItem `json:"items"`
}

type TrackItem struct {
	Track Track `json:"track"`
}

type Track struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	URI     string   `json:"uri"`
	Artists []Artist `json:"artists"`
}

type Artist struct {
	Name string `json:"name"`
}

type Playlist struct {
	Name   string         `json:"name"`
	ID     string         `json:"id"`
	URI    string         `json:"uri"`
	Tracks tracksResponse `json:"tracks"`
}

type User struct {
	ID string `json:"id"`
}

type playlistResponse struct {
	Items []Playlist `json:"items"`
}

type searchTracks struct {
	Items []Track `json:"items"`
}

type searchResponse struct {
	Tracks searchTracks `json:"tracks"`
}

type Credentials struct {
	AccessToken  string `json:"access_token" db:"access_token"`
	RefreshToken string `json:"refresh_token" db:"refresh_token"`
}

type Client struct {
	ID     string
	Secret string
}
