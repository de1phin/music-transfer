package youtube

type credentials struct {
	AccessToken string `json:"access_token"`
}

type pageInfo struct {
	TotalResults int64 `json:"totalResults"`
}

type playlist struct {
	ID      string  `json:"id"`
	Snippet snippet `json:"snippet"`
}

type playlistListResponse struct {
	PageInfo      pageInfo   `json:"pageInfo"`
	NextPageToken string     `json:"nextPageToken"`
	Items         []playlist `json:"items"`
}

type snippet struct {
	Title        string `json:"title"`
	ChannelTitle string `json:"channelTitle"`
}

type video struct {
	Snippet snippet `json:"snippet"`
}

type videoListResponse struct {
	Items         []video  `json:"items"`
	PageInfo      pageInfo `json:"pageInfo"`
	NextPageToken string   `json:"nextPageToken"`
}
