package mux

import (
	"fmt"
	"strings"
	"time"
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

	return true
}

func (mux *Mux) HandleChoosingSrc(from Interactor, msg Message) bool {
	for _, service := range mux.services {
		if service.Name() == msg.Content {
			if !service.Authorized(msg.UserID) {
				from.SendMessage(Message{
					UserID:    msg.UserID,
					UserState: ChoosingDst,
					Content: "<content><text>Please log into the service:</text><url><text>" +
						strings.Title(service.Name()) + "</text><link><![CDATA[" + service.GetAuthURL(msg.UserID) +
						"]]></link></url></content>",
				})

				timeLimit := time.Now().Add(time.Second * 60)
				for !service.Authorized(msg.UserID) {
					time.Sleep(3 * time.Second)
					if timeLimit.Before(time.Now()) {
						return true
					}
				}
			}

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
			if !service.Authorized(msg.UserID) {
				from.SendMessage(Message{
					UserID:    msg.UserID,
					UserState: Idle,
					Content: "<content><text>Please log into the service:</text><url><text>" +
						strings.Title(service.Name()) + "</text><link><![CDATA[" + service.GetAuthURL(msg.UserID) +
						"]]></link></url></content>",
				})

				timeLimit := time.Now().Add(time.Second * 60)
				for !service.Authorized(msg.UserID) {
					time.Sleep(3 * time.Second)
					if timeLimit.Before(time.Now()) {
						return true
					}
				}
			}

			src := mux.transferStorage.Get(msg.UserID)

			from.SendMessage(Message{
				UserID:    msg.UserID,
				UserState: Idle,
				Content: "<content><text>Transfering from " + strings.Title(src.Name()) +
					" to " + strings.Title(service.Name()) + "</text></content>",
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
