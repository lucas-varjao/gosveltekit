package version

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetPrefersBuildVersion(t *testing.T) {
	t.Setenv("APP_VERSION", "9.9.9")
	BuildVersion = "1.2.3"
	t.Cleanup(func() {
		BuildVersion = ""
	})

	if got := Get(); got != "1.2.3" {
		t.Fatalf("expected build version, got %q", got)
	}
}

func TestGetFallsBackToEnvVersion(t *testing.T) {
	t.Setenv("APP_VERSION", "2.3.4")
	BuildVersion = ""

	if got := Get(); got != "2.3.4" {
		t.Fatalf("expected env version, got %q", got)
	}
}

func TestGetFallsBackToRootVersionFile(t *testing.T) {
	BuildVersion = ""
	t.Setenv("APP_VERSION", "")

	rootDir := t.TempDir()
	backendDir := filepath.Join(rootDir, "backend")

	if err := os.MkdirAll(backendDir, 0o755); err != nil {
		t.Fatalf("failed to create backend dir: %v", err)
	}

	if err := os.WriteFile(filepath.Join(rootDir, "VERSION"), []byte("3.4.5\n"), 0o644); err != nil {
		t.Fatalf("failed to write VERSION file: %v", err)
	}

	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	if err := os.Chdir(backendDir); err != nil {
		t.Fatalf("failed to chdir to backend dir: %v", err)
	}

	t.Cleanup(func() {
		if err := os.Chdir(originalWD); err != nil {
			t.Fatalf("failed to restore working directory: %v", err)
		}
	})

	if got := Get(); got != "3.4.5" {
		t.Fatalf("expected root version file, got %q", got)
	}
}

func TestGetFallsBackToDev(t *testing.T) {
	BuildVersion = ""
	t.Setenv("APP_VERSION", "")

	rootDir := t.TempDir()
	backendDir := filepath.Join(rootDir, "backend")

	if err := os.MkdirAll(backendDir, 0o755); err != nil {
		t.Fatalf("failed to create backend dir: %v", err)
	}

	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	if err := os.Chdir(backendDir); err != nil {
		t.Fatalf("failed to chdir to backend dir: %v", err)
	}

	t.Cleanup(func() {
		if err := os.Chdir(originalWD); err != nil {
			t.Fatalf("failed to restore working directory: %v", err)
		}
	})

	if got := Get(); got != fallbackVersion {
		t.Fatalf("expected fallback version %q, got %q", fallbackVersion, got)
	}
}
