package server

import (
	"log"
	"net/http"

	"github.com/BrunoCardosoFS/NaxiServer/auth"
	"github.com/BrunoCardosoFS/NaxiServer/database"
	"github.com/gin-gonic/gin"
)

func (s *Server) registerApiRoutes() {
	log.Println("Registrando rotas da API...")

	api := s.router.Group("/api")
	{
		// Public Routes
		api.GET("/status", s.handleApiStatus())
		api.POST("/status", s.handleApiStatus())

		api.POST("/auth/login", s.handleLogin())

		// Protected Routes
		protected := api.Group("/")
		protected.Use(auth.AuthMiddleware())
		{

		}

		UsersRoutes := api.Group("/users")
		UsersRoutes.Use(auth.UsersMiddleware())
		{
			UsersRoutes.DELETE("/:username", s.handleDeleteUser())
			UsersRoutes.POST("/register", s.handleRegisterUser())
		}

		CatalogRoutes := api.Group("/catalog")
		CatalogRoutes.GET("/", s.handleApiGetCatalog())
		CatalogRoutes.Use(auth.CatalogMiddleware())
		{
			CatalogRoutes.POST("/", s.handleApiAddCatalogFolder())
			CatalogRoutes.DELETE("/:id", s.handleApiRemoveCatalogFolder())
		}
	}
}

func (s *Server) registerCdnRoutes() {
	cdnGroup := s.router.Group("/cdn")

	folder, err := database.GetCatalog()
	if err != nil {
		log.Println(err)
		return
	}

	for _, info := range folder {
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
