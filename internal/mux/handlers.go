package mux

import (
	"fmt"
	"strings"
)

type StateHandler struct {
	state   UserState
	handler Handler
}

func NewStateHandler(state UserState, handler Handler) Handler {
	stateHandler := &StateHandler{state: state, handler: handler}
	return func(state UserState, msg Message) bool {
		if state != stateHandler.state {
			return false
		}
		return stateHandler.handler(state, msg)
	}
}

func (mux *Mux) HandleIdle(state UserState, msg Message) bool {
	if msg.Content == "add service" {
		mux.stateStorage.Put(msg.UserID, ChoosingService)
		services := ""
		for _, service := range mux.services {
			services += fmt.Sprintf(`
	<either>
		<text>%s</text>
		<button><text>%s</text><metadata>%s</metadata></button>
	</either>`,
				strings.Title(service.Name()), strings.Title(service.Name()), strings.Title(service.Name()))
		}
		mux.interactor.SendMessage(Message{
			UserID:  msg.UserID,
			Content: "<content><text>Choose service to log in:</text>\n" + services + "</content>",
		})
	} else if msg.Content == "transfer" {
		mux.stateStorage.Put(msg.UserID, ChoosingSrc)
		services := ""
		for _, service := range mux.services {
			services += fmt.Sprintf(`
	<either>
		<text>%s</text>
		<button><text>%s</text><metadata>%s</metadata></button>
	</either>`,
				strings.Title(service.Name()), strings.Title(service.Name()), strings.Title(service.Name()))
		}
		mux.interactor.SendMessage(Message{
			UserID:  msg.UserID,
			Content: "<content>\n<text>Choose source service:</text>\n" + services + "</content>",
		})
	}

	return true
}

func (mux *Mux) HandleAuthorize(state UserState, msg Message) bool {
	for _, service := range mux.services {
		if service.Name() == msg.Content {
			url := service.GetAuthURL(msg.UserID)
			mux.stateStorage.Put(msg.UserID, Idle)
			mux.interactor.SendMessage(Message{
				UserID:  msg.UserID,
				Content: "<content><text>Please follow the link:</text><url><text>Log in</text><link><![CDATA[" + url + "]]></link></url></content>",
			})
			return true
		}
	}

	mux.interactor.SendMessage(Message{
		UserID:  msg.UserID,
		Content: "<content><text>Invalid service</text></content>",
	})
	return true
}

func (mux *Mux) HandleChoosingSrc(state UserState, msg Message) bool {
	for _, service := range mux.services {
		if service.Name() == msg.Content {
			mux.stateStorage.Put(msg.UserID, ChoosingDst)
			mux.transferStorage.Put(msg.UserID, service)
			services := ""
			for _, service := range mux.services {
				services += fmt.Sprintf(`
		<either>
			<text>%s</text>
			<button><text>%s</text><metadata>%s</metadata></button>
		</either>`,
					strings.Title(service.Name()), strings.Title(service.Name()), strings.Title(service.Name()))
			}
			mux.interactor.SendMessage(Message{
				UserID:  msg.UserID,
				Content: "<content><text>Choose destination service:</text>\n" + services + "</content>",
			})
			return true
		}
	}

	mux.stateStorage.Put(msg.UserID, ChoosingSrc)
	mux.interactor.SendMessage(Message{
		UserID:  msg.UserID,
		Content: "<content><text>Invalid service</text></content>",
	})
	return true
}

func (mux *Mux) HandleChoosingDst(state UserState, msg Message) bool {
	for _, service := range mux.services {
		if service.Name() == msg.Content {
			mux.stateStorage.Put(msg.UserID, Idle)
			src := mux.transferStorage.Get(msg.UserID)
			mux.interactor.SendMessage(Message{
				UserID:  msg.UserID,
				Content: "<content><text>Transfering from " + src.Name() + " to " + service.Name() + "</text></content>",
			})
			service.AddLiked(msg.UserID, src.GetLiked(msg.UserID))
			service.AddPlaylists(msg.UserID, src.GetPlaylists(msg.UserID))
			return true
		}
	}

	mux.stateStorage.Put(msg.UserID, ChoosingDst)
	mux.interactor.SendMessage(Message{
		UserID:  msg.UserID,
		Content: "<content><text>Invalid service</text></content>",
	})
	return true
}
