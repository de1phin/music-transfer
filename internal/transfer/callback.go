package transfer

import (
	"net/http"
)

func (transfer *Transfer) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	for _, service := range transfer.Services {
		if service.ValidAuthCallback(r) {
			userID, credentials := service.Authorize(r)
			transfer.Storage.PutServiceData(userID, service.Name(), credentials)
			user := transfer.Storage.GetUser(userID)
			user.State = Idle
			transfer.Storage.PutUser(user)
			transfer.handleLogged(Chat{user, ""})
		}
	}
}

func (transfer *Transfer) SetUpCallbackServers() {

	for _, service := range transfer.Services {
		endpoint, doSetup := service.InitCallbackServer(transfer.Config.GetCallbackURL())
		if doSetup {
			http.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
				transfer.CallbackHandler(w, r)
			})
		}
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})

	go http.ListenAndServe(transfer.Config.GetServerURL(), nil)
}
