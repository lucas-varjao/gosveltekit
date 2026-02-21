// backend/internal/config/config_test.go

package config

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func setupTestConfig(t *testing.T) func() {
	// Create a temporary config file for testing
	configContent := `
server:
  port: 8080
database:
  dsn: "postgresql://yaml-user:yaml-pass@localhost:5432/yaml-db?sslmode=disable"
auth:
  session_ttl: 720h
  refresh_threshold: 360h
  max_failed_attempts: 8
  lockout_duration: 45m
  allow_header_auth: true
  allow_cookie_auth: true
  cookie_secure: false
email:
  smtp_host: "smtp.yaml.local"
  smtp_port: 2525
  smtp_username: "yaml-user"
  smtp_password: "yaml-pass"
  from_email: "yaml@example.com"
  from_name: "YAML Sender"
  reset_url: "http://yaml/reset?token="
`
	err := os.MkdirAll("./configs", 0755)
	assert.NoError(t, err)

	err = os.WriteFile("./configs/app.yml", []byte(configContent), 0644)
	assert.NoError(t, err)

	// Return cleanup function
	return func() {
		os.RemoveAll("./configs")
		// Reset viper to avoid interference between tests
		viper.Reset()
		cfg = nil
	}
}

func TestLoadConfig(t *testing.T) {
	cleanup := setupTestConfig(t)
	defer cleanup()

	config, err := LoadConfig()
	assert.NoError(t, err)
	assert.NotNil(t, config)

	// Verify loaded values
	assert.Equal(t, 8080, config.Server.Port)
	assert.Equal(t, "postgresql://yaml-user:yaml-pass@localhost:5432/yaml-db?sslmode=disable", config.Database.DSN)
	assert.Equal(t, 720*time.Hour, config.Auth.SessionTTL)
	assert.Equal(t, 360*time.Hour, config.Auth.RefreshThreshold)
	assert.Equal(t, 8, config.Auth.MaxFailedAttempts)
	assert.Equal(t, 45*time.Minute, config.Auth.LockoutDuration)
	assert.True(t, config.Auth.AllowHeaderAuth)
	assert.True(t, config.Auth.AllowCookieAuth)
	assert.False(t, config.Auth.CookieSecure)
	assert.Equal(t, "smtp.yaml.local", config.Email.SMTPHost)
	assert.Equal(t, 2525, config.Email.SMTPPort)
	assert.Equal(t, "yaml-user", config.Email.SMTPUsername)
	assert.Equal(t, "yaml-pass", config.Email.SMTPPassword)
	assert.Equal(t, "yaml@example.com", config.Email.FromEmail)
	assert.Equal(t, "YAML Sender", config.Email.FromName)
	assert.Equal(t, "http://yaml/reset?token=", config.Email.ResetURL)
}

func TestLoadConfigError(t *testing.T) {
	// Reset viper to avoid interference from other tests
	viper.Reset()
	cfg = nil

	// Test with non-existent config file
	config, err := LoadConfig()
	assert.Error(t, err)
	assert.Nil(t, config)
}

func TestLoadConfigPrefersDatabaseDSNEnv(t *testing.T) {
	cleanup := setupTestConfig(t)
	defer cleanup()

	t.Setenv("DATABASE_DSN", "postgresql://env-user:env-pass@localhost:5432/env-db?sslmode=disable")
	t.Setenv("DATABASE_URL", "postgresql://url-user:url-pass@localhost:5432/url-db?sslmode=disable")

	config, err := LoadConfig()
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "postgresql://env-user:env-pass@localhost:5432/env-db?sslmode=disable", config.Database.DSN)
}

func TestLoadConfigUsesDatabaseURLAlias(t *testing.T) {
	cleanup := setupTestConfig(t)
	defer cleanup()

	t.Setenv("DATABASE_URL", "postgresql://url-user:url-pass@localhost:5432/url-db?sslmode=disable")

	config, err := LoadConfig()
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "postgresql://url-user:url-pass@localhost:5432/url-db?sslmode=disable", config.Database.DSN)
}

func TestLoadConfigWithoutFileWhenDatabaseDSNEnvIsSet(t *testing.T) {
	viper.Reset()
	cfg = nil
	t.Setenv("DATABASE_DSN", "postgresql://env-only:env-only@localhost:5432/env-only-db?sslmode=disable")

	config, err := LoadConfig()
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "postgresql://env-only:env-only@localhost:5432/env-only-db?sslmode=disable", config.Database.DSN)
}

func TestLoadConfigPrefersEnvForAllSections(t *testing.T) {
	cleanup := setupTestConfig(t)
	defer cleanup()

	t.Setenv("SERVER_PORT", "9090")
	t.Setenv("AUTH_MAX_FAILED_ATTEMPTS", "11")
	t.Setenv("AUTH_COOKIE_SECURE", "true")
	t.Setenv("EMAIL_SMTP_HOST", "smtp.env.local")
	t.Setenv("EMAIL_SMTP_PORT", "587")
	t.Setenv("EMAIL_FROM_NAME", "ENV Sender")

	config, err := LoadConfig()
	assert.NoError(t, err)
	assert.NotNil(t, config)

	assert.Equal(t, 9090, config.Server.Port)
	assert.Equal(t, 11, config.Auth.MaxFailedAttempts)
	assert.True(t, config.Auth.CookieSecure)
	assert.Equal(t, "smtp.env.local", config.Email.SMTPHost)
	assert.Equal(t, 587, config.Email.SMTPPort)
	assert.Equal(t, "ENV Sender", config.Email.FromName)
}

func TestGetConfig(t *testing.T) {
	cleanup := setupTestConfig(t)
	defer cleanup()

	// First, load the config
	_, err := LoadConfig()
	assert.NoError(t, err)

	// Then get it
	config := GetConfig()
	assert.NotNil(t, config)
	assert.Equal(t, 8080, config.Server.Port)
}

func TestGetConfigBeforeLoad(t *testing.T) {
	// Reset viper and cfg
	viper.Reset()
	cfg = nil

	// GetConfig should return nil if LoadConfig hasn't been called
	config := GetConfig()
	assert.Nil(t, config)
}
