package transfer

import (
	"errors"
	"log"
	"strings"
)

type Transfer struct {
	Interactor Interactor
	Storage    Storage
	Services   []MusicService
	Config     Config
}

func normalizeString(str string) string {
	return strings.ToLower(strings.Trim(str, " \n"+string(byte(13))))
}

func (transfer *Transfer) Run() {

	transfer.SetUpCallbackServers()

	for {
		id, message := transfer.Interactor.GetMessageFrom()
		chat := Chat{}
		if !transfer.Storage.HasUser(id) {
			chat = Chat{
				user: User{
					ID:          id,
					ServiceFrom: "",
					ServiceTo:   "",
					State:       Idle,
				},
				message: message,
			}
			transfer.Storage.PutUser(chat.user)
		} else {
			chat = Chat{
				user:    transfer.Storage.GetUser(id),
				message: message,
			}
		}
		go transfer.handle(chat)
	}

}

func (transfer *Transfer) getServiceByName(name string) (MusicService, error) {
	for _, service := range transfer.Services {
		if normalizeString(service.Name()) == normalizeString(name) {
			return service, nil
		}
	}

	return nil, errors.New("Invalid service")
}

func (transfer *Transfer) Transfer(user User, from MusicService, to MusicService) {
	log.Println("Transfering from", from.Name(), "to", to.Name())
	toData := transfer.Storage.GetServiceData(user.ID, to.Name())
	fromData := transfer.Storage.GetServiceData(user.ID, from.Name())
	to.AddFavourites(toData, from.GetFavourites(fromData))
	to.AddPlaylists(toData, from.GetPlaylists(fromData))
	transfer.Interactor.SendMessageTo(user.ID, "Success!")
	user.State = Idle
	transfer.Storage.PutUser(user)
}
