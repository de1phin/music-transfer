package transfer

import (
	"strings"
)

func normalizeString(str string) string {

	return strings.ToLower(strings.Trim(str, " \n"+string(byte(13))))
}

func (transfer *Transfer) handle(id int64, userState UserState, message string) {
	if userState == Idle {
		if message == "Transfer" {
			transfer.handlePickFirstService(id, userState, message)
		} else if normalizeString(message) == normalizeString("Add service") {
			transfer.handleAddService(id, userState, message)
		} else {
			transfer.handleIdle(id, userState, message)
		}
	} else if userState == ChoosingServiceToAdd {
		transfer.handleLogIntoService(id, userState, message)
	}
}

func (transfer *Transfer) handlePickFirstService(id int64, userState UserState, message string) {

}

func (transfer *Transfer) handleLogIntoService(id int64, userState UserState, message string) {
	validService := false
	for _, service := range transfer.Services {
		if normalizeString(service.Name()) == normalizeString(message) {
			validService = true
			transfer.Storage.PutUserState(id, LoggingIntoService)
			transfer.Interactor.SendURL(id, "Authorize via link:\n", service.GetAuthURL(id))
			break
		}
	}

	if !validService {
		transfer.Storage.PutUserState(id, Idle)
		transfer.Interactor.SendMessageTo(id, "Invalid service")
	}
}

func (transfer *Transfer) handleAddService(id int64, userState UserState, message string) {
	transfer.Storage.PutUserState(id, ChoosingServiceToAdd)
	names := []string{}
	for _, service := range transfer.Services {
		names = append(names, service.Name())
	}
	transfer.Interactor.ChooseFrom(id, "Choose service you want to add: \n", names)
}

func (transfer *Transfer) handleIdle(id int64, userState UserState, message string) {
	transfer.Interactor.SendMessageTo(id, "Choose one of the options:\nTransfer\nAdd service")
}
