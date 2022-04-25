package mux

import (
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
		Content: MessageContent{
			Text: "An error occured",
		},
	})
}

func (mux *Mux) HandleIdle(from Interactor, msg Message, internalID int64) bool {
	services := []string{}
	for _, service := range mux.services {
		services = append(services, strings.Title(service.Name()))
	}
	from.SendMessage(Message{
		UserState: ChoosingSrc,
		UserID:    msg.UserID,
		Content: MessageContent{
			Text:    "Choose source service:",
			Buttons: services,
		},
	})

	return true
}

func (mux *Mux) HandleChoosingSrc(from Interactor, msg Message, internalID int64) bool {
	for _, service := range mux.services {
		if service.Name() == msg.Content.Text {
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
					Content: MessageContent{
						Text: "Please log into the service",
						URLs: []URL{{
							Text: strings.Title(service.Name()),
							Link: authURL,
						}},
					},
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

			mux.transferStorage.Put(internalID, service)
			services := []string{}
			for _, service := range mux.services {
				services = append(services, strings.Title(service.Name()))
			}
			from.SendMessage(Message{
				UserState: ChoosingDst,
				UserID:    msg.UserID,
				Content: MessageContent{
					Text:    "Choose destination service:",
					Buttons: services,
				},
			})
			return true
		}
	}

	from.SendMessage(Message{
		UserState: Idle,
		UserID:    msg.UserID,
		Content:   MessageContent{},
	})
	return true
}

func (mux *Mux) HandleChoosingDst(from Interactor, msg Message, internalID int64) bool {
	for _, service := range mux.services {
		if service.Name() == msg.Content.Text {
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
					Content: MessageContent{
						Text: "Please log into the service",
						URLs: []URL{{
							Text: strings.Title(service.Name()),
							Link: authURL,
						}},
					},
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

			src, err := mux.transferStorage.Get(internalID)
			if err != nil {
				mux.handleError(err, from, msg)
				return true
			}

			from.SendMessage(Message{
				UserID:    msg.UserID,
				UserState: Idle,
				Content: MessageContent{
					Text: "Transfering from " + strings.Title(src.Name()) +
						" to " + strings.Title(service.Name()),
				},
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
		Content: MessageContent{
			Text: "Invalid service",
		},
	})
	return true
}
