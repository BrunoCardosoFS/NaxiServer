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
	// gin.SetMode(gin.ReleaseMode)
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
		Handler:           s.router,
		WriteTimeout:      time.Second * 5,
		ReadHeaderTimeout: time.Second * 5,
		ReadTimeout:       time.Second * 5,
		IdleTimeout:       time.Second * 60,
	}

	return server.ListenAndServe()
}

func (s *Server) registerFrontendRoutes() {
	// Falta implementar algo mais sofisticado
	s.router.NoRoute(gin.WrapH(http.FileServer(gin.Dir("static", false))))
}
