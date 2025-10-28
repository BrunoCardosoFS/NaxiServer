package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/BrunoCardosoFS/NaxiServer/models"
)

func (s *Server) registerCdnRoutes() {
	dbPath := "D:/Arquivos/Projetos/NaxiStudio/NaxiStudioApps/NaxiStudioFlow/build/db/"
	catalogPath := dbPath + "/catalog.json"

	dbServer := http.FileServer(http.Dir(dbPath))
	s.mux.Handle("/db/", http.StripPrefix("/db/", dbServer))

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

	for _, info := range categories {
		folder := http.FileServer(http.Dir(info.Path))
		prefix := "/cdn/" + info.ID + "/"

		s.mux.Handle(prefix, http.StripPrefix(prefix, folder))

		log.Printf("Server CDN: Title: %s | Prefix: %s",
			info.Title, prefix)
	}
}
