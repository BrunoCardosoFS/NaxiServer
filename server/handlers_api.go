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
			"status": "online",
		}

		c.JSON(http.StatusOK, response)
	}
}

func (s *Server) handleApiCatalog() gin.HandlerFunc {
	return func(c *gin.Context) {
		categories, err := database.GetCatalog()
		if err != nil {
			log.Printf("Error retrieving catalog from DB: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "The catalog could not be read."})
			return
		}

		c.JSON(http.StatusOK, categories)
	}
}

func (s *Server) handleApiAddFolderInCatalog() gin.HandlerFunc {
	return func(c *gin.Context) {
		var newCategory models.Folder
		if err := c.ShouldBindJSON(&newCategory); err != nil {
			log.Printf("Error binding JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
			return
		}

		if newCategory.ID == "" || newCategory.Title == "" || newCategory.Path == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required."})
			return
		}

		if err := database.AddFolderInCatalog(newCategory); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving to database: " + err.Error()})
			return
		}

		c.JSON(http.StatusCreated, newCategory)
	}
}
