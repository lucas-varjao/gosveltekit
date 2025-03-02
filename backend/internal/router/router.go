// router/router.go
package router

import (
	"gosveltekit/internal/auth"
	"gosveltekit/internal/handlers"
	"gosveltekit/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	authHandler *handlers.AuthHandler,
	tokenService *auth.TokenService,
	// outros handlers...
) *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Rotas públicas
	auth := r.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
	}

	// Rotas protegidas
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(tokenService))
	{
		api.POST("/logout", authHandler.Logout)

		// Rotas apenas para admin
		admin := api.Group("/admin")
		admin.Use(middleware.RoleMiddleware("admin"))
		{
			// Rotas de admin aqui
		}

		// Outras rotas protegidas
		// ...
	}

	return r
}
