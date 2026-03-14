// Package integration provides integration tests for the authentication flow.
package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gosveltekit/internal/auth"
	gormadapter "gosveltekit/internal/auth/adapter/gorm"
	"gosveltekit/internal/email"
	"gosveltekit/internal/handlers"
	"gosveltekit/internal/models"
	"gosveltekit/internal/router"
	"gosveltekit/internal/service"
	"gosveltekit/internal/testutil"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func setupIntegrationTest(t *testing.T) (*gin.Engine, *gorm.DB, *auth.AuthManager, *email.MockEmailService) {
	db := testutil.NewSQLiteTestDB(t, &models.User{}, &models.Session{})

	// Setup adapters
	userAdapter := gormadapter.NewUserAdapter(db)
	sessionAdapter := gormadapter.NewSessionAdapter(db)

	// Setup auth manager
	authConfig := auth.DefaultAuthConfig()
	authManager := auth.NewAuthManager(userAdapter, sessionAdapter, authConfig)

	// Setup services
	emailService := email.NewMockEmailService()
	authService := service.NewAuthService(authManager, sessionAdapter, userAdapter, emailService)
	authHandler := handlers.NewAuthHandler(authService)

	// Setup router
	r := router.SetupRouter(authHandler, authManager)
	return r, db, authManager, emailService
}

func TestCompleteAuthFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r, _, _, _ := setupIntegrationTest(t)

	// 1. Register user
	registration := map[string]any{
		"username":     "testuser",
		"email":        "test@example.com",
		"password":     "Test123!@#",
		"display_name": "Test User",
	}
	w := httptest.NewRecorder()
	jsonData, _ := json.Marshal(registration)
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "198.51.100.21:1234"
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 2. Login
	login := map[string]any{
		"username": "testuser",
		"password": "Test123!@#",
	}
	w = httptest.NewRecorder()
	jsonData, _ = json.Marshal(login)
	req, _ = http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "198.51.100.22:1234"
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var loginResponse map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &loginResponse)
	require.NoError(t, err)
	sessionID := loginResponse["session_id"].(string)
	assert.NotEmpty(t, sessionID)

	// 3. Access protected route via Authorization header
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/protected", nil)
	req.Header.Set("Authorization", "Bearer "+sessionID)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 4. Access protected route via X-Session-ID header
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/protected", nil)
	req.Header.Set("X-Session-ID", sessionID)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 5. Access protected route via cookie
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/protected", nil)
	req.Header.Set("Cookie", "session_id="+sessionID)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 6. Logout
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/logout", nil)
	req.Header.Set("Authorization", "Bearer "+sessionID)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 7. Attempt access after logout via Authorization (should fail)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/protected", nil)
	req.Header.Set("Authorization", "Bearer "+sessionID)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// 8. Attempt access after logout via cookie (should fail)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/protected", nil)
	req.Header.Set("Cookie", "session_id="+sessionID)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestPasswordResetFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r, db, _, mockEmailService := setupIntegrationTest(t)

	// 1. Create user directly in database
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

	// 2. Request password reset
	resetRequest := map[string]any{
		"email": "reset@example.com",
	}
	w := httptest.NewRecorder()
	jsonData, _ := json.Marshal(resetRequest)
	req, _ := http.NewRequest("POST", "/auth/password-reset-request", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "198.51.100.10:1234"
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify reset token was set and obtain plaintext token from sent email
	var updatedUser models.User
	err = db.First(&updatedUser, user.ID).Error
	require.NoError(t, err)
	assert.NotEmpty(t, updatedUser.ResetToken)

	sentEmails := mockEmailService.GetSentEmails()
	require.Len(t, sentEmails, 1)
	resetToken := sentEmails[0].Token
	require.NotEmpty(t, resetToken)

	// 3. Confirm password reset with received token
	resetConfirm := map[string]any{
		"token":            resetToken,
		"new_password":     "NewStrongPass123!",
		"confirm_password": "NewStrongPass123!",
	}
	w = httptest.NewRecorder()
	jsonData, _ = json.Marshal(resetConfirm)
	req, _ = http.NewRequest("POST", "/auth/password-reset", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "198.51.100.11:1234"
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	// 4. Old password should fail
	loginOldPassword := map[string]any{
		"username": "resetuser",
		"password": "oldpassword123",
	}
	w = httptest.NewRecorder()
	jsonData, _ = json.Marshal(loginOldPassword)
	req, _ = http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "198.51.100.12:1234"
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// 5. New password should succeed
	loginNewPassword := map[string]any{
		"username": "resetuser",
		"password": "NewStrongPass123!",
	}
	w = httptest.NewRecorder()
	jsonData, _ = json.Marshal(loginNewPassword)
	req, _ = http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "198.51.100.13:1234"
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetCurrentUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r, _, _, _ := setupIntegrationTest(t)

	// 1. Register and login
	registration := map[string]any{
		"username":     "meuser",
		"email":        "me@example.com",
		"password":     "Test123!@#",
		"display_name": "Me User",
	}
	w := httptest.NewRecorder()
	jsonData, _ := json.Marshal(registration)
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	login := map[string]any{
		"username": "meuser",
		"password": "Test123!@#",
	}
	w = httptest.NewRecorder()
	jsonData, _ = json.Marshal(login)
	req, _ = http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "198.51.100.23:1234"
	r.ServeHTTP(w, req)

	var loginResponse map[string]any
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	sessionID := loginResponse["session_id"].(string)

	// 2. Get current user
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/me", nil)
	req.Header.Set("Authorization", "Bearer "+sessionID)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	var userResponse map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &userResponse)
	require.NoError(t, err)
	assert.Equal(t, "meuser", userResponse["identifier"])
}

func TestAccountEndpointsFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r, _, _, _ := setupIntegrationTest(t)

	registration := map[string]any{
		"username":     "accountuser",
		"email":        "account@example.com",
		"password":     "Test123!@#",
		"display_name": "Account User",
	}

	w := httptest.NewRecorder()
	jsonData, _ := json.Marshal(registration)
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	login := map[string]any{
		"username": "accountuser",
		"password": "Test123!@#",
	}
	w = httptest.NewRecorder()
	jsonData, _ = json.Marshal(login)
	req, _ = http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "198.51.100.24:1234"
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var loginResponse map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &loginResponse)
	require.NoError(t, err)
	firstSessionID := loginResponse["session_id"].(string)

	// 1. Get profile
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/account/profile", nil)
	req.Header.Set("Authorization", "Bearer "+firstSessionID)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var profile map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &profile)
	require.NoError(t, err)
	assert.Equal(t, "accountuser", profile["identifier"])

	// 2. Update profile
	updateProfile := map[string]any{
		"display_name": "Account User Updated",
		"first_name":   "Account",
		"last_name":    "User",
	}
	w = httptest.NewRecorder()
	jsonData, _ = json.Marshal(updateProfile)
	req, _ = http.NewRequest("PATCH", "/api/account/profile", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+firstSessionID)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var updatedProfile map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &updatedProfile)
	require.NoError(t, err)
	assert.Equal(t, "Account User Updated", updatedProfile["display_name"])
	assert.Equal(t, "Account", updatedProfile["first_name"])
	assert.Equal(t, "User", updatedProfile["last_name"])

	// 3. Change password
	changePasswordPayload := map[string]any{
		"current_password": "Test123!@#",
		"new_password":     "ComplexN3w!A",
		"confirm_password": "ComplexN3w!A",
	}
	w = httptest.NewRecorder()
	jsonData, _ = json.Marshal(changePasswordPayload)
	req, _ = http.NewRequest("POST", "/api/account/change-password", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+firstSessionID)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	// Session should be invalid after password change.
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/account/profile", nil)
	req.Header.Set("Authorization", "Bearer "+firstSessionID)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)

	// 4. Login with new password
	loginWithNewPassword := map[string]any{
		"username": "accountuser",
		"password": "ComplexN3w!A",
	}
	w = httptest.NewRecorder()
	jsonData, _ = json.Marshal(loginWithNewPassword)
	req, _ = http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	err = json.Unmarshal(w.Body.Bytes(), &loginResponse)
	require.NoError(t, err)
	secondSessionID := loginResponse["session_id"].(string)

	// 5. Create one more session to test revocation.
	w = httptest.NewRecorder()
	jsonData, _ = json.Marshal(loginWithNewPassword)
	req, _ = http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	err = json.Unmarshal(w.Body.Bytes(), &loginResponse)
	require.NoError(t, err)
	thirdSessionID := loginResponse["session_id"].(string)

	// 6. List sessions
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/account/sessions", nil)
	req.Header.Set("Authorization", "Bearer "+thirdSessionID)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var sessions []map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &sessions)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(sessions), 2)

	// 7. Revoke one non-current session
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/account/sessions/"+secondSessionID, nil)
	req.Header.Set("Authorization", "Bearer "+thirdSessionID)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	// Ensure revoked session no longer exists in list.
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/account/sessions", nil)
	req.Header.Set("Authorization", "Bearer "+thirdSessionID)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	err = json.Unmarshal(w.Body.Bytes(), &sessions)
	require.NoError(t, err)
	for _, session := range sessions {
		assert.NotEqual(t, secondSessionID, session["id"])
	}
}

func TestAdminDashboardAccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r, db, _, _ := setupIntegrationTest(t)

	adminHash, err := bcrypt.GenerateFromPassword([]byte("Admin123!@#"), bcrypt.DefaultCost)
	require.NoError(t, err)

	adminUser := &models.User{
		Username:     "adminuser",
		Email:        "admin@example.com",
		PasswordHash: string(adminHash),
		DisplayName:  "Admin User",
		Active:       true,
		Role:         "admin",
	}
	require.NoError(t, db.Create(adminUser).Error)

	userRegistration := map[string]any{
		"username":     "regularuser",
		"email":        "regular@example.com",
		"password":     "Test123!@#",
		"display_name": "Regular User",
	}
	w := httptest.NewRecorder()
	jsonData, _ := json.Marshal(userRegistration)
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "198.51.100.31:1234"
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	regularLogin := map[string]any{
		"username": "regularuser",
		"password": "Test123!@#",
	}
	w = httptest.NewRecorder()
	jsonData, _ = json.Marshal(regularLogin)
	req, _ = http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "198.51.100.32:1234"
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var response map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	regularSession := response["session_id"].(string)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/admin/dashboard", nil)
	req.Header.Set("Authorization", "Bearer "+regularSession)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)

	adminLogin := map[string]any{
		"username": "adminuser",
		"password": "Admin123!@#",
	}
	w = httptest.NewRecorder()
	jsonData, _ = json.Marshal(adminLogin)
	req, _ = http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "198.51.100.33:1234"
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	adminSession := response["session_id"].(string)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/admin/dashboard", nil)
	req.Header.Set("Authorization", "Bearer "+adminSession)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAdminUsersListing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r, db, _, _ := setupIntegrationTest(t)

	adminHash, err := bcrypt.GenerateFromPassword([]byte("Admin123!@#"), bcrypt.DefaultCost)
	require.NoError(t, err)

	adminUser := &models.User{
		Username:     "adminlist",
		Email:        "adminlist@example.com",
		PasswordHash: string(adminHash),
		DisplayName:  "Admin List",
		Active:       true,
		Role:         "admin",
	}
	require.NoError(t, db.Create(adminUser).Error)

	users := []*models.User{
		{
			Username:     "alice",
			Email:        "alice@example.com",
			DisplayName:  "Alice Doe",
			PasswordHash: string(adminHash),
			Active:       true,
			Role:         "user",
		},
		{
			Username:     "bob",
			Email:        "bob@example.com",
			DisplayName:  "Bob Doe",
			PasswordHash: string(adminHash),
			Active:       false,
			Role:         "manager",
		},
		{
			Username:     "charlie",
			Email:        "charlie@example.com",
			DisplayName:  "Charlie Doe",
			PasswordHash: string(adminHash),
			Active:       true,
			Role:         "user",
		},
	}
	require.NoError(t, db.Create(&users).Error)

	adminLogin := map[string]any{
		"username": "adminlist",
		"password": "Admin123!@#",
	}
	w := httptest.NewRecorder()
	jsonData, _ := json.Marshal(adminLogin)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "198.51.100.40:1234"
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var authResponse map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &authResponse))
	adminSession := authResponse["session_id"].(string)

	regularLogin := map[string]any{
		"username": "alice",
		"password": "Admin123!@#",
	}
	w = httptest.NewRecorder()
	jsonData, _ = json.Marshal(regularLogin)
	req, _ = http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.RemoteAddr = "198.51.100.41:1234"
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &authResponse))
	regularSession := authResponse["session_id"].(string)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/admin/users?page=1&page_size=2&sort=email&order=asc", nil)
	req.Header.Set("Authorization", "Bearer "+regularSession)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusForbidden, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/admin/users?page=1&page_size=2&sort=email&order=asc", nil)
	req.Header.Set("Authorization", "Bearer "+adminSession)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var pageOne struct {
		Items      []map[string]any `json:"items"`
		Page       float64          `json:"page"`
		PageSize   float64          `json:"page_size"`
		TotalItems float64          `json:"total_items"`
		TotalPages float64          `json:"total_pages"`
		Sort       map[string]any   `json:"sort"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &pageOne))
	require.Len(t, pageOne.Items, 2)
	assert.Equal(t, float64(1), pageOne.Page)
	assert.Equal(t, float64(2), pageOne.PageSize)
	assert.GreaterOrEqual(t, int(pageOne.TotalItems), 4)
	assert.Equal(t, "asc", pageOne.Sort["direction"])
	assert.Equal(t, "adminlist@example.com", pageOne.Items[0]["email"])

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/admin/users?page=1&page_size=10&search=charlie&sort=display_name&order=desc", nil)
	req.Header.Set("Authorization", "Bearer "+adminSession)
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var searchResponse struct {
		Items []map[string]any `json:"items"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &searchResponse))
	require.Len(t, searchResponse.Items, 1)
	assert.Equal(t, "charlie", searchResponse.Items[0]["identifier"])
}
