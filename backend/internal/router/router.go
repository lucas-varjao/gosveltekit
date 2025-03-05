// backend/internal/router/router.go

package router

import (
	"gosveltekit/internal/auth"
	"gosveltekit/internal/handlers"
	"gosveltekit/internal/middleware"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func SetupRouter(
	authHandler *handlers.AuthHandler,
	tokenService *auth.TokenService,
) *gin.Engine {
	r := gin.Default()

	// Add CORS middleware
	r.Use(middleware.CorsMiddleware())

	// Root route
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello GoSvelteKit",
		})
	})

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

	// Limiter mais restritivo para rotas de autenticação (prevenção de força bruta)
	authLimiter := middleware.NewIPRateLimiter(rate.Limit(1), 3, time.Hour)

	// Rotas públicas de autenticação
	auth := r.Group("/auth")
	auth.Use(middleware.RateLimitMiddleware(authLimiter))
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/register", authHandler.Register)
		auth.POST("/password-reset-request", authHandler.RequestPasswordReset)
		auth.POST("/password-reset", authHandler.ResetPassword)
	}

	// Limiter para API normal (mais permissivo)
	apiLimiter := middleware.NewIPRateLimiter(rate.Limit(10), 20, time.Hour)

	// Rotas protegidas
	api := r.Group("/api")
	api.Use(middleware.RateLimitMiddleware(apiLimiter))
	api.Use(middleware.AuthMiddleware(tokenService))
	{
		// Rota de teste protegida
		api.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "Esta é uma rota protegida",
			})
		})

		api.POST("/logout", authHandler.Logout)

		// Rotas apenas para admin
		admin := api.Group("/admin")
		admin.Use(middleware.RoleMiddleware("admin"))
		{
			admin.GET("/dashboard", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "Admin Dashboard",
				})
			})
		}
	}

	return r
}
