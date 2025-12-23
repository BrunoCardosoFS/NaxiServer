package server

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	dbPath string
}

func NewServer(databasePath string) *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	server := &Server{
		router: router,
		dbPath: databasePath,
	}

	server.registerApiRoutes()
	server.registerCdnRoutes()
	server.registerFrontendRoutes()

	return server
}

func (s *Server) Start(addr string) error {
	log.Println("Server: Online server on port", addr)

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
