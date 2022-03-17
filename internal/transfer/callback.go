package transfer

import (
	"log"
	"net/http"
)

func (transfer *Transfer) SetUpCallbackServers(url string) {

	for _, service := range transfer.services {
		http.HandleFunc("/"+service.URLName(), func(w http.ResponseWriter, r *http.Request) {
			log.Println("Got a response on", r.URL)
		})
	}

	go http.ListenAndServe(url, nil)
}
