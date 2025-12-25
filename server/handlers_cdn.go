package server

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/BrunoCardosoFS/NaxiServer/models"
	"github.com/gin-gonic/gin"
)

func (s *Server) handleCdnList(path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		entries, err := os.ReadDir(path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "The directory could not be read."})
			log.Printf("Error reading directory %s: %v", path, err)
			return
		}

		var files []models.CdnFileEntry
		for _, entry := range entries {
			files = append(files, models.CdnFileEntry{
				Name:     entry.Name(),
				FileType: strings.ToLower(filepath.Ext(entry.Name())),
				IsDir:    entry.IsDir(),
			})
		}

		c.JSON(http.StatusOK, files)
	}
}
