package transfer

import "net/http"

type UserState int

const (
	Idle UserState = iota
	ChoosingServiceToAdd
	LoggingIntoService
)

type Chat struct {
	userID    int64
	userState UserState
	message   string
}

type MusicService interface {
	Name() string
	URLName() string
	GetAuthURL(int64) string
	Authorize(callback *http.Request) (int64, interface{})
	ValidAuthCallback(callback *http.Request) bool
}

type Storage interface {
	GetUserState(int64) UserState
	PutUserState(int64, UserState)
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
