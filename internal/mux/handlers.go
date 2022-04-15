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

func (mux *Mux) handleError(err error, from Interactor, msg Message) {
	mux.logger.Log("Mux.handleError:", err)
	from.SendMessage(Message{
		UserID:    msg.UserID,
		UserState: msg.UserState,
		Content:   "<content><text>An error occured</text></content>",
	})
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
			authorized, err := service.Authorized(internalID)
			if err != nil {
				mux.handleError(err, from, msg)
				return false
			}
			if !authorized {
				authURL, err := service.GetAuthURL(internalID)
				if err != nil {
					mux.handleError(err, from, msg)
					return false
				}

				from.SendMessage(Message{
					UserID:    msg.UserID,
					UserState: ChoosingDst,
					Content: "<content><text>Please log into the service:</text><url><text>" +
						strings.Title(service.Name()) + "</text><link><![CDATA[" + authURL +
						"]]></link></url></content>",
				})
				log.Println("wait for", internalID)
				timeLimit := time.Now().Add(time.Second * 60)
				for {
					time.Sleep(3 * time.Second)
					if timeLimit.Before(time.Now()) {
						return true
					}
					authorized, err := service.Authorized(internalID)

					if err != nil {
						mux.handleError(err, from, msg)
						return false
					}
					if authorized {
						break
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
			authorized, err := service.Authorized(internalID)
			if err != nil {
				mux.handleError(err, from, msg)
				return false
			}
			if !authorized {
				authURL, err := service.GetAuthURL(internalID)
				if err != nil {
					mux.handleError(err, from, msg)
					return false
				}

				from.SendMessage(Message{
					UserID:    msg.UserID,
					UserState: Idle,
					Content: "<content><text>Please log into the service:</text><url><text>" +
						strings.Title(service.Name()) + "</text><link><![CDATA[" + authURL +
						"]]></link></url></content>",
				})

				timeLimit := time.Now().Add(time.Second * 60)
				for {
					time.Sleep(3 * time.Second)
					if timeLimit.Before(time.Now()) {
						return true
					}
					authorized, err := service.Authorized(internalID)
					if err != nil {
						mux.handleError(err, from, msg)
						return false
					}
					if authorized {
						break
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

			liked, err := src.GetLiked(internalID)
			if err != nil {
				mux.handleError(err, from, msg)
				return false
			}
			err = service.AddLiked(internalID, liked)
			if err != nil {
				mux.handleError(err, from, msg)
				return false
			}
			playlists, err := src.GetPlaylists(internalID)
			if err != nil {
				mux.handleError(err, from, msg)
				return false
			}
			err = service.AddPlaylists(internalID, playlists)
			if err != nil {
				mux.handleError(err, from, msg)
				return false
			}

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
