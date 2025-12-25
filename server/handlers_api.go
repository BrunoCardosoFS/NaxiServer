package server

import (
	"log"
	"net/http"

	"github.com/BrunoCardosoFS/NaxiServer/database"
	"github.com/BrunoCardosoFS/NaxiServer/models"
	"github.com/gin-gonic/gin"
)

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

func (s *Server) handleApiGetCatalog() gin.HandlerFunc {
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

func (s *Server) handleApiAddCatalogFolder() gin.HandlerFunc {
	return func(c *gin.Context) {
		var newFolder models.Folder

		err := c.ShouldBindJSON(&newFolder)
		if err != nil {
			log.Printf("Error binding JSON: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
			return
		}

		if newFolder.ID == "" || newFolder.Title == "" || newFolder.Path == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required."})
			return
		}

		if err := database.AddFolderInCatalog(newFolder); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving to database: " + err.Error()})
			return
		}

		c.JSON(http.StatusCreated, newFolder)
	}
}

func (s *Server) handleApiRemoveCatalogFolder() gin.HandlerFunc {
	return func(c *gin.Context) {
		folderIdToDelete := c.Param("id")

		err := database.RemoveCatalogFolder(folderIdToDelete)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving to database: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Folder successfully deleted."})
	}
}
