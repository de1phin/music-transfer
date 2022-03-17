package transfer

type Transfer struct {
	Interactor Interactor
	Storage    Storage
	Services   []MusicService
	Config     Config
}

func (transfer *Transfer) Run() {

	transfer.SetUpCallbackServers(transfer.Config.GetCallbackURL())

	for {
		id, message := transfer.Interactor.GetMessageFrom()
		go transfer.handle(id, transfer.Storage.GetUserState(id), message)
	}

}
