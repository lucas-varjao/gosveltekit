package auth

import (
	"sportbetsim/internal/config"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestConfig() *config.Config {
	return &config.Config{
		JWT: config.JWTConfig{
			SecretKey:       "test-secret-key",
			AccessTokenTTL:  15 * time.Minute,
			RefreshTokenTTL: 24 * time.Hour,
			Issuer:          "test-issuer",
		},
	}
}

func TestTokenService_GenerateAccessToken(t *testing.T) {
	cfg := createTestConfig()
	tokenService := NewTokenService(cfg)

	// Test parameters
	userID := uint(123)
	role := "admin"

	// Generate token
	token, expiresAt, err := tokenService.GenerateAccessToken(userID, role)

	// Assertions
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.WithinDuration(t, time.Now().Add(cfg.JWT.AccessTokenTTL), expiresAt, time.Second)

	// Verify token can be parsed correctly
	claims := &Claims{}
	jwtToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWT.SecretKey), nil
	})

	require.NoError(t, err)
	require.True(t, jwtToken.Valid)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, role, claims.Role)
	assert.Equal(t, cfg.JWT.Issuer, claims.Issuer)
	assert.Equal(t, "123", claims.Subject)
}

func TestTokenService_ValidateToken_Success(t *testing.T) {
	cfg := createTestConfig()
	tokenService := NewTokenService(cfg)

	// Generate a valid token
	userID := uint(456)
	role := "user"
	tokenString, _, err := tokenService.GenerateAccessToken(userID, role)
	require.NoError(t, err)

	// Validate the token
	claims, err := tokenService.ValidateToken(tokenString)

	// Assertions
	require.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, role, claims.Role)
}

func TestTokenService_ValidateToken_Expired(t *testing.T) {
	cfg := createTestConfig()
	// Use a very short expiration for testing expired tokens
	cfg.JWT.AccessTokenTTL = -time.Minute // Already expired
	tokenService := NewTokenService(cfg)

	// Generate an expired token
	tokenString, _, err := tokenService.GenerateAccessToken(uint(789), "user")
	require.NoError(t, err)

	// Validate the expired token
	_, err = tokenService.ValidateToken(tokenString)

	// Assertions
	assert.ErrorIs(t, err, ErrExpiredToken)
}

func TestTokenService_ValidateToken_Invalid(t *testing.T) {
	cfg := createTestConfig()
	tokenService := NewTokenService(cfg)

	// Test with invalid token string
	_, err := tokenService.ValidateToken("invalid.token.string")

	// Assertions
	assert.ErrorIs(t, err, ErrInvalidToken)
}

func TestTokenService_GenerateRefreshToken(t *testing.T) {
	cfg := createTestConfig()
	tokenService := NewTokenService(cfg)

	// Generate refresh token
	token, expiresAt, err := tokenService.GenerateRefreshToken()

	// Assertions
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.WithinDuration(t, time.Now().Add(cfg.JWT.RefreshTokenTTL), expiresAt, time.Second)
}
