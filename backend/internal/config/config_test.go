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
  dsn: "test.db"
auth:
  session_ttl: 720h
  refresh_threshold: 360h
  max_failed_attempts: 8
  lockout_duration: 45m
  allow_header_auth: true
  allow_cookie_auth: true
  cookie_secure: false
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
	assert.Equal(t, "test.db", config.Database.DSN)
	assert.Equal(t, 720*time.Hour, config.Auth.SessionTTL)
	assert.Equal(t, 360*time.Hour, config.Auth.RefreshThreshold)
	assert.Equal(t, 8, config.Auth.MaxFailedAttempts)
	assert.Equal(t, 45*time.Minute, config.Auth.LockoutDuration)
	assert.True(t, config.Auth.AllowHeaderAuth)
	assert.True(t, config.Auth.AllowCookieAuth)
	assert.False(t, config.Auth.CookieSecure)
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
