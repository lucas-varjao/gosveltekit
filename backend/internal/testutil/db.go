package testutil

import (
	"fmt"
	"strings"
	"sync/atomic"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var sqliteDBCounter uint64

// NewSQLiteTestDB creates an isolated in-memory SQLite database for tests.
//
// It configures a single pooled connection so tests behave consistently with
// SQLite in-memory databases, and optionally runs auto-migrations.
func NewSQLiteTestDB(t testing.TB, autoMigrateModels ...any) *gorm.DB {
	t.Helper()

	dbName := strings.NewReplacer("/", "_", " ", "_").Replace(strings.ToLower(t.Name()))
	dbID := atomic.AddUint64(&sqliteDBCounter, 1)
	dsn := fmt.Sprintf("file:%s_%d?mode=memory&cache=shared&_fk=1", dbName, dbID)

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite test database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get sql.DB from gorm database: %v", err)
	}

	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetConnMaxLifetime(0)
	sqlDB.SetConnMaxIdleTime(0)

	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	if len(autoMigrateModels) > 0 {
		if err := db.AutoMigrate(autoMigrateModels...); err != nil {
			t.Fatalf("failed to auto-migrate sqlite test database: %v", err)
		}
	}

	return db
}
