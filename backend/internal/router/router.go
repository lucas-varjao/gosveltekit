// Package router sets up the HTTP routes for the application.
package router

import (
	"net/http"
	"time"

	"gosveltekit/internal/auth"
	"gosveltekit/internal/handlers"
	"gosveltekit/internal/middleware"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

const (
	authRateLimitPerSecond = 1
	authRateBurst          = 3
	apiRateLimitPerSecond  = 10
	apiRateBurst           = 20
)

// SetupRouter configures all routes for the application
func SetupRouter(
	authHandler *handlers.AuthHandler,
	authManager *auth.AuthManager,
	options ...middleware.AuthMiddlewareOptions,
) *gin.Engine {
	authOptions := middleware.DefaultAuthMiddlewareOptions()
	if len(options) > 0 {
		authOptions = options[0]
	}

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

	// Rate limiter for auth routes (brute force prevention)
	authLimiter := middleware.NewIPRateLimiter(
		rate.Limit(authRateLimitPerSecond),
		authRateBurst,
		time.Hour,
	)

	// Public auth routes
	authRoutes := r.Group("/auth")
	authRoutes.Use(middleware.RateLimitMiddleware(authLimiter))
	authRoutes.POST("/login", authHandler.Login)
	authRoutes.POST("/register", authHandler.Register)
	authRoutes.POST("/password-reset-request", authHandler.RequestPasswordReset)
	authRoutes.POST("/password-reset", authHandler.ResetPassword)

	// Rate limiter for API (more permissive)
	apiLimiter := middleware.NewIPRateLimiter(
		rate.Limit(apiRateLimitPerSecond),
		apiRateBurst,
		time.Hour,
	)

	// Protected routes
	api := r.Group("/api")
	api.Use(middleware.RateLimitMiddleware(apiLimiter))
	api.Use(middleware.AuthMiddleware(authManager, authOptions))
	// Test protected route
	api.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Esta é uma rota protegida",
		})
	})

	api.GET("/me", authHandler.GetCurrentUser)
	api.POST("/logout", authHandler.Logout)

	// Admin only routes
	admin := api.Group("/admin")
	admin.Use(middleware.RoleMiddleware("admin"))
	admin.GET("/dashboard", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Admin Dashboard",
		})
	})

	return r
}
