package transfer

import (
	"strings"
)

func normalizeString(str string) string {
	return strings.ToLower(strings.Trim(str, " \n"+string(byte(13))))
}

func (transfer *Transfer) handle(chat Chat) {
	switch chat.userState {
	case Idle:
		if chat.message == "Transfer" {
			transfer.handlePickFirstService(chat)
		} else if normalizeString(chat.message) == normalizeString("Add service") {
			transfer.handleAddService(chat)
		} else {
			transfer.handleIdle(chat)
		}
	case ChoosingServiceToAdd:
		transfer.handleLogIntoService(chat)
	case LoggingIntoService:
		transfer.handleCouldntLog(chat)
	}
}

func (transfer *Transfer) handlePickFirstService(chat Chat) {

}

func (transfer *Transfer) handleLogged(chat Chat) {
	transfer.Interactor.SendMessageTo(chat.userID, "Successfully authorized")
	transfer.Storage.PutUserState(chat.userID, Idle)
}

func (transfer *Transfer) handleCouldntLog(chat Chat) {
	transfer.Interactor.SendMessageTo(chat.userID, "Couldn't authorize into the service so far, please wait or try again")
	transfer.Storage.PutUserState(chat.userID, Idle)
}

func (transfer *Transfer) handleLogIntoService(chat Chat) {
	validService := false
	for _, service := range transfer.Services {
		if normalizeString(service.Name()) == normalizeString(chat.message) {
			validService = true
			transfer.Storage.PutUserState(chat.userID, LoggingIntoService)
			transfer.Interactor.SendURL(chat.userID, "Authorize via link:\n", service.GetAuthURL(chat.userID))
			break
		}
	}

	if !validService {
		transfer.Storage.PutUserState(chat.userID, Idle)
		transfer.Interactor.SendMessageTo(chat.userID, "Invalid service")
	}
}

func (transfer *Transfer) handleAddService(chat Chat) {
	transfer.Storage.PutUserState(chat.userID, ChoosingServiceToAdd)
	names := []string{}
	for _, service := range transfer.Services {
		names = append(names, service.Name())
	}
	transfer.Interactor.ChooseFrom(chat.userID, "Choose service you want to add: \n", names)
}

func (transfer *Transfer) handleIdle(chat Chat) {
	transfer.Interactor.SendMessageTo(chat.userID, "Choose one of the options:\nTransfer\nAdd service")
}
