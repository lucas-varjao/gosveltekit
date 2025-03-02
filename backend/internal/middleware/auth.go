// middleware/auth.go
package middleware

import (
	"net/http"
	"strings"

	"gosveltekit/internal/auth"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(tokenService *auth.TokenService) gin.HandlerFunc {
	return func(context *gin.Context) {
		authHeader := context.GetHeader("Authorization")
		if authHeader == "" {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "autorização necessária"})
			return
		}

		// Formato: "Bearer {token}"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "formato inválido"})
			return
		}

		claims, err := tokenService.ValidateToken(parts[1])
		if err != nil {
			status := http.StatusUnauthorized
			message := "token inválido"

			if err == auth.ErrExpiredToken {
				message = "token expirado"
			}

			context.AbortWithStatusJSON(status, gin.H{"error": message})
			return
		}

		// Armazena os claims no contexto para uso posterior
		context.Set("userID", claims.UserID)
		context.Set("role", claims.Role)

		context.Next()
	}
}

// Middleware opcional para verificar roles
func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(context *gin.Context) {
		userRole, exists := context.Get("role")
		if !exists {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "usuário não autenticado"})
			return
		}

		for _, role := range roles {
			if role == userRole {
				context.Next()
				return
			}
		}

		context.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "acesso negado"})
	}
}
