package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gosveltekit/internal/bootstrap"
	"gosveltekit/internal/config"

	"github.com/pressly/goose/v3"
)

const migrationTemplate = `-- +goose Up
-- +goose StatementBegin
SELECT 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 1;
-- +goose StatementEnd
`

func main() {
	if len(os.Args) < 2 {
		exitf("usage: go run ./cmd/migrate [up|down|create <name>]")
	}

	command := os.Args[1]
	migrationsDir := filepath.Join("db", "migrations")

	switch command {
	case "create":
		if len(os.Args) < 3 {
			exitf("usage: go run ./cmd/migrate create <name>")
		}

		filename, err := createMigrationFile(migrationsDir, os.Args[2])
		if err != nil {
			exitf("failed to create migration: %v", err)
		}

		fmt.Println(filename)
	case "up", "down":
		cfg, err := config.LoadConfig()
		if err != nil {
			exitf("failed to load config: %v", err)
		}

		db, err := bootstrap.OpenSQL(cfg)
		if err != nil {
			exitf("failed to open database: %v", err)
		}
		defer db.Close()

		if err := goose.SetDialect("postgres"); err != nil {
			exitf("failed to set goose dialect: %v", err)
		}

		switch command {
		case "up":
			if err := goose.Up(db, migrationsDir); err != nil {
				exitf("failed to apply migrations: %v", err)
			}
		case "down":
			if err := goose.Down(db, migrationsDir); err != nil {
				exitf("failed to revert migration: %v", err)
			}
		}
	default:
		exitf("unknown command %q", command)
	}
}

func createMigrationFile(migrationsDir, name string) (string, error) {
	sanitizedName := strings.TrimSpace(name)
	if sanitizedName == "" {
		return "", fmt.Errorf("migration name cannot be empty")
	}

	sanitizedName = strings.ToLower(sanitizedName)
	sanitizedName = strings.ReplaceAll(sanitizedName, " ", "_")
	filename := fmt.Sprintf("%s_%s.sql", time.Now().UTC().Format("20060102150405"), sanitizedName)

	if err := os.MkdirAll(migrationsDir, 0o755); err != nil {
		return "", err
	}

	fullPath := filepath.Join(migrationsDir, filename)
	if err := os.WriteFile(fullPath, []byte(migrationTemplate), 0o644); err != nil {
		return "", err
	}

	return fullPath, nil
}

func exitf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
