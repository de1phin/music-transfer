package mux

import (
	"fmt"

	"github.com/de1phin/music-transfer/internal/storage"
)

type Handler func(Interactor, Message, int64) bool

type Interactor interface {
	Name() string
	SendMessage(Message)
	GetMessage() Message
}

type Service interface {
	Name() string
	Authorized(int64) bool
	GetAuthURL(int64) string
	GetLiked(int64) Playlist
	AddLiked(int64, Playlist)
	GetPlaylists(int64) []Playlist
	AddPlaylists(int64, []Playlist)
}

type Mux struct {
	services        []Service
	transferStorage storage.Storage[int64, Transfer]
	idStorage       storage.Storage[string, int64]
	interactors     []Interactor
	handlers        []Handler
	IDGenerator     IDGenerator
}

func NewMux(services []Service, interactors []Interactor, transferStorage storage.Storage[int64, Transfer], idStorage storage.Storage[string, int64]) *Mux {
	mux := Mux{
		services:        services,
		interactors:     interactors,
		transferStorage: transferStorage,
		idStorage:       idStorage,
		IDGenerator:     IDGenerator{nextID: 0},
	}
	mux.handlers = []Handler{
		NewStateHandler(Idle, mux.HandleIdle),
		NewStateHandler(ChoosingSrc, mux.HandleChoosingSrc),
		NewStateHandler(ChoosingDst, mux.HandleChoosingDst),
	}
	return &mux
}

func (mux *Mux) Run(quit <-chan struct{}) {
	for _, interactor := range mux.interactors {
		go func(interactor Interactor) {
			for {
				msg := interactor.GetMessage()
				key := interactor.Name() + ":" + fmt.Sprintf("%d", msg.UserID)
				var internalID int64
				if mux.idStorage.Exist(key) {
					internalID = mux.idStorage.Get(key)
				} else {
					internalID = mux.IDGenerator.NextID()
					mux.idStorage.Put(key, internalID)
				}

				for _, handler := range mux.handlers {
					if handler(interactor, msg, internalID) {
						break
					}
				}
			}
		}(interactor)
	}

	<-quit
}
