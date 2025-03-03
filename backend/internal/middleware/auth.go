// middleware/auth.go
package middleware

import (
	"net/http"
	"strings"

	"gosveltekit/internal/auth"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates a Gin middleware for JWT authentication.
//
// It validates the Authorization header containing a JWT token and extracts
// user information from the token claims. The middleware expects the header
// format "Bearer {token}". If validation succeeds, it adds the user ID and role
// to the request context for use by subsequent request handlers.
//
// Parameters:
//   - tokenService: A pointer to auth.TokenService used to validate tokens
//
// Returns:
//   - gin.HandlerFunc: A middleware function that can be used with Gin router
//
// Authentication failures result in:
//   - 401 Unauthorized: When header is missing or malformed
//   - 401 Unauthorized: When token is invalid or expired (with specific error messages)
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

		tokenString := parts[1]
		claims, err := tokenService.ValidateToken(tokenString)
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
		context.Set("accessToken", tokenString)

		context.Next()
	}
}

// Middleware opcional para verificar roles
// RoleMiddleware creates a new Gin middleware to verify that the current user's role
// matches at least one of the required roles. This middleware expects the user's role
// to be already set in the Gin context with the key "role" (typically by a previous
// authentication middleware).
//
// Parameters:
//   - roles: A variadic list of role names that are allowed to access the protected route.
//
// Returns:
//   - A Gin middleware function that:
//   - Returns 401 Unauthorized if no role is found in the context
//   - Returns 403 Forbidden if the user's role doesn't match any of the required roles
//   - Calls context.Next() if the user's role matches any of the required roles
//
// Usage example:
//
//	router.GET("/admin", RoleMiddleware("admin"), adminHandler)
//	router.GET("/api", RoleMiddleware("admin", "user"), apiHandler)
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
