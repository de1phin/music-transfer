package server

import (
	"net/http"
)

type Server struct {
	HostName string
	*http.ServeMux
}

func NewServer(config Config) *Server {
	return &Server{
		HostName: config.Hostname,
		ServeMux: http.NewServeMux(),
	}
}

func (cs *Server) Run() {
	http.ListenAndServe(cs.HostName, cs)
}
