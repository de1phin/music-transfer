package spotify

type tracksResponse struct {
	Items []trackItem `json:"items"`
}

type trackItem struct {
	Track track `json:"track"`
}

type track struct {
	Name string `json:"name"`
}
