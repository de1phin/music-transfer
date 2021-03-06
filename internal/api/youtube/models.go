package youtube

type Credentials struct {
	AccessToken string `json:"access_token" db:"access_token"`
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

type PlaylistItem struct {
	ID      string  `json:"id"`
	Snippet snippet `json:"snippet"`
}

type playlistItemListResponse struct {
	NextPageToken string         `json:"nextPageToken"`
	PageInfo      pageInfo       `json:"pageInfo"`
	Items         []PlaylistItem `json:"items"`
}

type resourceID struct {
	Kind    string `json:"kind"`
	VideoID string `json:"videoId"`
}

type addVideosRequest struct {
	Snippet snippet `json:"snippet"`
}

type snippet struct {
	PlaylistID             string     `json:"playlistId"`
	Title                  string     `json:"title"`
	ChannelTitle           string     `json:"channelTitle"`
	VideoOwnerChannelTitle string     `json:"videoOwnerChannelTitle"`
	ResourceID             resourceID `json:"resourceId"`
}

type Video struct {
	ID      string  `json:"id"`
	Snippet snippet `json:"snippet"`
}

type mergedID struct {
	Kind       string `json:"kind"`
	VideoID    string `json:"videoId"`
	PlaylistID string `json:"playlistId"`
	ChannelID  string `json:"channelId"`
}

type SearchResult struct {
	ID      mergedID `json:"id"`
	Snippet snippet  `json:"snippet"`
}

type searchListResponse struct {
	PageInfo      pageInfo       `json:"pageInfo"`
	NextPageToken string         `json:"nextPageToken"`
	Items         []SearchResult `json:"items"`
}

type videoListResponse struct {
	Items         []Video  `json:"items"`
	PageInfo      pageInfo `json:"pageInfo"`
	NextPageToken string   `json:"nextPageToken"`
}
