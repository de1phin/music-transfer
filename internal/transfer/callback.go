package transfer

import (
	"net/http"
)

func (transfer *Transfer) SetUpCallbackServers() {

	for _, service := range transfer.Services {
		http.HandleFunc("/"+service.URLName(), func(w http.ResponseWriter, r *http.Request) {
			serviceURLName := r.URL.EscapedPath()[1:]
			for _, service := range transfer.Services {
				if service.URLName() == serviceURLName {
					userID, credentials := service.Authorize(r)
					transfer.Storage.PutServiceData(userID, service.Name(), credentials)
					transfer.Storage.PutUserState(userID, Idle)
					break
				}
			}
		})
	}

	go http.ListenAndServe(transfer.Config.GetServerURL(), nil)
}
