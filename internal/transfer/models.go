package transfer

import "net/http"

type UserState int

const (
	Idle UserState = iota
	ChoosingServiceToAdd
	LoggingIntoService
	PickingFirstService
	PickingSecondService
	Transfering
)

type User struct {
	ID          int64
	State       UserState
	ServiceFrom string
	ServiceTo   string
}

type Chat struct {
	user    User
	message string
}

type Playlist struct {
	Name  string
	Songs []Song
}

type Song struct {
	Name    string
	Artists string
}

type MusicService interface {
	Name() string
	URLName() string
	GetAuthURL(int64) string
	Authorize(callback *http.Request) (int64, interface{})
	ValidAuthCallback(callback *http.Request) bool
	GetFavourites(interface{}) Playlist
	AddFavourites(interface{}, Playlist)
	GetPlaylists(interface{}) []Playlist
	AddPlaylists(interface{}, []Playlist)
}

type Storage interface {
	HasUser(int64) bool
	GetUser(int64) User
	PutUser(User)
	AddService(string)
	PutServiceData(int64, string, interface{})
	GetServiceData(int64, string) interface{}
}

type Interactor interface {
	GetMessageFrom() (int64, string)
	SendMessageTo(int64, string)
	ChooseFrom(int64, string, []string)
	SendURL(int64, string, string)
}

type Config interface {
	GetCallbackURL() string
	GetServerURL() string
}
