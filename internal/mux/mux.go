package mux

import (
	"github.com/de1phin/music-transfer/internal/interactor"
	"github.com/de1phin/music-transfer/internal/storage"
)

type Handler func(UserState, interactor.Message) bool

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
	stateStorage    storage.Storage[UserState]
	transferStorage storage.Storage[Transfer]
	interactor      interactor.InteractorSpec
	handlers        []Handler
}

func NewMux(services []Service, interactor interactor.InteractorSpec, stateStorage storage.Storage[UserState], transferStorage storage.Storage[Transfer]) *Mux {
	mux := Mux{
		services:        services,
		interactor:      interactor,
		stateStorage:    stateStorage,
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
		userState := Idle
		if !mux.stateStorage.Exist(msg.UserID) {
			mux.stateStorage.Put(msg.UserID, Idle)
		} else {
			userState = mux.stateStorage.Get(msg.UserID)
		}

		for _, handler := range mux.handlers {
			if handler(userState, msg) {
				break
			}
		}
	}
}
