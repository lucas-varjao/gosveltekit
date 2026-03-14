package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAuthHandler_ListMockPaginationItems(t *testing.T) {
	tests := []struct {
		name           string
		rawQuery       string
		expectedStatus int
	}{
		{
			name:           "offset success",
			rawQuery:       "pagination_mode=offset&page=1&page_size=5&sort=title&order=asc",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "cursor success",
			rawQuery:       "pagination_mode=cursor&page_size=5&sort=created_at&order=desc",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing mode",
			rawQuery:       "page=1",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid mode",
			rawQuery:       "pagination_mode=invalid",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "cursor unsupported sort",
			rawQuery:       "pagination_mode=cursor&page_size=5&sort=priority",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "cursor conflicting params",
			rawQuery:       "pagination_mode=cursor&page_size=5&after=a&before=b",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			rec := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(rec)
			req, _ := http.NewRequest(http.MethodGet, "/api/examples/pagination/items?"+tt.rawQuery, nil)
			c.Request = req

			handler := NewAuthHandler(&MockAuthService{})
			handler.ListMockPaginationItems(c)

			if rec.Code != tt.expectedStatus {
				t.Fatalf("expected %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestAuthHandler_ListMockPaginationItems_CursorRoundTrip(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewAuthHandler(&MockAuthService{})

	firstResponse := httptest.NewRecorder()
	firstContext, _ := gin.CreateTestContext(firstResponse)
	firstRequest, _ := http.NewRequest(
		http.MethodGet,
		"/api/examples/pagination/items?pagination_mode=cursor&page_size=4&sort=created_at&order=asc",
		nil,
	)
	firstContext.Request = firstRequest
	handler.ListMockPaginationItems(firstContext)
	if firstResponse.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", firstResponse.Code)
	}

	var pageOne struct {
		Pagination struct {
			NextCursor string `json:"next_cursor"`
		} `json:"pagination"`
	}
	if err := json.Unmarshal(firstResponse.Body.Bytes(), &pageOne); err != nil {
		t.Fatalf("unmarshal first page: %v", err)
	}
	if pageOne.Pagination.NextCursor == "" {
		t.Fatal("expected next cursor")
	}

	secondResponse := httptest.NewRecorder()
	secondContext, _ := gin.CreateTestContext(secondResponse)
	secondRequest, _ := http.NewRequest(
		http.MethodGet,
		"/api/examples/pagination/items?pagination_mode=cursor&page_size=4&sort=created_at&order=asc&after="+pageOne.Pagination.NextCursor,
		nil,
	)
	secondContext.Request = secondRequest
	handler.ListMockPaginationItems(secondContext)
	if secondResponse.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", secondResponse.Code)
	}
}
