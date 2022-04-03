package mux

import (
	"github.com/de1phin/music-transfer/internal/interactor"
)

type StateHandler struct {
	state   UserState
	handler Handler
}

func NewStateHandler(state UserState, handler Handler) Handler {
	stateHandler := &StateHandler{state: state, handler: handler}
	return func(state UserState, msg interactor.Message) bool {
		if state != stateHandler.state {
			return false
		}
		return stateHandler.handler(state, msg)
	}
}

func (mux *Mux) HandleIdle(state UserState, msg interactor.Message) bool {
	if msg.Text == "add service" {
		mux.stateStorage.Put(msg.UserID, ChoosingService)
		services := ""
		for _, service := range mux.services {
			services += "\n" + service.Name()
		}
		mux.interactor.SendMessage(interactor.Message{
			UserID: msg.UserID,
			Text:   "Choose service to authorize:" + services,
		})
	} else if msg.Text == "transfer" {
		mux.stateStorage.Put(msg.UserID, ChoosingSrc)
		services := ""
		for _, service := range mux.services {
			services += "\n" + service.Name()
		}
		mux.interactor.SendMessage(interactor.Message{
			UserID: msg.UserID,
			Text:   "Choose source service:" + services,
		})
	}

	return true
}

func (mux *Mux) HandleAuthorize(state UserState, msg interactor.Message) bool {
	for _, service := range mux.services {
		if service.Name() == msg.Text {
			url := service.GetAuthURL(msg.UserID)
			mux.interactor.SendMessage(interactor.Message{
				UserID: msg.UserID,
				Text:   url,
			})
			mux.stateStorage.Put(msg.UserID, Idle)

			return true
		}
	}

	mux.interactor.SendMessage(interactor.Message{
		UserID: msg.UserID,
		Text:   "Invalid service",
	})
	return true
}

func (mux *Mux) HandleChoosingSrc(state UserState, msg interactor.Message) bool {
	for _, service := range mux.services {
		if service.Name() == msg.Text {
			mux.stateStorage.Put(msg.UserID, ChoosingDst)
			mux.transferStorage.Put(msg.UserID, service)
			services := ""
			for _, service := range mux.services {
				services += "\n" + service.Name()
			}
			mux.interactor.SendMessage(interactor.Message{
				UserID: msg.UserID,
				Text:   "Choose destination service:" + services,
			})
			return true
		}
	}

	mux.stateStorage.Put(msg.UserID, ChoosingSrc)
	mux.interactor.SendMessage(interactor.Message{
		UserID: msg.UserID,
		Text:   "Invalid service",
	})
	return true
}

func (mux *Mux) HandleChoosingDst(state UserState, msg interactor.Message) bool {
	for _, service := range mux.services {
		if service.Name() == msg.Text {
			mux.stateStorage.Put(msg.UserID, Idle)
			src := mux.transferStorage.Get(msg.UserID)
			mux.interactor.SendMessage(interactor.Message{
				UserID: msg.UserID,
				Text:   "Transfering from " + src.Name() + " to " + service.Name(),
			})
			service.AddLiked(msg.UserID, src.GetLiked(msg.UserID))
			service.AddPlaylists(msg.UserID, src.GetPlaylists(msg.UserID))
			return true
		}
	}

	mux.stateStorage.Put(msg.UserID, ChoosingDst)
	mux.interactor.SendMessage(interactor.Message{
		UserID: msg.UserID,
		Text:   "Invalid service",
	})
	return true
}
