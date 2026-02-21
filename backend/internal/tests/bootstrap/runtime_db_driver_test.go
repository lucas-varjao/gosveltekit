package bootstrap

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func backendRoot(t *testing.T) string {
	t.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to resolve runtime caller path")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", ".."))
}

func TestMainUsesPostgresDriver(t *testing.T) {
	root := backendRoot(t)
	mainPath := filepath.Join(root, "main.go")

	content, err := os.ReadFile(mainPath)
	if err != nil {
		t.Fatalf("failed to read %s: %v", mainPath, err)
	}

	mainContent := string(content)
	if !strings.Contains(mainContent, "gorm.io/driver/postgres") {
		t.Fatalf("%s must import gorm.io/driver/postgres", mainPath)
	}
	if !strings.Contains(mainContent, "postgres.Open(") {
		t.Fatalf("%s must open database with postgres.Open", mainPath)
	}
	if strings.Contains(mainContent, "gorm.io/driver/sqlite") || strings.Contains(mainContent, "sqlite.Open(") {
		t.Fatalf("%s must not reference sqlite driver in runtime bootstrap", mainPath)
	}
}

func TestRuntimeDependenciesDoNotIncludeSQLiteDriver(t *testing.T) {
	root := backendRoot(t)
	cmd := exec.Command("go", "list", "-deps", ".")
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to inspect runtime dependencies: %v\n%s", err, string(out))
	}

	dependencies := "\n" + string(out) + "\n"
	if strings.Contains(dependencies, "\ngorm.io/driver/sqlite\n") {
		t.Fatal("runtime dependencies must not include gorm.io/driver/sqlite")
	}
}
