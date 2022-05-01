package mux

import (
	"fmt"
	"strings"
	"time"
)

type handlerWrapper func(Interactor, Message, int64) bool
type handler func(Interactor, Message, int64)

func newHandler(state UserState, handler handler) handlerWrapper {
	return func(from Interactor, msg Message, internalID int64) bool {
		if msg.UserState != state {
			return false
		}

		handler(from, msg, internalID)
		return true
	}
}

func (mux *Mux) handleError(err error, from Interactor, msg Message) {
	mux.logger.Log(fmt.Errorf("mux: Unable to handle %v: %w", msg, err))
	err = from.SendMessage(Message{
		UserID:    msg.UserID,
		UserState: Idle,
		Content:   MessageContent{Text: "An error occured"},
	})
	if err != nil {
		mux.logger.Log(fmt.Errorf("mux: Unable to send message: %w", err))
		return
	}
}

func (mux *Mux) handleIdle(from Interactor, msg Message, internalID int64) {
	err := from.SendMessage(Message{
		UserID:    msg.UserID,
		UserState: ChooseSource,
		Content: MessageContent{
			Text:    "Choose source service:",
			Buttons: mux.GetServicesNames(),
		},
	})
	if err != nil {
		mux.handleError(fmt.Errorf("Unable to send message: %w", err), from, msg)
		return
	}
}

func (mux *Mux) handleChooseSource(from Interactor, msg Message, internalID int64) {
	source := mux.GetServiceByName(msg.Content.Text)
	if source == nil {
		err := from.SendMessage(Message{
			UserID:    msg.UserID,
			UserState: ChooseSource,
			Content: MessageContent{
				Text: "Invalid service",
			},
		})
		if err != nil {
			mux.handleError(fmt.Errorf("Unable to send message: %w", err), from, msg)
			return
		}
		return
	}

	authorized, err := source.Authorized(internalID)
	if err != nil {
		mux.handleError(fmt.Errorf("Unable to check authorization: %w", err), from, msg)
		return
	}

	if authorized {
		err = mux.transferStorage.Set(internalID, TransferState{
			sourceServiceName:       source.Name(),
			sourceServiceAuthorized: true,
			activeInteractorName:    from.Name(),
			interactorUserID:        msg.UserID,
		})
		if err != nil {
			mux.handleError(fmt.Errorf("Unable to put transfer state: %w", err), from, msg)
			return
		}
		mux.handleAuthorizeSource(from, msg, internalID)
	} else {
		url, err := source.GetAuthURL(internalID)
		if err != nil {
			mux.handleError(fmt.Errorf("Unable to get auth url: %w", err), from, msg)
			return
		}
		err = mux.transferStorage.Set(internalID, TransferState{
			sourceServiceName:       source.Name(),
			sourceServiceAuthorized: false,
			activeInteractorName:    from.Name(),
			interactorUserID:        msg.UserID,
		})
		if err != nil {
			mux.handleError(fmt.Errorf("Unable to put transfer state: %w", err), from, msg)
			return
		}
		err = from.SendMessage(Message{
			UserID:    msg.UserID,
			UserState: AuthorizeSource,
			Content: MessageContent{
				Text: "Please log into the service",
				URLs: []URL{{
					Text: strings.Title(source.Name()),
					Link: url,
				}},
			},
		})
		if err != nil {
			mux.handleError(fmt.Errorf("Unable to send message: %w", err), from, msg)
			return
		}
	}
}

func (mux *Mux) handleAuthorizeSource(from Interactor, msg Message, internalID int64) {
	err := from.SendMessage(Message{
		UserID:    msg.UserID,
		UserState: ChooseDestination,
		Content: MessageContent{
			Text:    "Choose destination service:",
			Buttons: mux.GetServicesNames(),
		},
	})
	if err != nil {
		mux.handleError(fmt.Errorf("Unable to send message: %w", err), from, msg)
		return
	}
}

func (mux *Mux) handleChooseDestination(from Interactor, msg Message, internalID int64) {
	destination := mux.GetServiceByName(msg.Content.Text)
	if destination == nil {
		err := from.SendMessage(Message{
			UserID:    msg.UserID,
			UserState: ChooseDestination,
			Content: MessageContent{
				Text: "Invalid service",
			},
		})
		if err != nil {
			mux.handleError(fmt.Errorf("Unable to send message: %w", err), from, msg)
			return
		}
		return
	}

	authorized, err := destination.Authorized(internalID)
	if err != nil {
		mux.handleError(fmt.Errorf("Unable to check authorization: %w", err), from, msg)
		return
	}

	if authorized {
		transferState, err := mux.transferStorage.Get(internalID)
		if err != nil {
			mux.handleError(fmt.Errorf("Unable to get transfer state: %w", err), from, msg)
			return
		}
		transferState.destinationServiceName = destination.Name()
		err = mux.transferStorage.Set(internalID, transferState)
		if err != nil {
			mux.handleError(fmt.Errorf("Unable to put transfer state: %w", err), from, msg)
			return
		}
		mux.handleAuthorizeDestination(from, msg, internalID)
	} else {
		url, err := destination.GetAuthURL(internalID)
		if err != nil {
			mux.handleError(fmt.Errorf("Unable to get auth url: %w", err), from, msg)
			return
		}
		err = mux.transferStorage.Set(internalID, TransferState{
			sourceServiceName:       destination.Name(),
			sourceServiceAuthorized: false,
			activeInteractorName:    from.Name(),
			interactorUserID:        msg.UserID,
		})
		if err != nil {
			mux.handleError(fmt.Errorf("Unable to put transfer state: %w", err), from, msg)
			return
		}
		err = from.SendMessage(Message{
			UserID:    msg.UserID,
			UserState: AuthorizeDestination,
			Content: MessageContent{
				Text: "Please log into the service",
				URLs: []URL{{
					Text: strings.Title(destination.Name()),
					Link: url,
				}},
			},
		})
		if err != nil {
			mux.handleError(fmt.Errorf("Unable to send message: %w", err), from, msg)
			return
		}
	}
}

func (mux *Mux) handleAuthorizeDestination(from Interactor, msg Message, internalID int64) {
	err := from.SendMessage(Message{
		UserID:    msg.UserID,
		UserState: Transfer,
		Content: MessageContent{
			Text: "Transfering",
		},
	})
	if err != nil {
		mux.handleError(fmt.Errorf("Unable to send message: %w", err), from, msg)
		return
	}

	start := time.Now()
	transferState, err := mux.transferStorage.Get(internalID)
	if err != nil {
		mux.handleError(fmt.Errorf("Unable to get transfer state: %w", err), from, msg)
		return
	}
	source := mux.GetServiceByName(transferState.sourceServiceName)
	if source == nil {
		mux.handleError(fmt.Errorf("Source service is nil: %w", err), from, msg)
		return
	}
	destination := mux.GetServiceByName(transferState.destinationServiceName)
	if destination == nil {
		mux.handleError(fmt.Errorf("Destination service is nil: %w", err), from, msg)
		return
	}
	err = mux.transfer(source, destination, internalID)
	if err != nil {
		mux.handleError(fmt.Errorf("Unable to transfer: %w", err), from, msg)
		return
	}

	err = from.SendMessage(Message{
		UserID:    msg.UserID,
		UserState: Idle,
		Content: MessageContent{
			Text: "Finished transfering",
		},
	})
	if err != nil {
		mux.handleError(fmt.Errorf("Unable to send message: %w", err), from, msg)
		return
	}

	mux.logger.Log("mux: Transfer done in", time.Since(start))
}

func (mux *Mux) handleTransfer(from Interactor, msg Message, internalID int64) {
	err := from.SendMessage(Message{
		UserID:    msg.UserID,
		UserState: Transfer,
		Content: MessageContent{
			Text: "Still transfering, please wait",
		},
	})
	if err != nil {
		mux.handleError(fmt.Errorf("Unable to send message: %w", err), from, msg)
		return
	}
}
