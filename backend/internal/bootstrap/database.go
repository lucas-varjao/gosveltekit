package bootstrap

import (
	"database/sql"
	"fmt"

	"gosveltekit/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// OpenGorm opens the runtime PostgreSQL connection used by the application.
func OpenGorm(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar ao banco PostgreSQL: %w", err)
	}

	return db, nil
}

// OpenSQL opens a database/sql connection for migration tooling.
func OpenSQL(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.Database.DSN)
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir conexão SQL: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("falha ao conectar ao banco PostgreSQL: %w", err)
	}

	return db, nil
}
