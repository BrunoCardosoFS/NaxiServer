package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) registerApiRoutes() {
	log.Println("Registrando rotas da API...")

	api := s.router.Group("/api")
	{
		api.GET("/status", s.handleApiStatus())
		api.POST("/status", s.handleApiStatus()) // O log mostra acessos POST
	}
}

// handleApiStatus
func (s *Server) handleApiStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Server API: '/api/status' accessed: " + c.Request.Method)

		response := map[string]string{
			"status":  "online",
			"message": "NaxiServer API está funcionando!",
		}

		c.JSON(http.StatusOK, response)
	}
}
