package yandex

import (
	"net/http"
)

type Credentials struct {
	UID     string `db:"yandex_id"`
	Cookies string `db:"cookies"`
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

type Accounts struct {
	Users []User `json:"accounts"`
}

type User struct {
	ID          string      `json:"uid"`
	Login       string      `json:"login"`
	DisplayName DisplayName `json:"displayName"`
}

type DisplayName struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Name      string `json:"name"`
}
