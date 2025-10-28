package server

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
}

func NewServer() *Server {
	router := gin.Default()

	srv := &Server{
		router: router,
	}

	srv.registerApiRoutes()
	srv.registerCdnRoutes()
	srv.registerFrontendRoutes()

	return srv
}

func (s *Server) Start(addr string) error {
	log.Println("Server: Servidor online na porta", addr)

	server := &http.Server{
		Addr:              addr,
		Handler:           s.router, // <-- AQUI
		WriteTimeout:      time.Second * 5,
		ReadHeaderTimeout: time.Second * 5,
		ReadTimeout:       time.Second * 5,
		IdleTimeout:       time.Second * 60,
	}

	return server.ListenAndServe()
}

func (s *Server) registerFrontendRoutes() {
	s.router.StaticFile("/favicon.ico", "./icons/icon.ico")
	s.router.Static("/app", "./static")
}
