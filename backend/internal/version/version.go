package version

import (
	"os"
	"path/filepath"
	"strings"
)

const fallbackVersion = "dev"

// BuildVersion can be injected at build time using Go ldflags.
var BuildVersion string

func Get() string {
	if version := strings.TrimSpace(BuildVersion); version != "" {
		return version
	}

	if version := strings.TrimSpace(os.Getenv("APP_VERSION")); version != "" {
		return version
	}

	if version := readVersionFile(filepath.Join("..", "VERSION")); version != "" {
		return version
	}

	return fallbackVersion
}

func readVersionFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(data))
}
