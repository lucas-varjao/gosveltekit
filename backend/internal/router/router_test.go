// backend/internal/router/router_test.go

package router

import (
	"encoding/json"
	"gosveltekit/internal/auth"
	"gosveltekit/internal/config"
	"gosveltekit/internal/handlers"
	"gosveltekit/internal/models"
	"gosveltekit/internal/service"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

// MockAuthService implements service.AuthServiceInterface
type MockAuthService struct {
	service.AuthServiceInterface
}

func (m *MockAuthService) Login(username, password, ip, userAgent string) (*service.LoginResponse, error) {
	return &service.LoginResponse{}, nil
}

func (m *MockAuthService) Logout(userID uint, accessToken string) error {
	return nil
}

func (m *MockAuthService) RefreshToken(refreshToken string) (*service.LoginResponse, error) {
	return &service.LoginResponse{}, nil
}

func (m *MockAuthService) Register(username, email, password, displayName string) (*models.User, error) {
	return &models.User{}, nil
}

func (m *MockAuthService) RequestPasswordReset(email string) error {
	return nil
}

func (m *MockAuthService) ResetPassword(token, newPassword string) error {
	return nil
}

func NewMockAuthHandler() *handlers.AuthHandler {
	mockAuthService := &MockAuthService{}
	return handlers.NewAuthHandler(mockAuthService)
}

// MockTokenService implements a mock version of auth.TokenService
type MockTokenService struct {
	auth.TokenService
}

func NewMockTokenService() *auth.TokenService {
	cfg := &config.Config{
		JWT: config.JWTConfig{
			SecretKey:        "test-secret",
			AccessTokenTTL:   time.Hour,
			RefreshTokenTTL:  time.Hour * 24,
			PasswordResetTTL: time.Hour,
			Issuer:           "test",
		},
	}
	return auth.NewTokenService(cfg)
}

func TestSetupRouter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup
	mockAuthHandler := NewMockAuthHandler()
	mockTokenService := NewMockTokenService()
	router := SetupRouter(mockAuthHandler, mockTokenService)

	// Test cases structure
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedBody   map[string]string
	}{
		{
			name:           "Root endpoint",
			method:         "GET",
			path:           "/",
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]string{"message": "Hello GoSvelteKit"},
		},
		{
			name:           "Ping endpoint",
			method:         "GET",
			path:           "/ping",
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]string{"message": "pong"},
		},
		{
			name:           "Health endpoint",
			method:         "GET",
			path:           "/health",
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]string{"status": "ok"},
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			router.ServeHTTP(w, req)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check response body
			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			for key, expectedValue := range tt.expectedBody {
				if actualValue := response[key]; actualValue != expectedValue {
					t.Errorf("Expected %s to be %s, got %s", key, expectedValue, actualValue)
				}
			}
		})
	}
}

func TestRateLimiting(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup
	mockAuthHandler := NewMockAuthHandler()
	mockTokenService := NewMockTokenService()
	router := SetupRouter(mockAuthHandler, mockTokenService)

	// Test auth routes rate limiting
	t.Run("Auth routes rate limiting", func(t *testing.T) {
		path := "/auth/login"

		// Make multiple requests to trigger rate limiting
		for i := 0; i < 5; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", path, nil)
			req.RemoteAddr = "192.0.2.1:1234" // Set a consistent IP for testing
			router.ServeHTTP(w, req)

			if i < 3 {
				if w.Code == http.StatusTooManyRequests {
					t.Errorf("Request %d should not be rate limited", i+1)
				}
			} else {
				if w.Code != http.StatusTooManyRequests {
					t.Errorf("Request %d should be rate limited", i+1)
				}
			}
		}
	})

	// Test API routes rate limiting
	t.Run("API routes rate limiting", func(t *testing.T) {
		path := "/api/logout"

		// Add mock auth token to pass auth middleware
		authToken := "Bearer mock-token"

		// Make multiple requests to trigger rate limiting
		for i := 0; i < 25; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", path, nil)
			req.RemoteAddr = "192.0.2.1:1234" // Set a consistent IP for testing
			req.Header.Set("Authorization", authToken)
			router.ServeHTTP(w, req)

			if i < 20 {
				if w.Code == http.StatusTooManyRequests {
					t.Errorf("Request %d should not be rate limited", i+1)
				}
			} else {
				if w.Code != http.StatusTooManyRequests {
					t.Errorf("Request %d should be rate limited", i+1)
				}
			}
		}
	})
}

func TestProtectedRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup with properly initialized services
	mockAuthHandler := NewMockAuthHandler()
	mockTokenService := NewMockTokenService()
	router := SetupRouter(mockAuthHandler, mockTokenService)

	// Generate a valid token for testing
	token, _, err := mockTokenService.GenerateAccessToken(1, "user")
	if err != nil {
		t.Fatalf("Failed to generate test token: %v", err)
	}

	tests := []struct {
		name           string
		method         string
		path           string
		withAuth       bool
		token          string
		expectedStatus int
	}{
		{
			name:           "Access protected route without auth",
			method:         "POST",
			path:           "/api/logout",
			withAuth:       false,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Access protected route with valid auth",
			method:         "POST",
			path:           "/api/logout",
			withAuth:       true,
			token:          token,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Access admin route without auth",
			method:         "GET",
			path:           "/api/admin/dashboard",
			withAuth:       false,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, tt.path, nil)

			if tt.withAuth {
				req.Header.Set("Authorization", "Bearer "+tt.token)
				// Set required context values that would normally be set by middleware
				req.Header.Set("X-User-ID", "1")
				req.Header.Set("X-User-Role", "user")
			}

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
