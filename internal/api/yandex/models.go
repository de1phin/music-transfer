package yandex

import (
	"net/http"
)

type Credentials struct {
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
	Type    string   `json:"type"`
}

type Artist struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type PlaylistResponse struct {
	Playlist Playlist `json:"playlist"`
}

type Library struct {
	Owner       Owner   `json:"owner"`
	PlaylistIDs []int64 `json:"playlistIds"`
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

type authResponse struct {
	Status string `json:"status"`
	UID    int64  `json:"default_uid"`
}
