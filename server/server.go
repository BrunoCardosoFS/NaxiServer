package server

import (
	"log"
	"net/http"
	"time"
)

type Server struct {
	mux *http.ServeMux
}

func NewServer() *Server {
	srv := &Server{
		mux: http.NewServeMux(),
	}

	srv.registerApiRoutes()
	srv.registerCdnRoutes()
	srv.registerFrontendRoutes()

	return srv
}

func (s *Server) Start(addr string) error {
	server := &http.Server{
		Addr:              addr,
		Handler:           s.mux,
		WriteTimeout:      time.Second * 5,
		ReadHeaderTimeout: time.Second * 5,
		ReadTimeout:       time.Second * 5,
		IdleTimeout:       time.Second * 60,
	}

	log.Println("Server: Servidor online na porta", addr)
	return server.ListenAndServe()
}

func (s *Server) registerFrontendRoutes() {
	pageRoot := http.FileServer(http.Dir("./static"))
	s.mux.Handle("/", pageRoot)
}
