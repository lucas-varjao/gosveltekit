// Package main is the entry point for the GoSvelteKit backend server.
package main

import (
	"fmt"
	"log/slog"
	"os"

	"gosveltekit/internal/auth"
	gormadapter "gosveltekit/internal/auth/adapter/gorm"
	"gosveltekit/internal/config"
	"gosveltekit/internal/email"
	"gosveltekit/internal/handlers"
	"gosveltekit/internal/middleware"
	"gosveltekit/internal/models"
	"gosveltekit/internal/router"
	"gosveltekit/internal/service"
	"gosveltekit/internal/version"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("Falha ao carregar as configurações")
	}

	dbDSN := cfg.Database.DSN

	// Connect to PostgreSQL
	db, err := gorm.Open(postgres.Open(dbDSN), &gorm.Config{})
	if err != nil {
		panic(
			"Falha ao conectar ao banco PostgreSQL. Verifique DATABASE_DSN/DATABASE_URL ou database.dsn no app.yml",
		)
	}

	// Migrate tables (including new Session table)
	if err := db.AutoMigrate(&models.User{}, &models.Session{}); err != nil {
		panic("Falha ao migrar tabelas")
	}

	// Create admin user if not exists
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("failed to hash admin password", "err", err)
	}

	result := db.Where(models.User{Username: "admin"}).FirstOrCreate(&models.User{
		Username:     "admin",
		Email:        "onyx.views5004@eagereverest.com",
		DisplayName:  "Administrator",
		PasswordHash: string(passwordHash),
		Role:         "admin",
	})
	if result.Error != nil {
		slog.Error("failed to create or find admin user", "err", result.Error)
	}
	slog.Info("admin user ready", "rows_affected", result.RowsAffected)

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
