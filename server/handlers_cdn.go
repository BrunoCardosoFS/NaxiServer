package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/BrunoCardosoFS/NaxiServer/models"
	"github.com/gin-gonic/gin"
)

func (s *Server) registerCdnRoutes() {
	dbPath := "D:/Arquivos/Projetos/NaxiStudio/NaxiStudioApps/NaxiStudioFlow/build/db/"
	catalogPath := dbPath + "/catalog.json"

	file, err := os.Open(catalogPath)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	var categories []models.Folder
	err = json.NewDecoder(file).Decode(&categories)
	if err != nil {
		log.Panic(err)
	}

	cdnGroup := s.router.Group("/cdn")

	for _, info := range categories {
		categoryID := info.ID
		categoryPath := info.Path
		cdnURLPrefix := "/cdn/" + categoryID

		cdnGroup.GET("/"+categoryID, s.handleCdnList(categoryPath))

		fileRoute := "/" + categoryID + "/*filepath"
		fileServerHandler := http.FileServer(http.Dir(categoryPath))
		strippedHandler := http.StripPrefix(cdnURLPrefix, fileServerHandler)

		cdnGroup.GET(fileRoute, gin.WrapH(strippedHandler))

		log.Printf("Server CDN: Title: %s | API: %s | Files: %s/* | Path: %s",
			info.Title, cdnURLPrefix, cdnURLPrefix, categoryPath)
	}
}

func (s *Server) handleCdnList(path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		entries, err := os.ReadDir(path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "The directory could not be read."})
			log.Printf("Erro ao ler diretório %s: %v", path, err)
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
