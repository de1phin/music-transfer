package transfer

type Transfer struct {
	Interactor Interactor
	Storage    Storage
	Services   []MusicService
	Config     Config
}

func (transfer *Transfer) Run() {

	transfer.SetUpCallbackServers()

	for {
		id, message := transfer.Interactor.GetMessageFrom()
		chat := Chat{
			userID:    id,
			userState: transfer.Storage.GetUserState(id),
			message:   message,
		}
		go transfer.handle(chat)
	}

}
