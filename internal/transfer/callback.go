package transfer

import (
	"log"
	"net/http"
)

func (transfer *Transfer) SetUpCallbackServers() {

	for _, service := range transfer.Services {
		http.HandleFunc("/"+service.URLName(), func(w http.ResponseWriter, r *http.Request) {
			serviceURLName := r.URL.EscapedPath()[1:]
			for _, service := range transfer.Services {
				if service.URLName() == serviceURLName {
					if service.ValidAuthCallback(r) {
						log.Println("Good callback")
						userID, credentials := service.Authorize(r)
						transfer.Storage.PutServiceData(userID, service.Name(), credentials)
						transfer.Storage.PutUserState(userID, Idle)
						transfer.handleLogged(Chat{userID, transfer.Storage.GetUserState(userID), ""})
					} else {
						log.Println("Bad callback")
					}
				}
			}
		})
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})

	go http.ListenAndServe(transfer.Config.GetServerURL(), nil)
}
