// backend/internal/handlers/auth_handler.go

package handlers

import (
	"net/http"

	"gosveltekit/internal/service"
	"gosveltekit/internal/validation"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthServiceInterface
}

func NewAuthHandler(authService service.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RegistrationRequest struct {
	Username    string `json:"username" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
	DisplayName string `json:"display_name" binding:"required"`
}

type PasswordResetRequest struct {
	Token           string `json:"token" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

// Login handles user authentication with input validation
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate input data before attempting login
	if err := validation.ValidateLoginRequest(req.Username, req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get client IP and user agent
	ip := "N/A"
	ip = c.ClientIP()
	userAgent := "N/A"
	userAgent = c.Request.UserAgent()

	response, err := h.authService.Login(req.Username, req.Password, ip, userAgent)
	if err != nil {
		status := http.StatusUnauthorized
		message := "credenciais inválidas"

		if err == service.ErrUserNotActive {
			message = "usuário inativo"
		} else if err.Error() == "conta temporariamente bloqueada, tente novamente mais tarde" {
			message = err.Error()
		}

		c.JSON(status, gin.H{"error": message})
		return
	}

	c.JSON(http.StatusOK, response)
}

// RefreshToken handles token refresh with validation
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate refresh token
	if err := validation.ValidateRefreshToken(req.RefreshToken); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		status := http.StatusUnauthorized
		message := "token de atualização inválido"

		if err == service.ErrExpiredToken {
			message = "token expirado"
		}

		c.JSON(status, gin.H{"error": message})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "não autenticado"})
		return
	}

	// Obtém o token atual do contexto
	accessToken, exists := c.Get("accessToken")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token não encontrado"})
		return
	}

	if err := h.authService.Logout(userID.(uint), accessToken.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao fazer logout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logout realizado com sucesso"})
}

// Register handles new user registration with comprehensive validation
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate all registration data
	if err := validation.ValidateRegistrationRequest(
		req.Username,
		req.Email,
		req.Password,
		req.DisplayName,
	); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Forward to service layer
	user, err := h.authService.Register(req.Username, req.Email, req.Password, req.DisplayName)
	if err != nil {
		status := http.StatusBadRequest
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// RequestPasswordReset handles password reset requests
func (h *AuthHandler) RequestPasswordReset(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate email
	if err := validation.ValidateEmail(req.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.RequestPasswordReset(req.Email); err != nil {
		if err.Error() == "invalid email format" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Don't reveal if email exists for security reasons
		c.JSON(http.StatusOK, gin.H{"message": "se o email existir, um link de recuperação será enviado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "se o email existir, um link de recuperação será enviado"})
}

// ResetPassword handles password reset with token validation
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req PasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate password reset request
	if err := validation.ValidatePasswordReset(req.Token, req.NewPassword, req.ConfirmPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.ResetPassword(req.Token, req.NewPassword); err != nil {
		status := http.StatusBadRequest
		message := "falha ao redefinir senha"

		if err == service.ErrInvalidToken {
			message = "token inválido"
		} else if err == service.ErrExpiredToken {
			message = "token expirado"
		}

		c.JSON(status, gin.H{"error": message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "senha redefinida com sucesso"})
}
