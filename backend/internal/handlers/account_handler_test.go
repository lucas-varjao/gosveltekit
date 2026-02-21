package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"gosveltekit/internal/service"

	"github.com/gin-gonic/gin"
)

func TestAuthHandler_GetAccountProfile(t *testing.T) {
	tests := []struct {
		name           string
		setupContext   func(*gin.Context)
		setupMock      func(*MockAuthService)
		expectedStatus int
	}{
		{
			name: "success",
			setupContext: func(c *gin.Context) {
				c.Set("userID", "1")
			},
			setupMock: func(m *MockAuthService) {
				m.GetProfileFunc = func(userID string) (*service.AccountProfile, error) {
					return &service.AccountProfile{ID: userID, Identifier: "testuser"}, nil
				}
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "unauthorized",
			setupContext: func(c *gin.Context) {
			},
			setupMock: func(m *MockAuthService) {
				m.GetProfileFunc = func(userID string) (*service.AccountProfile, error) {
					return nil, nil
				}
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := setupTestRouter()
			mockService := &MockAuthService{}
			tt.setupMock(mockService)
			handler := NewAuthHandler(mockService)
			tt.setupContext(c)

			handler.GetAccountProfile(c)

			if w.Code != tt.expectedStatus {
				t.Fatalf("expected %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestAuthHandler_UpdateAccountProfile(t *testing.T) {
	c, w := setupTestRouter()
	mockService := &MockAuthService{}
	mockService.UpdateProfileFunc = func(userID string, input service.UpdateProfileInput) (*service.AccountProfile, error) {
		return &service.AccountProfile{ID: userID, DisplayName: "Updated"}, nil
	}
	handler := NewAuthHandler(mockService)

	c.Set("userID", "1")
	payload := map[string]any{"display_name": "Updated"}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPatch, "/api/account/profile", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	handler.UpdateAccountProfile(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, w.Code)
	}
}

func TestAuthHandler_ChangeAccountPassword(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockAuthService)
		expectedStatus int
	}{
		{
			name: "wrong current password",
			setupMock: func(m *MockAuthService) {
				m.ChangePasswordFunc = func(userID string, input service.ChangePasswordInput) error {
					return service.ErrWrongPassword
				}
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "success",
			setupMock: func(m *MockAuthService) {
				m.ChangePasswordFunc = func(userID string, input service.ChangePasswordInput) error {
					return nil
				}
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := setupTestRouter()
			mockService := &MockAuthService{}
			tt.setupMock(mockService)
			handler := NewAuthHandler(mockService)

			c.Set("userID", "1")
			payload := map[string]any{
				"current_password": "old",
				"new_password":     "Secur3!PassA",
				"confirm_password": "Secur3!PassA",
			}
			body, _ := json.Marshal(payload)
			req, _ := http.NewRequest(http.MethodPost, "/api/account/change-password", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			handler.ChangeAccountPassword(c)

			if w.Code != tt.expectedStatus {
				t.Fatalf("expected %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestAuthHandler_ListAccountSessions(t *testing.T) {
	c, w := setupTestRouter()
	mockService := &MockAuthService{}
	mockService.ListSessionsFunc = func(userID, currentSessionID string) ([]service.SessionInfo, error) {
		return []service.SessionInfo{
			{
				ID:        currentSessionID,
				CreatedAt: time.Now().Add(-time.Hour),
				ExpiresAt: time.Now().Add(time.Hour),
				IsCurrent: true,
			},
		}, nil
	}
	handler := NewAuthHandler(mockService)

	c.Set("userID", "1")
	c.Set("sessionID", "current-session")

	handler.ListAccountSessions(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, w.Code)
	}
}

func TestAuthHandler_RevokeAccountSession(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockAuthService)
		expectedStatus int
	}{
		{
			name: "forbidden",
			setupMock: func(m *MockAuthService) {
				m.RevokeSessionFunc = func(userID, sessionID, currentSessionID string) error {
					return service.ErrAccessDenied
				}
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name: "success",
			setupMock: func(m *MockAuthService) {
				m.RevokeSessionFunc = func(userID, sessionID, currentSessionID string) error {
					return nil
				}
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, rec := setupTestRouter()
			mockService := &MockAuthService{}
			tt.setupMock(mockService)
			handler := NewAuthHandler(mockService)

			c.Set("userID", "1")
			c.Set("sessionID", "current-session")
			c.Params = gin.Params{{Key: "session_id", Value: "session-to-revoke"}}
			c.Request, _ = http.NewRequest(http.MethodDelete, "/api/account/sessions/session-to-revoke", nil)

			handler.RevokeAccountSession(c)

			if rec.Code != tt.expectedStatus {
				t.Fatalf("expected %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestAuthHandler_RevokeAccountSession_MissingParam(t *testing.T) {
	c, w := setupTestRouter()
	mockService := &MockAuthService{}
	mockService.RevokeSessionFunc = func(userID, sessionID, currentSessionID string) error {
		return errors.New("should not be called")
	}
	handler := NewAuthHandler(mockService)

	c.Set("userID", "1")
	c.Set("sessionID", "current-session")
	c.Request, _ = http.NewRequest(http.MethodDelete, "/api/account/sessions/", nil)

	handler.RevokeAccountSession(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}
