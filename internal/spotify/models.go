package spotify

type tracksResponse struct {
	Items []trackItem `json:"items"`
}

type trackItem struct {
	Track track `json:"track"`
}

type track struct {
	Name    string   `json:"name"`
	Artists []artist `json:"artists"`
}

type artist struct {
	Name string `json:"name"`
}

type playlist struct {
	Name   string         `json:"name"`
	ID     string         `json:"ID"`
	Tracks tracksResponse `json:"tracks"`
}

type playlistResponse struct {
	Items []playlist `json:"items"`
}
