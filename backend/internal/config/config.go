// backend/internal/config/config.go

package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const defaultMaxFailedAttempts = 5

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type DatabaseConfig struct {
	DSN string `mapstructure:"dsn"`
}

type AuthConfig struct {
	SessionTTL        time.Duration `mapstructure:"session_ttl"`
	RefreshThreshold  time.Duration `mapstructure:"refresh_threshold"`
	MaxFailedAttempts int           `mapstructure:"max_failed_attempts"`
	LockoutDuration   time.Duration `mapstructure:"lockout_duration"`
	AllowHeaderAuth   bool          `mapstructure:"allow_header_auth"`
	AllowCookieAuth   bool          `mapstructure:"allow_cookie_auth"`
	CookieSecure      bool          `mapstructure:"cookie_secure"`
}

// EmailConfig contém configurações para envio de email
type EmailConfig struct {
	SMTPHost     string `mapstructure:"smtp_host"`
	SMTPPort     int    `mapstructure:"smtp_port"`
	SMTPUsername string `mapstructure:"smtp_username"`
	SMTPPassword string `mapstructure:"smtp_password"`
	FromEmail    string `mapstructure:"from_email"`
	FromName     string `mapstructure:"from_name"`
	ResetURL     string `mapstructure:"reset_url"`
}

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Auth     AuthConfig     `mapstructure:"auth"`
	Email    EmailConfig    `mapstructure:"email"`
}

var cfg *Config

var envConfigKeys = []string{
	"server.port",
	"database.dsn",
	"auth.session_ttl",
	"auth.refresh_threshold",
	"auth.max_failed_attempts",
	"auth.lockout_duration",
	"auth.allow_header_auth",
	"auth.allow_cookie_auth",
	"auth.cookie_secure",
	"email.smtp_host",
	"email.smtp_port",
	"email.smtp_username",
	"email.smtp_password",
	"email.from_email",
	"email.from_name",
	"email.reset_url",
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("app")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./configs")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := bindEnvConfigKeys(); err != nil {
		return nil, fmt.Errorf("falha ao vincular variáveis de ambiente: %w", err)
	}

	viper.SetDefault("auth.session_ttl", "720h")
	viper.SetDefault("auth.refresh_threshold", "360h")
	viper.SetDefault("auth.max_failed_attempts", defaultMaxFailedAttempts)
	viper.SetDefault("auth.lockout_duration", "30m")
	viper.SetDefault("auth.allow_header_auth", true)
	viper.SetDefault("auth.allow_cookie_auth", true)
	viper.SetDefault("auth.cookie_secure", false)

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return nil, fmt.Errorf("falha ao ler o arquivo de configuração: %w", err)
		}

		// Allow env-only execution when config file is absent.
		if viper.GetString("database.dsn") == "" {
			return nil, fmt.Errorf("falha ao ler o arquivo de configuração: %w", err)
		}
	}

	cfg = &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("falha ao carregar as configurações: %w", err)
	}

	return cfg, nil
}

func GetConfig() *Config {
	return cfg
}

func bindEnvConfigKeys() error {
	for _, key := range envConfigKeys {
		// Keep DATABASE_URL as a compatibility alias for database.dsn.
		if key == "database.dsn" {
			if err := viper.BindEnv(key, "DATABASE_DSN", "DATABASE_URL"); err != nil {
				return err
			}
			continue
		}

		if err := viper.BindEnv(key); err != nil {
			return err
		}
	}

	return nil
}
