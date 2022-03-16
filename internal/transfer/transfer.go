package transfer

type Transfer struct {
	interactor Interactor
	storage    Storage
	services   []MusicService
}

func Run(interactor Interactor, storage Storage, services []MusicService) {

	transfer := &Transfer{
		interactor: interactor,
		storage:    storage,
		services:   services,
	}

	for {
		id, message := interactor.GetMessageFrom()
		go transfer.handle(id, storage.GetUserState(id), message)
	}

}
