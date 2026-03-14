package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"gosveltekit/internal/pagination"
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
			name:     "offset success",
			rawQuery: "pagination_mode=offset&page=2&page_size=5&search=adm&sort=email&order=asc",
			setupMock: func(m *MockAuthService) {
				m.ListAdminUsersFunc = func(
					input service.ListAdminUsersInput,
				) (*pagination.Response[service.AdminUserRow], error) {
					if input.PaginationMode != pagination.ModeOffset || input.Offset == nil {
						t.Fatalf("unexpected input: %#v", input)
					}
					if input.Offset.Page != 2 || input.Offset.PageSize != 5 || input.Offset.Search != "adm" || input.Offset.Sort != "email" || input.Offset.Order != pagination.SortAsc {
						t.Fatalf("unexpected input: %#v", input)
					}

					return &pagination.Response[service.AdminUserRow]{
						Items:          []service.AdminUserRow{},
						PaginationMode: pagination.ModeOffset,
						Sort: pagination.Sort{
							Field:     "email",
							Direction: pagination.SortAsc,
						},
						Pagination: pagination.OffsetMetadata{Page: 2, PageSize: 5, TotalPages: 1},
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "cursor success",
			rawQuery: "pagination_mode=cursor&page_size=5&search=adm&sort=email&order=asc&after=test-cursor",
			setupMock: func(m *MockAuthService) {
				m.ListAdminUsersFunc = func(
					input service.ListAdminUsersInput,
				) (*pagination.Response[service.AdminUserRow], error) {
					if input.PaginationMode != pagination.ModeCursor || input.Cursor == nil {
						t.Fatalf("unexpected input: %#v", input)
					}
					if input.Cursor.PageSize != 5 || input.Cursor.Search != "adm" || input.Cursor.Sort != "email" || input.Cursor.Order != pagination.SortAsc || input.Cursor.After != "test-cursor" {
						t.Fatalf("unexpected input: %#v", input)
					}

					return &pagination.Response[service.AdminUserRow]{
						Items:          []service.AdminUserRow{},
						PaginationMode: pagination.ModeCursor,
						Sort: pagination.Sort{
							Field:     "email",
							Direction: pagination.SortAsc,
						},
						Pagination: pagination.CursorMetadata{PageSize: 5, HasNext: true, HasPrev: true},
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing pagination mode",
			rawQuery:       "page=1",
			setupMock:      func(m *MockAuthService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid pagination mode",
			rawQuery:       "pagination_mode=unknown",
			setupMock:      func(m *MockAuthService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid query",
			rawQuery:       "pagination_mode=offset&page=abc",
			setupMock:      func(m *MockAuthService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "service validation failure",
			rawQuery: "pagination_mode=offset&page=1&page_size=10&sort=unknown",
			setupMock: func(m *MockAuthService) {
				m.ListAdminUsersFunc = func(
					input service.ListAdminUsersInput,
				) (*pagination.Response[service.AdminUserRow], error) {
					return nil, service.ErrInvalidAdminUsersQuery
				}
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "unsupported cursor sort",
			rawQuery: "pagination_mode=cursor&page_size=10&sort=display_name",
			setupMock: func(m *MockAuthService) {
				m.ListAdminUsersFunc = func(
					input service.ListAdminUsersInput,
				) (*pagination.Response[service.AdminUserRow], error) {
					return nil, service.ErrUnsupportedCursorSort
				}
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "service failure",
			rawQuery: "pagination_mode=offset&page=1&page_size=10",
			setupMock: func(m *MockAuthService) {
				m.ListAdminUsersFunc = func(
					input service.ListAdminUsersInput,
				) (*pagination.Response[service.AdminUserRow], error) {
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
