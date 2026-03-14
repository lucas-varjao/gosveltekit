// Package main is the entry point for the GoSvelteKit backend server.
package main

import (
	"fmt"
	"log/slog"
	"os"

	"gosveltekit/internal/auth"
	gormadapter "gosveltekit/internal/auth/adapter/gorm"
	"gosveltekit/internal/bootstrap"
	"gosveltekit/internal/config"
	"gosveltekit/internal/email"
	"gosveltekit/internal/handlers"
	"gosveltekit/internal/middleware"
	"gosveltekit/internal/router"
	"gosveltekit/internal/service"
	"gosveltekit/internal/version"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("Falha ao carregar as configurações")
	}

	db, err := bootstrap.OpenGorm(cfg)
	if err != nil {
		panic(
			"Falha ao conectar ao banco PostgreSQL. Verifique DATABASE_DSN/DATABASE_URL ou database.dsn no app.yml",
		)
	}

	// Initialize adapters
	userAdapter := gormadapter.NewUserAdapter(db)
	sessionAdapter := gormadapter.NewSessionAdapter(db)

	// Initialize auth manager from config
	authConfig := auth.DefaultAuthConfig()
	if cfg.Auth.SessionTTL > 0 {
		authConfig.SessionDuration = cfg.Auth.SessionTTL
	}
	if cfg.Auth.RefreshThreshold > 0 {
		authConfig.RefreshThreshold = cfg.Auth.RefreshThreshold
	}
	if cfg.Auth.MaxFailedAttempts > 0 {
		authConfig.MaxFailedAttempts = cfg.Auth.MaxFailedAttempts
	}
	if cfg.Auth.LockoutDuration > 0 {
		authConfig.LockoutDuration = cfg.Auth.LockoutDuration
	}

	authManager := auth.NewAuthManager(userAdapter, sessionAdapter, authConfig)
	authMiddlewareOptions := middleware.AuthMiddlewareOptions{
		AllowHeaderAuth: cfg.Auth.AllowHeaderAuth,
		AllowCookieAuth: cfg.Auth.AllowCookieAuth,
		CookieSecure:    cfg.Auth.CookieSecure,
	}

	// Initialize services
	emailService := email.NewEmailService(cfg)
	authService := service.NewAuthService(authManager, sessionAdapter, userAdapter, emailService)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService, cfg.Auth.CookieSecure)

	// Setup router
	r := router.SetupRouter(authHandler, authManager, authMiddlewareOptions)

	// Start server
	listenAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	slog.Info("starting server", "addr", listenAddr, "version", "v"+version.Get())
	if err := r.Run(listenAddr); err != nil {
		slog.Error("failed to start server", "err", err)
		os.Exit(1)
	}
}
