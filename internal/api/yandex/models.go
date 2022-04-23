package yandex

import (
	"net/http"
)

type Credentials struct {
	Login   string `db:"yandex_login"`
	UID     string `db:"yandex_id"`
	Cookies string `db:"cookies"`
}

type Accounts struct {
	Users []User `json:"accounts"`
}

type User struct {
	ID          string      `json:"uid"`
	Login       string      `json:"login"`
	DisplayName DisplayName `json:"displayName"`
}

type Owner struct {
	UID   string `json:"uid"`
	Name  string `json:"name"`
	Login string `json:"login"`
}

type DisplayName struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Name      string `json:"name"`
}

type playlistPatchDifference struct {
	Operation string         `json:"op"`
	At        int            `json:"at"`
	Tracks    []TrackSnippet `json:"tracks"`
}

type playlistAddResponse struct {
	Playlist PlaylistSnippet `json:"playlist"`
}

type PlaylistSnippet struct {
	Title      string `json:"title"`
	Kind       int64  `json:"kind"`
	TrackCount int64  `json:"trackCount"`
}

type Playlist struct {
	Title      string   `json:"title"`
	Kind       int      `json:"kind"`
	TrackCount int      `json:"trackCount"`
	TrackIDs   []string `json:"trackIds"`
	Tracks     []Track  `json:"tracks"`
}

type Track struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Artists []Artist `json:"artists"`
	Albums  []Album  `json:"albums"`
	Type    string   `json:"type"`
}

type TrackSnippet struct {
	ID      string `json:"id"`
	AlbumID int64  `json:"albumId"`
}

type Artist struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Library struct {
	Owner       Owner      `json:"owner"`
	Playlists   []Playlist `json:"playlists"`
	PlaylistIDs []int64    `json:"playlistIds"`
}

type PlaylistResponse struct {
	Playlist Playlist `json:"playlist"`
}

type Album struct {
	Title string `json:"title"`
	ID    int64  `json:"id"`
}

type TrackItem struct {
	ID     int64    `json:"id"`
	Title  string   `json:"title"`
	Artist []Artist `json:"artists"`
	Albums []Album  `json:"albums"`
	Type   string   `json:"music"`
}

type TrackItemArray struct {
	Items []TrackItem `json:"items"`
}

type SearchResponse struct {
	Tracks TrackItemArray `json:"tracks"`
}

type loginForm struct {
	Sender   string `json:"sender"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginFormTokens struct {
	csrf        string
	processUUID string
	cookies     []*http.Cookie
}

type submitResponse struct {
	CsrfToken string `json:"csrf_token"`
	Status    string `json:"status"`
	TrackID   string `json:"track_id"`
}

type authTokensResponse struct {
	User AuthTokens `json:"user"`
}

type AuthTokens struct {
	Sign string `json:"sign"`
}

type authResponse struct {
	Status string `json:"status"`
	UID    int64  `json:"default_uid"`
}
