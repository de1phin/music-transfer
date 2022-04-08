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
	return func(from Interactor, msg Message) bool {
		if msg.UserState != stateHandler.state {
			return false
		}
		return stateHandler.handler(from, msg)
	}
}

func (mux *Mux) HandleIdle(from Interactor, msg Message) bool {
	if msg.Content == "add service" {
		services := ""
		for _, service := range mux.services {
			services += fmt.Sprintf(`
	<either>
		<text>%s</text>
		<button>%s</button>
	</either>`,
				strings.Title(service.Name()), strings.Title(service.Name()))
		}
		from.SendMessage(Message{
			UserState: ChoosingService,
			UserID:    msg.UserID,
			Content:   "<content><text>Choose service to log in:</text>\n" + services + "</content>",
		})
	} else if msg.Content == "transfer" {
		services := ""
		for _, service := range mux.services {
			services += fmt.Sprintf(`
	<either>
		<text>%s</text>
		<button>%s</button>
	</either>`,
				strings.Title(service.Name()), strings.Title(service.Name()))
		}
		from.SendMessage(Message{
			UserState: ChoosingSrc,
			UserID:    msg.UserID,
			Content:   "<content>\n<text>Choose source service:</text>\n" + services + "</content>",
		})
	}

	return true
}

func (mux *Mux) HandleAuthorize(from Interactor, msg Message) bool {
	for _, service := range mux.services {
		if service.Name() == msg.Content {
			url := service.GetAuthURL(msg.UserID)
			from.SendMessage(Message{
				UserState: Idle,
				UserID:    msg.UserID,
				Content:   "<content><text>Please follow the link:</text><url><text>Log in</text><link><![CDATA[" + url + "]]></link></url></content>",
			})
			return true
		}
	}

	from.SendMessage(Message{
		UserState: Idle,
		UserID:    msg.UserID,
		Content:   "<content><text>Invalid service</text></content>",
	})
	return true
}

func (mux *Mux) HandleChoosingSrc(from Interactor, msg Message) bool {
	for _, service := range mux.services {
		if service.Name() == msg.Content {
			mux.transferStorage.Put(msg.UserID, service)
			services := ""
			for _, service := range mux.services {
				services += fmt.Sprintf(`
		<either>
			<text>%s</text>
			<button>%s</button>
		</either>`,
					strings.Title(service.Name()), strings.Title(service.Name()))
			}
			from.SendMessage(Message{
				UserState: ChoosingDst,
				UserID:    msg.UserID,
				Content:   "<content><text>Choose destination service:</text>\n" + services + "</content>",
			})
			return true
		}
	}

	from.SendMessage(Message{
		UserState: Idle,
		UserID:    msg.UserID,
		Content:   "<content><text>Invalid service</text></content>",
	})
	return true
}

func (mux *Mux) HandleChoosingDst(from Interactor, msg Message) bool {
	for _, service := range mux.services {
		if service.Name() == msg.Content {
			src := mux.transferStorage.Get(msg.UserID)
			from.SendMessage(Message{
				UserState: Idle,
				UserID:    msg.UserID,
				Content:   "<content><text>Transfering from " + src.Name() + " to " + service.Name() + "</text></content>",
			})
			service.AddLiked(msg.UserID, src.GetLiked(msg.UserID))
			service.AddPlaylists(msg.UserID, src.GetPlaylists(msg.UserID))
			return true
		}
	}

	from.SendMessage(Message{
		UserState: Idle,
		UserID:    msg.UserID,
		Content:   "<content><text>Invalid service</text></content>",
	})
	return true
}
