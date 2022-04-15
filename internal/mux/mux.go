package mux

import (
	"fmt"
	"time"

	"github.com/de1phin/music-transfer/internal/log"
	"github.com/de1phin/music-transfer/internal/storage"
)

type Handler func(Interactor, Message, int64) bool

type Interactor interface {
	Name() string
	SendMessage(Message) error
	GetMessage() (Message, error)
}

type Service interface {
	Name() string
	Authorized(int64) (bool, error)
	GetAuthURL(int64) (string, error)
	GetLiked(int64) (Playlist, error)
	AddLiked(int64, Playlist) error
	GetPlaylists(int64) ([]Playlist, error)
	AddPlaylists(int64, []Playlist) error
}

type Mux struct {
	services        []Service
	transferStorage storage.Storage[int64, Transfer]
	idStorage       storage.Storage[string, int64]
	interactors     []Interactor
	handlers        []Handler
	IDGenerator     IDGenerator
	logger          log.Logger
}

func NewMux(services []Service, interactors []Interactor, transferStorage storage.Storage[int64, Transfer], idStorage storage.Storage[string, int64], logger log.Logger) *Mux {
	mux := Mux{
		services:        services,
		interactors:     interactors,
		transferStorage: transferStorage,
		idStorage:       idStorage,
		IDGenerator:     IDGenerator{nextID: 0},
		logger:          logger,
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
		go mux.handleInteractor(interactor)
	}

	<-quit
}

func (mux *Mux) handleInteractor(interactor Interactor) {
	for {
		msg, err := interactor.GetMessage()
		mux.logger.Log("Mux.handleInteractor: New message via " + interactor.Name())
		if err != nil {
			mux.logger.Log(err)
			continue
		}
		key := interactor.Name() + ":" + fmt.Sprintf("%d", msg.UserID)
		var internalID int64
		if mux.idStorage.Exist(key) {
			internalID = mux.idStorage.Get(key)
		} else {
			internalID = mux.IDGenerator.NextID()
			mux.idStorage.Put(key, internalID)
		}

		start := time.Now()
		for _, handler := range mux.handlers {
			if handler(interactor, msg, internalID) {
				break
			}
		}
		mux.logger.Log("Mux.handleInteractor: Message handled in " + time.Since(start).String())
	}
}
