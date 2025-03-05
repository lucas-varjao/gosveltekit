// backend/internal/tests/integration/auth_flow_test.go

package integration

import (
	"bytes"
	"encoding/json"
	"gosveltekit/internal/auth"
	"gosveltekit/internal/config"
	"gosveltekit/internal/handlers"
	"gosveltekit/internal/models"
	"gosveltekit/internal/repository"
	"gosveltekit/internal/router"
	"gosveltekit/internal/service"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupIntegrationTest(t *testing.T) (*gin.Engine, *gorm.DB) {
	// Configurar banco de dados de teste
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.User{})
	require.NoError(t, err)

	// Configurar serviços
	cfg := &config.Config{
		JWT: config.JWTConfig{
			SecretKey:        "test-secret-key",
			AccessTokenTTL:   15 * time.Minute,
			RefreshTokenTTL:  24 * time.Hour,
			PasswordResetTTL: 1 * time.Hour,
			Issuer:           "test-issuer",
		},
	}

	userRepo := repository.NewUserRepository(db)
	tokenService := auth.NewTokenService(cfg)
	authService := service.NewAuthService(userRepo, tokenService)
	authHandler := handlers.NewAuthHandler(authService)

	// Configurar router
	r := router.SetupRouter(authHandler, tokenService)
	return r, db
}

func TestCompleteAuthFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r, _ := setupIntegrationTest(t)

	// 1. Registro de usuário
	registration := map[string]interface{}{
		"username":     "testuser",
		"email":        "test@example.com",
		"password":     "Test123!@#",
		"display_name": "Test User",
	}
	w := httptest.NewRecorder()
	jsonData, _ := json.Marshal(registration)
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 2. Login
	login := map[string]interface{}{
		"username": "testuser",
		"password": "Test123!@#",
	}
	w = httptest.NewRecorder()
	jsonData, _ = json.Marshal(login)
	req, _ = http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var loginResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &loginResponse)
	require.NoError(t, err)
	accessToken := loginResponse["access_token"].(string)

	// 3. Acesso a rota protegida
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/protected", bytes.NewBuffer(nil))
	req.Header.Set("Authorization", "Bearer "+accessToken)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 4. Logout
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/logout", bytes.NewBuffer(nil))
	req.Header.Set("Authorization", "Bearer "+accessToken)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 5. Tentativa de acesso após logout (deve falhar)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/protected", bytes.NewBuffer(nil))
	req.Header.Set("Authorization", "Bearer "+accessToken)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestTokenRefreshFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r, _ := setupIntegrationTest(t)

	// 1. Registro e login inicial
	registration := map[string]interface{}{
		"username":     "refreshuser",
		"email":        "refresh@example.com",
		"password":     "Test123!@#",
		"display_name": "Refresh User",
	}
	w := httptest.NewRecorder()
	jsonData, _ := json.Marshal(registration)
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	login := map[string]interface{}{
		"username": "refreshuser",
		"password": "Test123!@#",
	}
	w = httptest.NewRecorder()
	jsonData, _ = json.Marshal(login)
	req, _ = http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	var loginResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &loginResponse)
	require.NoError(t, err)
	refreshToken := loginResponse["refresh_token"].(string)

	// 2. Refresh token
	refresh := map[string]interface{}{
		"refresh_token": refreshToken,
	}
	w = httptest.NewRecorder()
	jsonData, _ = json.Marshal(refresh)
	req, _ = http.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var refreshResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &refreshResponse)
	require.NoError(t, err)
	assert.NotEmpty(t, refreshResponse["access_token"])
	assert.NotEmpty(t, refreshResponse["refresh_token"])
}

func TestPasswordResetFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r, db := setupIntegrationTest(t)

	// 1. Criar usuário diretamente no banco
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("oldpassword123"), bcrypt.DefaultCost)
	require.NoError(t, err)

	user := &models.User{
		Username:     "resetuser",
		Email:        "reset@example.com",
		PasswordHash: string(hashedPassword),
		DisplayName:  "Reset User",
		Active:       true,
		Role:         "user",
	}
	err = db.Create(user).Error
	require.NoError(t, err)

	// 2. Solicitar reset de senha
	resetRequest := map[string]interface{}{
		"email": "reset@example.com",
	}
	w := httptest.NewRecorder()
	jsonData, _ := json.Marshal(resetRequest)
	req, _ := http.NewRequest("POST", "/auth/password-reset-request", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Buscar o usuário atualizado no banco
	var updatedUser models.User
	err = db.First(&updatedUser, user.ID).Error
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser.ResetToken)

	// 3. Reset de senha usando o serviço diretamente para obter o token
	cfg := &config.Config{
		JWT: config.JWTConfig{
			SecretKey:        "test-secret-key",
			PasswordResetTTL: 1 * time.Hour,
			Issuer:           "test-issuer",
		},
	}
	tokenService := auth.NewTokenService(cfg)

	// Gerar um novo token e atualizar no banco
	plaintextToken, hashedToken, expiresAt, err := tokenService.GeneratePasswordResetToken(user.ID)
	require.NoError(t, err)

	// Atualizar o token no banco
	updatedUser.ResetToken = hashedToken
	updatedUser.ResetTokenExpiry = expiresAt
	err = db.Save(&updatedUser).Error
	require.NoError(t, err)

	// Fazer o request de reset com o token
	resetPassword := map[string]interface{}{
		"token":            plaintextToken,
		"new_password":     "NewTest123!@#",
		"confirm_password": "NewTest123!@#",
	}
	w = httptest.NewRecorder()
	jsonData, _ = json.Marshal(resetPassword)
	req, _ = http.NewRequest("POST", "/auth/password-reset", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 4. Tentar login com a nova senha
	login := map[string]interface{}{
		"username": "resetuser",
		"password": "NewTest123!@#",
	}
	w = httptest.NewRecorder()
	jsonData, _ = json.Marshal(login)
	req, _ = http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
