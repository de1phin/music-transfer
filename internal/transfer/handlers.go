package transfer

func (transfer *Transfer) handle(id int64, userState UserState, message string) {
	if userState == Idle {
		if message == "Transfer" {
			transfer.handlePickFirstService(id, userState, message)
		} else if message == "Add service" {
			transfer.handleAddService(id, userState, message)
		} else {
			transfer.handleIdle(id, userState, message)
		}
	}
}

func (transfer *Transfer) handlePickFirstService(id int64, userState UserState, message string) {

}

func (transfer *Transfer) handleAddService(id int64, userState UserState, message string) {

}

func (transfer *Transfer) handleIdle(id int64, userState UserState, message string) {
	transfer.interactor.SendMessageTo(id, "Choose one of the options:\nTransfer\nAdd service")
}
