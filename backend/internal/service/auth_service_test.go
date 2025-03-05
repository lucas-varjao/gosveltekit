// backend/internal/service/auth_service_test.go

package service

import (
	"gosveltekit/internal/auth"
	"gosveltekit/internal/config"
	"gosveltekit/internal/models"
	"gosveltekit/internal/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Test helpers
func setupTest(t *testing.T) (*AuthService, *repository.UserRepository, *auth.TokenService, *gorm.DB) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.User{})
	require.NoError(t, err)

	cfg := &config.Config{
		JWT: config.JWTConfig{
			SecretKey:       "test-secret-key",
			AccessTokenTTL:  15 * time.Minute,
			RefreshTokenTTL: 24 * time.Hour,
			Issuer:          "test-issuer",
		},
	}

	userRepo := repository.NewUserRepository(db)
	tokenService := auth.NewTokenService(cfg)
	authService := NewAuthService(userRepo, tokenService)

	return authService, userRepo, tokenService, db
}

func createTestUser(t *testing.T, db *gorm.DB) *models.User {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	require.NoError(t, err)

	user := &models.User{
		Username:     "testuser",
		Email:        "test@example.com",
		DisplayName:  "Test User",
		PasswordHash: string(hashedPassword),
		Active:       true,
		Role:         "user",
	}

	err = db.Create(user).Error
	require.NoError(t, err)

	return user
}

// Tests
func TestNewAuthService(t *testing.T) {
	authService, _, _, _ := setupTest(t)
	assert.NotNil(t, authService)
	assert.NotNil(t, authService.userRepo)
	assert.NotNil(t, authService.tokenService)
	assert.NotNil(t, authService.failedLoginAttempts)
}

func TestAuthService_Login_Success(t *testing.T) {
	authService, _, _, db := setupTest(t)
	user := createTestUser(t, db)

	response, err := authService.Login("testuser", "password123", "127.0.0.1", "test-agent")

	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.AccessToken)
	assert.NotEmpty(t, response.RefreshToken)
	assert.NotZero(t, response.ExpiresAt)
	assert.Equal(t, user.ID, response.User.ID)
	assert.Equal(t, user.Username, response.User.Username)
	assert.Empty(t, response.User.PasswordHash)
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	authService, _, _, db := setupTest(t)
	_ = createTestUser(t, db)

	testCases := []struct {
		name     string
		username string
		password string
		wantErr  error
	}{
		{
			name:     "wrong password",
			username: "testuser",
			password: "wrongpass",
			wantErr:  ErrInvalidCredentials,
		},
		{
			name:     "non-existent user",
			username: "nonexistent",
			password: "password123",
			wantErr:  ErrInvalidCredentials,
		},
		{
			name:     "empty credentials",
			username: "",
			password: "",
			wantErr:  ErrInvalidCredentials,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response, err := authService.Login(tc.username, tc.password, "127.0.0.1", "test-agent")
			assert.Nil(t, response)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestAuthService_Login_AccountLocked(t *testing.T) {
	authService, _, _, db := setupTest(t)
	_ = createTestUser(t, db)

	// Attempt to login with wrong password 5 times
	for i := 0; i < 5; i++ {
		_, _ = authService.Login("testuser", "wrongpass", "127.0.0.1", "test-agent")
	}

	// Try one more time
	response, err := authService.Login("testuser", "wrongpass", "127.0.0.1", "test-agent")
	assert.Nil(t, response)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "conta temporariamente bloqueada")
}

func TestAuthService_Login_InactiveUser(t *testing.T) {
	authService, _, _, db := setupTest(t)
	user := createTestUser(t, db)

	// Deactivate user
	user.Active = false
	db.Save(user)

	response, err := authService.Login("testuser", "password123", "127.0.0.1", "test-agent")
	assert.Nil(t, response)
	assert.ErrorIs(t, err, ErrUserNotActive)
}

func TestAuthService_RefreshToken_Success(t *testing.T) {
	authService, _, _, db := setupTest(t)
	user := createTestUser(t, db)

	// First login to get a refresh token
	loginResp, err := authService.Login("testuser", "password123", "127.0.0.1", "test-agent")
	require.NoError(t, err)

	// Try to refresh the token
	refreshResp, err := authService.RefreshToken(loginResp.RefreshToken)

	require.NoError(t, err)
	assert.NotNil(t, refreshResp)
	assert.NotEmpty(t, refreshResp.AccessToken)
	assert.NotEmpty(t, refreshResp.RefreshToken)
	assert.NotZero(t, refreshResp.ExpiresAt)
	assert.Equal(t, user.ID, refreshResp.User.ID)
	assert.Empty(t, refreshResp.User.PasswordHash)
}

func TestAuthService_RefreshToken_Invalid(t *testing.T) {
	authService, _, _, _ := setupTest(t)

	response, err := authService.RefreshToken("invalid-refresh-token")
	assert.Nil(t, response)
	assert.ErrorIs(t, err, ErrInvalidCredentials)
}

func TestAuthService_Logout(t *testing.T) {
	authService, _, _, db := setupTest(t)
	user := createTestUser(t, db)

	// First login to get tokens
	loginResp, err := authService.Login("testuser", "password123", "127.0.0.1", "test-agent")
	require.NoError(t, err)

	// Logout
	err = authService.Logout(user.ID, loginResp.AccessToken)
	require.NoError(t, err)

	// Verify that refresh token is cleared
	var updatedUser models.User
	err = db.First(&updatedUser, user.ID).Error
	require.NoError(t, err)
	assert.Empty(t, updatedUser.RefreshToken)
	assert.True(t, updatedUser.RefreshTokenExpiry.IsZero())
}

func TestAuthService_Register_Success(t *testing.T) {
	authService, _, _, _ := setupTest(t)

	user, err := authService.Register("newuser", "new@example.com", "password123", "New User")

	require.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "newuser", user.Username)
	assert.Equal(t, "new@example.com", user.Email)
	assert.Equal(t, "New User", user.DisplayName)
	assert.NotEmpty(t, user.PasswordHash)
	assert.True(t, user.Active)
}

func TestAuthService_Register_DuplicateUser(t *testing.T) {
	authService, _, _, db := setupTest(t)
	_ = createTestUser(t, db)

	// Try to register with same username
	user, err := authService.Register("testuser", "another@example.com", "password123", "Another User")
	assert.Nil(t, user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "username already exists")

	// Try to register with same email
	user, err = authService.Register("anotheruser", "test@example.com", "password123", "Another User")
	assert.Nil(t, user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "email already exists")
}

func TestAuthService_RequestPasswordReset(t *testing.T) {
	authService, _, _, db := setupTest(t)
	user := createTestUser(t, db)

	err := authService.RequestPasswordReset(user.Email)
	require.NoError(t, err)

	// Verify that reset token was set
	var updatedUser models.User
	err = db.First(&updatedUser, user.ID).Error
	require.NoError(t, err)
	assert.NotEmpty(t, updatedUser.ResetToken)
	assert.False(t, updatedUser.ResetTokenExpiry.IsZero())
}

func TestAuthService_ResetPassword_Success(t *testing.T) {
	authService, _, tokenService, db := setupTest(t)
	user := createTestUser(t, db)

	// Request password reset to get a valid token
	err := authService.RequestPasswordReset(user.Email)
	require.NoError(t, err)

	// Get the stored token hash from the database
	var userWithToken models.User
	err = db.First(&userWithToken, user.ID).Error
	require.NoError(t, err)

	// Generate a valid password reset token
	plaintextToken, hashedToken, _, err := tokenService.GeneratePasswordResetToken(user.ID)
	require.NoError(t, err)

	// Update user with our known token hash
	userWithToken.ResetToken = hashedToken
	userWithToken.ResetTokenExpiry = time.Now().Add(1 * time.Hour)
	err = db.Save(&userWithToken).Error
	require.NoError(t, err)

	// Reset password using the plaintext token
	err = authService.ResetPassword(plaintextToken, "newpassword123")
	require.NoError(t, err)

	// Verify password was changed and token was cleared
	var userAfterReset models.User
	err = db.First(&userAfterReset, user.ID).Error
	require.NoError(t, err)
	assert.NotEqual(t, user.PasswordHash, userAfterReset.PasswordHash)
	assert.Empty(t, userAfterReset.ResetToken)
	assert.True(t, userAfterReset.ResetTokenExpiry.IsZero())

	// Verify we can login with the new password
	loginResp, err := authService.Login(user.Username, "newpassword123", "127.0.0.1", "test-agent")
	require.NoError(t, err)
	assert.NotNil(t, loginResp)
}

func TestAuthService_ResetPassword_InvalidToken(t *testing.T) {
	authService, _, _, _ := setupTest(t)

	err := authService.ResetPassword("invalid.token", "newpassword123")
	assert.ErrorIs(t, err, ErrInvalidToken)
}
