package mux

import (
	"github.com/de1phin/music-transfer/internal/storage"
)

type Handler func(Interactor, Message) bool

type Interactor interface {
	SendMessage(Message)
	GetMessage() Message
}

type Service interface {
	Name() string
	GetAuthURL(int64) string
	GetLiked(int64) Playlist
	AddLiked(int64, Playlist)
	GetPlaylists(int64) []Playlist
	AddPlaylists(int64, []Playlist)
}

type Mux struct {
	services        []Service
	transferStorage storage.Storage[Transfer]
	interactor      Interactor
	handlers        []Handler
}

func NewMux(services []Service, interactor Interactor, transferStorage storage.Storage[Transfer]) *Mux {
	mux := Mux{
		services:        services,
		interactor:      interactor,
		transferStorage: transferStorage,
	}
	mux.handlers = []Handler{
		NewStateHandler(Idle, mux.HandleIdle),
		NewStateHandler(ChoosingService, mux.HandleAuthorize),
		NewStateHandler(ChoosingSrc, mux.HandleChoosingSrc),
		NewStateHandler(ChoosingDst, mux.HandleChoosingDst),
	}
	return &mux
}

func (mux *Mux) Run() {
	for {
		msg := mux.interactor.GetMessage()

		for _, handler := range mux.handlers {
			if handler(mux.interactor, msg) {
				break
			}
		}
	}
}
