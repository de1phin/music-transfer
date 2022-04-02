package callback

import "net/http"

type CallbackServer struct {
	HostName string
	*http.ServeMux
}

func NewCallbackServer(hostname string) *CallbackServer {
	return &CallbackServer{
		HostName: hostname,
		ServeMux: http.NewServeMux(),
	}
}

func (cs *CallbackServer) Run() {
	http.ListenAndServe(cs.HostName, cs)
}
