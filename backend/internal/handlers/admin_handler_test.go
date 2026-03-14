package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"gosveltekit/internal/service"

	"github.com/gin-gonic/gin"
)

func TestAuthHandler_ListAdminUsers(t *testing.T) {
	tests := []struct {
		name           string
		rawQuery       string
		setupMock      func(*MockAuthService)
		expectedStatus int
	}{
		{
			name:     "success",
			rawQuery: "page=2&page_size=5&search=adm&sort=email&order=asc",
			setupMock: func(m *MockAuthService) {
				m.ListAdminUsersFunc = func(
					input service.ListAdminUsersInput,
				) (*service.PaginatedResult[service.AdminUserRow], error) {
					if input.Page != 2 || input.PageSize != 5 || input.Search != "adm" || input.Sort != "email" || input.Order != "asc" {
						t.Fatalf("unexpected input: %#v", input)
					}

					return &service.PaginatedResult[service.AdminUserRow]{
						Items:      []service.AdminUserRow{},
						Page:       2,
						PageSize:   5,
						TotalPages: 1,
						Sort: service.AdminUsersSort{
							Field:     "email",
							Direction: "asc",
						},
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "invalid query",
			rawQuery: "page=abc",
			setupMock: func(m *MockAuthService) {
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "service validation failure",
			rawQuery: "page=1&page_size=10&sort=unknown",
			setupMock: func(m *MockAuthService) {
				m.ListAdminUsersFunc = func(
					input service.ListAdminUsersInput,
				) (*service.PaginatedResult[service.AdminUserRow], error) {
					return nil, service.ErrInvalidAdminUsersQuery
				}
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "service failure",
			rawQuery: "page=1&page_size=10",
			setupMock: func(m *MockAuthService) {
				m.ListAdminUsersFunc = func(
					input service.ListAdminUsersInput,
				) (*service.PaginatedResult[service.AdminUserRow], error) {
					return nil, errors.New("boom")
				}
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)

			mockService := &MockAuthService{}
			tt.setupMock(mockService)
			handler := NewAuthHandler(mockService)

			req, _ := http.NewRequest(http.MethodGet, "/api/admin/users?"+tt.rawQuery, nil)
			c.Request = req

			handler.ListAdminUsers(c)

			if rec.Code != tt.expectedStatus {
				t.Fatalf("expected %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}
