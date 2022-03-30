package transfer

func (transfer *Transfer) handle(chat Chat) {
	switch chat.user.State {
	case Idle:
		if normalizeString(chat.message) == normalizeString("Transfer") {
			transfer.handleStartTransfer(chat)
		} else if normalizeString(chat.message) == normalizeString("Add service") {
			transfer.handleAddService(chat)
		} else {
			transfer.handleIdle(chat)
		}
	case ChoosingServiceToAdd:
		transfer.handleLogIntoService(chat)
	case LoggingIntoService:
		transfer.handleCouldntLog(chat)
	case PickingFirstService:
		transfer.handlePickFirstService(chat)
	case PickingSecondService:
		transfer.handlePickSecondService(chat)
	}
}

func (transfer *Transfer) handlePickSecondService(chat Chat) {
	service, err := transfer.getServiceByName(chat.message)
	if err != nil {
		chat.user.State = Idle
		transfer.Storage.PutUser(chat.user)
		transfer.Interactor.SendMessageTo(chat.user.ID, "Invalid service")
		return
	}

	chat.user.State = Transfering
	chat.user.ServiceTo = service.Name()
	transfer.Storage.PutUser(chat.user)
	serviceFrom, _ := transfer.getServiceByName(chat.user.ServiceFrom)
	transfer.Transfer(chat.user, serviceFrom, service)
}

func (transfer *Transfer) handlePickFirstService(chat Chat) {
	service, err := transfer.getServiceByName(chat.message)
	if err != nil {
		chat.user.State = Idle
		transfer.Storage.PutUser(chat.user)
		transfer.Interactor.SendMessageTo(chat.user.ID, "Invalid service")
		return
	}

	chat.user.State = PickingSecondService
	chat.user.ServiceFrom = service.Name()
	transfer.Storage.PutUser(chat.user)
	names := []string{}
	for _, service := range transfer.Services {
		names = append(names, service.Name())
	}
	transfer.Interactor.ChooseFrom(chat.user.ID, "Pick service you want to transfer to:\n", names)
}

func (transfer *Transfer) handleStartTransfer(chat Chat) {
	changedUser := chat.user
	changedUser.State = PickingFirstService
	transfer.Storage.PutUser(changedUser)
	names := []string{}
	for _, service := range transfer.Services {
		names = append(names, service.Name())
	}
	transfer.Interactor.ChooseFrom(chat.user.ID, "Pick service you want to transfer from:\n", names)
}

func (transfer *Transfer) handleLogged(chat Chat) {
	transfer.Interactor.SendMessageTo(chat.user.ID, "Successfully authorized")
	changedUser := chat.user
	changedUser.State = Idle
	transfer.Storage.PutUser(changedUser)
}

func (transfer *Transfer) handleCouldntLog(chat Chat) {
	transfer.Interactor.SendMessageTo(chat.user.ID, "Couldn't authorize into the service so far, please wait or try again")
	changedUser := chat.user
	changedUser.State = Idle
	transfer.Storage.PutUser(changedUser)
}

func (transfer *Transfer) handleLogIntoService(chat Chat) {
	validService := false
	for _, service := range transfer.Services {
		if normalizeString(service.Name()) == normalizeString(chat.message) {
			validService = true
			changedUser := chat.user
			changedUser.State = LoggingIntoService
			transfer.Storage.PutUser(changedUser)
			transfer.Interactor.SendURL(chat.user.ID, "Authorize via link:\n", service.GetAuthURL(chat.user.ID))
			break
		}
	}

	if !validService {
		changedUser := chat.user
		changedUser.State = Idle
		transfer.Storage.PutUser(changedUser)
		transfer.Interactor.SendMessageTo(chat.user.ID, "Invalid service")
	}
}

func (transfer *Transfer) handleAddService(chat Chat) {
	changedUser := chat.user
	changedUser.State = ChoosingServiceToAdd
	transfer.Storage.PutUser(changedUser)
	names := []string{}
	for _, service := range transfer.Services {
		names = append(names, service.Name())
	}
	transfer.Interactor.ChooseFrom(chat.user.ID, "Choose service you want to add: \n", names)
}

func (transfer *Transfer) handleIdle(chat Chat) {
	transfer.Interactor.SendMessageTo(chat.user.ID, "Choose one of the options:\nTransfer\nAdd service")
}
