package mux

import (
	"fmt"
	"time"

	"github.com/de1phin/music-transfer/internal/log"
	"github.com/de1phin/music-transfer/internal/storage"
)

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

type OnAuthorized func(Service, int64)

type Mux struct {
	services        []Service
	transferStorage storage.Storage[int64, TransferState]
	idStorage       storage.Storage[string, int64]
	interactors     []Interactor
	handlers        []handlerWrapper
	IDGenerator     IDGenerator
	logger          log.Logger
}

func NewMux(services []Service, interactors []Interactor, transferStorage storage.Storage[int64, TransferState], idStorage storage.Storage[string, int64], logger log.Logger) *Mux {
	mux := Mux{
		services:        services,
		interactors:     interactors,
		transferStorage: transferStorage,
		idStorage:       idStorage,
		IDGenerator:     IDGenerator{nextID: 0},
		logger:          logger,
	}
	mux.handlers = []handlerWrapper{
		newHandler(Idle, mux.handleIdle),
		newHandler(ChooseSource, mux.handleChooseSource),
		newHandler(AuthorizeSource, mux.handleAuthorizeSource),
		newHandler(ChooseDestination, mux.handleChooseDestination),
		newHandler(AuthorizeDestination, mux.handleAuthorizeDestination),
		newHandler(Transfer, mux.handleTransfer),
	}
	return &mux
}

func (mux *Mux) Run(quit <-chan struct{}) {
	for _, interactor := range mux.interactors {
		go mux.listenInteractor(interactor)
	}

	<-quit
}

func (mux *Mux) listenInteractor(interactor Interactor) {
	for {
		msg, err := interactor.GetMessage()
		mux.logger.Log("Mux.handleInteractor: New message via " + interactor.Name())
		if err != nil {
			mux.logger.Log(err)
			continue
		}
		key := interactor.Name() + ":" + fmt.Sprintf("%d", msg.UserID)
		var internalID int64
		ok, err := mux.idStorage.Exist(key)
		if err != nil {
			mux.logger.Log("Mux.handleInteractor:", err)
		}
		if ok {
			internalID, err = mux.idStorage.Get(key)
			if err != nil {
				mux.logger.Log("Mux.handleInteractor:", err)
			}
		} else {
			internalID = mux.IDGenerator.NextID()
			err = mux.idStorage.Set(key, internalID)
			if err != nil {
				mux.logger.Log("Mux.handleInteractor:", err)
			}
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

func (mux *Mux) transfer(from Service, to Service, userID int64) error {
	liked, err := from.GetLiked(userID)
	if err != nil {
		return fmt.Errorf("Unable to get liked: %w", err)
	}
	err = to.AddLiked(userID, liked)
	if err != nil {
		return fmt.Errorf("Unable to add liked: %w", err)
	}

	playlists, err := from.GetPlaylists(userID)
	if err != nil {
		return fmt.Errorf("Unable to get playlists: %w", err)
	}
	err = to.AddPlaylists(userID, playlists)
	if err != nil {
		return fmt.Errorf("Unable to add playlists: %w", err)
	}

	return nil
}

func (mux *Mux) OnAuthorized(from Service, internalID int64) {
	fmt.Println("Notified from", from.Name(), "about", internalID)
	transferState, err := mux.transferStorage.Get(internalID)
	if err != nil {
		mux.logger.Log("mux: OnAuthorized: Unable to get transfer state: %w", err)
		return
	}

	interactor := mux.GetInteractorByName(transferState.activeInteractorName)
	if interactor == nil {
		mux.logger.Log("mux: OnAuthorized: Active Interactor is nil")
		return
	}
	fmt.Println("Interactor:", interactor.Name())

	if transferState.sourceServiceAuthorized {
		fmt.Println("HandleAuthorizeDestination")
		mux.handleAuthorizeDestination(interactor, Message{UserID: transferState.interactorUserID}, internalID)
	} else {
		transferState.sourceServiceAuthorized = true
		err = mux.transferStorage.Set(internalID, transferState)
		if err != nil {
			mux.handleError(fmt.Errorf("OnAuthorized: Unable to put transfer state: %w", err),
				interactor, Message{
					UserID: transferState.interactorUserID,
				})
			return
		}
		fmt.Println("HandleAuthorizeSource")
		mux.handleAuthorizeSource(interactor, Message{UserID: transferState.interactorUserID}, internalID)
	}
	fmt.Println("Done")
}
