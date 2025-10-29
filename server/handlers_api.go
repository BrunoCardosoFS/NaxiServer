package server

import (
	"log"
	"net/http"

	"github.com/BrunoCardosoFS/NaxiServer/auth"
	"github.com/BrunoCardosoFS/NaxiServer/database"
	"github.com/BrunoCardosoFS/NaxiServer/models"
	"github.com/gin-gonic/gin"
)

func (s *Server) registerApiRoutes() {
	log.Println("Registrando rotas da API...")

	api := s.router.Group("/api")
	{
		// Public Routes
		api.GET("/status", s.handleApiStatus())
		api.POST("/status", s.handleApiStatus())

		api.POST("/auth/register", s.handleRegister())
		api.POST("/auth/login", s.handleLogin())

		api.GET("/catalog", s.handleApiCatalog())

		// Protected Routes
		protected := api.Group("/")
		protected.Use(auth.AuthMiddleware())
		{
			protected.POST("/catalog", s.handleApiAddFolderInCatalog())
			protected.DELETE("/users/:username", s.handleDeleteUser())
		}
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

func (s *Server) handleApiCatalog() gin.HandlerFunc {
	return func(c *gin.Context) {
		categories, err := database.GetCatalog()
		if err != nil {
			log.Printf("Erro ao buscar catálogo no DB: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Não foi possível ler o catálogo"})
			return
		}

		c.JSON(http.StatusOK, categories)
	}
}

func (s *Server) handleApiAddFolderInCatalog() gin.HandlerFunc {
	return func(c *gin.Context) {
		var newCategory models.Folder
		if err := c.ShouldBindJSON(&newCategory); err != nil {
			log.Printf("Erro ao vincular JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido: " + err.Error()})
			return
		}

		if newCategory.ID == "" || newCategory.Title == "" || newCategory.Path == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Todos os campos são obrigatórios"})
			return
		}

		if err := database.AddFolderInCatalog(newCategory); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar no banco: " + err.Error()})
			return
		}

		c.JSON(http.StatusCreated, newCategory)
	}
}
