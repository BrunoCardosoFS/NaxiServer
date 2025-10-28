package server

import (
	"encoding/json"
	"log"
	"net/http"
)

func (s *Server) registerApiRoutes() {
	log.Println("Registrando rotas da API...")
	s.mux.HandleFunc("/api/status", s.handleApiStatus())
}

// handleApiStatus
func (s *Server) handleApiStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Println("Server API: '/api/status' accessed: " + r.Method)

		response := map[string]string{
			"status":  "online",
			"message": "NaxiServer API está funcionando!",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
