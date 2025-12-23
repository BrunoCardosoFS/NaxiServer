package server

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/BrunoCardosoFS/NaxiServer/auth"
	"github.com/BrunoCardosoFS/NaxiServer/database"
	"github.com/gin-gonic/gin"
)

type userCredentials struct {
	Username string `json:"user"`
	Password string `json:"password"`
}

type userCredentialsRegister struct {
	Username string `json:"user"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Type     int    `json:"type"`
}

func (s *Server) handleRegister() gin.HandlerFunc {
	return func(c *gin.Context) {
		var creds userCredentialsRegister
		if err := c.ShouldBindJSON(&creds); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
			return
		}

		// Falta desenvolver uma validação mais robusta
		if len(creds.Password) < 8 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "The password must be at least 8 characters long."})
			return
		}

		hashedPassword, err := auth.HashPassword(creds.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing password."})
			return
		}

		err = database.CreateUser(creds.Username, creds.Name, creds.Email, hashedPassword, creds.Type)

		if err != nil {
			// Falta tratar os erros
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Username already exists or there's an error in the database."})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User created successfully."})
	}
}

func (s *Server) handleLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var creds userCredentials
		if err := c.ShouldBindJSON(&creds); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
			return
		}

		user, err := database.GetUserByUsername(creds.Username)

		if err != nil {
			if err == sql.ErrNoRows {
				log.Println("Login: User not found:", creds.Username)
			} else {
				log.Printf("Login: Database Error: %v", err)
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials."})
			return
		}

		if !auth.CheckPasswordHash(creds.Password, user.PasswordHash) {
			log.Println("Login: Invalid password for", creds.Username)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		tokenString, err := auth.GenerateJWT(user.Username)
		if err != nil {
			log.Printf("Login: Error generating JWT: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token."})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": tokenString, "user": user.Username, "name": user.Name, "email": user.Email, "type": user.Type})
	}
}

func (s *Server) handleDeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		usernameToDelete := c.Param("username")
		usernameFromToken, _ := c.Get("username")
		if usernameFromToken == usernameToDelete {
			log.Println("Auth: User attempting to self-delete:", usernameToDelete)
			c.JSON(http.StatusForbidden, gin.H{"error": "You cannot delete your own account."})
			return
		}

		err := database.DeleteUser(usernameToDelete)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found."})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting user."})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User successfully deleted."})
	}
}
