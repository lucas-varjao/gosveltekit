// middleware/auth_test.go
package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"gosveltekit/internal/auth"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Setup test environment
func init() {
	gin.SetMode(gin.TestMode)
}

// MockTokenService embeds auth.TokenService for compatibility while overriding methods for testing
type MockTokenService struct {
	ValidateTokenFunc      func(tokenString string) (*auth.Claims, error)
	BlacklistTokenFunc     func(tokenString string) error
	IsTokenBlacklistedFunc func(tokenString string) bool
}

// ValidateToken implements the TokenService interface method
func (m *MockTokenService) ValidateToken(tokenString string) (*auth.Claims, error) {
	return m.ValidateTokenFunc(tokenString)
}

func (m *MockTokenService) BlacklistToken(tokenString string) error {
	return m.BlacklistTokenFunc(tokenString)
}

func (m *MockTokenService) IsTokenBlacklisted(tokenString string) bool {
	return m.IsTokenBlacklistedFunc(tokenString)
}

// Test cases for AuthMiddleware
func TestAuthMiddleware(t *testing.T) {
	t.Run("Missing Authorization Header", func(t *testing.T) {
		// Setup

		// Create a test router and handler to properly test middleware chain
		r := gin.New()
		r.Use(func(c *gin.Context) {
			// Convert our mock to the interface type required by AuthMiddleware
			// In a real test, we'd need to properly inject our mock implementation
			AuthMiddleware(&auth.TokenService{})(c)
		})

		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		// Make request with no Authorization header
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "autorização necessária")
	})

	t.Run("Invalid Authorization Format", func(t *testing.T) {
		// Setup test router
		r := gin.New()
		r.Use(func(c *gin.Context) {
			AuthMiddleware(&auth.TokenService{})(c)
		})

		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		// Make request with invalid authorization format
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "InvalidFormat")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "formato inválido")
	})

	t.Run("Invalid Bearer Format", func(t *testing.T) {
		// Setup test router
		r := gin.New()
		r.Use(func(c *gin.Context) {
			AuthMiddleware(&auth.TokenService{})(c)
		})

		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		// Make request with invalid bearer format (missing token part)
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "formato inválido")
	})

	// Additional test cases for auth middleware that need to use the TokenService
	// These would require a more sophisticated mock approach to implement
	// properly in a real test suite
}

// Test cases for RoleMiddleware
func TestRoleMiddleware(t *testing.T) {
	t.Run("No Role in Context", func(t *testing.T) {
		// Create test router with middleware
		r := gin.New()
		r.Use(RoleMiddleware("admin"))

		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		// Make request (no role in context)
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "usuário não autenticado")
	})

	t.Run("Role Matches Required", func(t *testing.T) {
		// Create test router that sets the role first
		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("role", "admin")
			c.Next()
		})
		r.Use(RoleMiddleware("admin", "manager"))

		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		// Make request
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("One of Multiple Roles Matches", func(t *testing.T) {
		// Create test router that sets the role first
		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("role", "user")
			c.Next()
		})
		r.Use(RoleMiddleware("admin", "user", "manager"))

		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		// Make request
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Role Doesn't Match Required", func(t *testing.T) {
		// Create test router that sets a non-matching role
		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("role", "user")
			c.Next()
		})
		r.Use(RoleMiddleware("admin", "superuser"))

		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		// Make request
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), "acesso negado")
	})

	t.Run("Empty Roles List", func(t *testing.T) {
		// Create test router with no required roles
		r := gin.New()
		r.Use(func(c *gin.Context) {
			c.Set("role", "admin")
			c.Next()
		})
		r.Use(RoleMiddleware())

		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		// Make request
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusForbidden, w.Code)
		assert.Contains(t, w.Body.String(), "acesso negado")
	})
}
