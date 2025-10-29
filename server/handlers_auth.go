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
			c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido: " + err.Error()})
			return
		}

		// Falta desenvolver uma validação mais robusta
		if len(creds.Password) < 8 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Senha deve ter pelo menos 8 caracteres"})
			return
		}

		hashedPassword, err := auth.HashPassword(creds.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar senha"})
			return
		}

		err = database.CreateUser(creds.Username, creds.Name, creds.Email, hashedPassword, creds.Type)

		if err != nil {
			// Falta tratar os erros
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Nome de usuário já existe ou erro no banco"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Usuário criado com sucesso"})
	}
}

func (s *Server) handleLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var creds userCredentials
		if err := c.ShouldBindJSON(&creds); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido: " + err.Error()})
			return
		}

		user, err := database.GetUserByUsername(creds.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Println("Login: Usuário não encontrado:", creds.Username)
			} else {
				log.Printf("Login: Erro no DB: %v", err)
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
			return
		}

		if !auth.CheckPasswordHash(creds.Password, user.PasswordHash) {
			log.Println("Login: Senha inválida para:", creds.Username)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciais inválidas"})
			return
		}

		tokenString, err := auth.GenerateJWT(user.Username)
		if err != nil {
			log.Printf("Login: Erro ao gerar JWT: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	}
}

func (s *Server) handleDeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		usernameToDelete := c.Param("username")
		usernameFromToken, _ := c.Get("username")
		if usernameFromToken == usernameToDelete {
			log.Println("Auth: Usuário tentando se auto-deletar:", usernameToDelete)
			c.JSON(http.StatusForbidden, gin.H{"error": "Você não pode deletar sua própria conta"})
			return
		}

		err := database.DeleteUser(usernameToDelete)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar usuário"})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Usuário deletado com sucesso"})
	}
}
