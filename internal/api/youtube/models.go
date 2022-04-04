package youtube

type Credentials struct {
	AccessToken string `json:"access_token"`
}

type YoutubeConfig struct {
	APIKey       string
	ClientID     string
	ClientSecret string
	Scopes       string
	RedirectURI  string
}

type pageInfo struct {
	TotalResults int64 `json:"totalResults"`
}

type Playlist struct {
	ID      string  `json:"id"`
	Snippet snippet `json:"snippet"`
}

type playlistListResponse struct {
	PageInfo      pageInfo   `json:"pageInfo"`
	NextPageToken string     `json:"nextPageToken"`
	Items         []Playlist `json:"items"`
}

type snippet struct {
	Title        string `json:"title"`
	ChannelTitle string `json:"channelTitle"`
}

type Video struct {
	Snippet snippet `json:"snippet"`
}

type videoListResponse struct {
	Items         []Video  `json:"items"`
	PageInfo      pageInfo `json:"pageInfo"`
	NextPageToken string   `json:"nextPageToken"`
}
