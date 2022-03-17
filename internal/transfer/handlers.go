package transfer

import (
	"strings"
)

func normalizeString(str string) string {

	return strings.ToLower(strings.Trim(str, " \n"+string(13)))
}

func (transfer *Transfer) handle(id int64, userState UserState, message string) {
	if userState == Idle {
		if message == "Transfer" {
		} else if normalizeString(message) == normalizeString("Add service") {
			transfer.handleAddService(id, userState, message)
		} else {
			transfer.handleIdle(id, userState, message)
		}
	} else if userState == ChoosingServiceToAdd {
		transfer.handleLogIntoService(id, userState, message)
	}

}

func (transfer *Transfer) handleLogIntoService(id int64, userState UserState, message string) {
	validService := false
	for _, service := range transfer.services {
		if normalizeString(service.Name()) == normalizeString(message) {
			validService = true
			transfer.storage.PutUserState(id, LoggingIntoService)
			transfer.interactor.SendURL(id, "Authorize via link:\n", service.GetAuthURL(id, transfer.callbackURL))
			break
		}
	}

	if !validService {
		transfer.storage.PutUserState(id, Idle)
		transfer.interactor.SendMessageTo(id, "Invalid service")
	}
}

func (transfer *Transfer) handleAddService(id int64, userState UserState, message string) {
	transfer.storage.PutUserState(id, ChoosingServiceToAdd)
	names := []string{}
	for _, service := range transfer.services {
		names = append(names, service.Name())
	}
	transfer.interactor.ChooseFrom(id, "Choose service you want to add: \n", names)
}

func (transfer *Transfer) handleIdle(id int64, userState UserState, message string) {
	transfer.interactor.SendMessageTo(id, "Choose one of the options:\nTransfer\nAdd service")
}
