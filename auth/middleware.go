package auth

import (
	"log"
	"net/http"
	"strings"

	"github.com/BrunoCardosoFS/NaxiServer/database"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Println("Auth: Header de autorização ausente")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Autorização necessária"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			log.Println("Auth: Formato 'Bearer' esperado")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Formato de token inválido"})
			return
		}

		claims, err := ValidateToken(tokenString)
		if err != nil {
			log.Printf("Auth: Token inválido: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token inválido ou expirado"})
			return
		}

		_, err = database.GetUserByUsername(claims.Username)

		if err != nil {
			log.Printf("Auth: Erro ao consultar usuário no banco de dados: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Usuário não existe"})
		}

		c.Set("username", claims.Username)

		c.Next()
	}
}
