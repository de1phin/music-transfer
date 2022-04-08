package mux

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type StateHandler struct {
	state   UserState
	handler Handler
}

func NewStateHandler(state UserState, handler Handler) Handler {
	stateHandler := &StateHandler{state: state, handler: handler}
	return func(from Interactor, msg Message, internalID int64) bool {
		if msg.UserState != stateHandler.state {
			return false
		}
		return stateHandler.handler(from, msg, internalID)
	}
}

func (mux *Mux) HandleIdle(from Interactor, msg Message, internalID int64) bool {
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

func (mux *Mux) HandleChoosingSrc(from Interactor, msg Message, internalID int64) bool {
	for _, service := range mux.services {
		if service.Name() == msg.Content {
			if !service.Authorized(internalID) {
				from.SendMessage(Message{
					UserID:    msg.UserID,
					UserState: ChoosingDst,
					Content: "<content><text>Please log into the service:</text><url><text>" +
						strings.Title(service.Name()) + "</text><link><![CDATA[" + service.GetAuthURL(internalID) +
						"]]></link></url></content>",
				})
				log.Println("wait for", internalID)
				timeLimit := time.Now().Add(time.Second * 60)
				for !service.Authorized(internalID) {
					time.Sleep(3 * time.Second)
					if timeLimit.Before(time.Now()) {
						return true
					}
				}
				log.Println("stop")
			}

			mux.transferStorage.Put(internalID, service)
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

func (mux *Mux) HandleChoosingDst(from Interactor, msg Message, internalID int64) bool {
	for _, service := range mux.services {
		if service.Name() == msg.Content {
			if !service.Authorized(internalID) {
				from.SendMessage(Message{
					UserID:    msg.UserID,
					UserState: Idle,
					Content: "<content><text>Please log into the service:</text><url><text>" +
						strings.Title(service.Name()) + "</text><link><![CDATA[" + service.GetAuthURL(internalID) +
						"]]></link></url></content>",
				})

				timeLimit := time.Now().Add(time.Second * 60)
				for !service.Authorized(internalID) {
					time.Sleep(3 * time.Second)
					if timeLimit.Before(time.Now()) {
						return true
					}
				}
			}

			src := mux.transferStorage.Get(internalID)

			from.SendMessage(Message{
				UserID:    msg.UserID,
				UserState: Idle,
				Content: "<content><text>Transfering from " + strings.Title(src.Name()) +
					" to " + strings.Title(service.Name()) + "</text></content>",
			})

			service.AddLiked(internalID, src.GetLiked(internalID))
			service.AddPlaylists(internalID, src.GetPlaylists(internalID))

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
