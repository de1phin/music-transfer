package transfer

type Transfer struct {
	interactor  Interactor
	storage     Storage
	services    []MusicService
	callbackURL string
}

func Run(interactor Interactor, storage Storage, services []MusicService) {

	transfer := &Transfer{
		interactor:  interactor,
		storage:     storage,
		services:    services,
		callbackURL: "localhost:8081",
	}

	transfer.SetUpCallbackServers(transfer.callbackURL)

	for {
		id, message := interactor.GetMessageFrom()
		go transfer.handle(id, storage.GetUserState(id), message)
	}

}
